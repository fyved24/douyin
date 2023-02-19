package services

import (
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/fyved24/douyin/models"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

var FavoriteSrv = &FavoriteService{}

type FavoriteService struct {
}

var (
	pre_user_favorite   = "user_favorite"
	pre_user_favorited  = "user_favorited"
	pre_video_favorited = "video_favorited"
)

// FavoriteAction 点赞操作  1-点赞，2-取消点赞
func (*FavoriteService) FavoriteAction(userId int64, videoId int64, actionType int64) (err error) {
	// 1 点赞

	switch actionType {
	case 1:
		err := FavoriteSrv.addFavorite(userId, videoId)
		if err != nil {
			return err
		}
	case 2:
		err := FavoriteSrv.deleteFavorite(userId, videoId)
		if err != nil {
			return err
		}
	default:
		return errors.New("操作异常")
	}

	return nil
}

func (*FavoriteService) addFavorite(userId int64, videoId int64) error {
	// 判断用户的喜欢列表 redis缓存是否存在
	pre_key := "user_favorite"
	key := fmt.Sprintf("%s%d", pre_key, userId)
	isExist, err := models.RedisDB.SIsMember(models.Ctx, key, videoId).Result()

	if !isExist { // 缓存不存在 扫mysql
		err = models.DB.Transaction(func(tx *gorm.DB) error {
			var favoriteExit = &models.Favorite{}
			res := tx.Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteExit)
			if res.Error != nil { // favorite表中不存在
				res = tx.Create(&models.Favorite{UserID: userId, VideoID: videoId, Status: 1})
				if res.Error != nil {
					return res.Error
				}
				updateCount(models.Ctx, userId, videoId, 1)
			} else { // favorite表中存在
				if favoriteExit.Status == 0 { // 修改状态
					tx.Model(favoriteExit).Update("status", 1)
					updateCount(models.Ctx, userId, videoId, 1)
				}

			}

			return nil
		})

		// 添加到缓存
		models.RedisDB.SAdd(models.Ctx, key, videoId)
	} else {
		// 缓存存在不做任何操作
	}

	return err
}

func updateCount(ctx context.Context, UserId int64, VideoId int64, incrNum int64) {
	// 1.更新用户自己的喜欢总数
	// 2.更新视频的喜欢个数
	// 3.更新视频作者的被喜欢总数

	updateUserFavoriteCount(ctx, UserId, incrNum)
	res := updateVideoCount(ctx, VideoId, incrNum)
	updateVideoAuthorFavoritedCount(ctx, VideoId, res, incrNum)
}

// redis 更新用户自己的喜欢总数
func updateUserFavoriteCount(ctx context.Context, UserId int64, incrNum int64) {
	pre_user_count := "user_count"
	key := fmt.Sprintf("%s%d", pre_user_count, UserId)
	//val, err := models.RedisDB.HGetAll(ctx, key).Result()
	isExist, err := models.RedisDB.Exists(ctx, key).Result()
	var user = new(models.User)
	if err == nil {
		if isExist == 0 { // 对应的user_count 不存在查找
			// 分布式锁
			lock, err := models.RedisLock.Obtain(ctx, key, 100*time.Millisecond, nil)
			if err == redislock.ErrNotObtained {
				fmt.Println("Could not obtain lock!")
			} else if err != nil {
				log.Fatalln(err)
			}
			defer lock.Release(ctx)
			// 释放锁
			models.DB.Where("id=?", user.ID).First(&user)
			// 将查到的数据存入redis
			user_count := make(map[string]interface{})
			user_count["favorite"] = incrNum + int64(user.FavoriteCount)
			user_count["favorited"] = user.TotalFavorited

			models.RedisDB.HMSet(ctx, key, user_count)
		} else {
			models.RedisDB.HIncrBy(ctx, key, "favorite", incrNum)
		}
	}
}

// 添加视频的喜欢总数
func updateVideoCount(ctx context.Context, VideoId int64, incrNum int64) map[string]interface{} {
	resMap := make(map[string]interface{})
	pre_key := "video_count"
	key := fmt.Sprintf("%s%d", pre_key, VideoId)
	rows, err := models.RedisDB.Exists(ctx, key).Result()
	if err == nil { // 不存在
		if rows == 0 {
			lock, err := models.RedisLock.Obtain(ctx, key, 100*time.Millisecond, nil)
			if err == redislock.ErrNotObtained {
				fmt.Println("Could not obtain lock!")
			} else if err != nil {
				log.Fatalln(err)
			}
			defer lock.Release(ctx)

			var video = new(models.Video)
			models.DB.Where("id = ?", VideoId).First(&video)
			// 将查到的数据放入到video_count+id : count_num ; 0不会过期
			models.RedisDB.Set(ctx, key, video.FavoriteCount+incrNum, 0).Err()
			// 把作者信息放入到map防止2次查表
			resMap["AuthorID"] = video.AuthorID
		} else { // 缓存在
			models.RedisDB.IncrBy(ctx, key, incrNum)
		}

	}
	return resMap
}

//
func updateVideoAuthorFavoritedCount(ctx context.Context, VideoId int64, resMap map[string]interface{}, incrNum int64) {
	// 先查询视频所属用户
	UserId, ok := resMap["AuthorID"] //
	pre_user_count := "user_count"
	if !ok {
		var video = new(models.Video)
		models.DB.Where("id=?", VideoId).First(&video)
		UserId = video.AuthorID
	}

	key := fmt.Sprintf("%s%d", pre_user_count, UserId)
	//val, err := models.RedisDB.HGetAll(ctx, key).Result()
	isExist, err := models.RedisDB.Exists(ctx, key).Result()
	var user = new(models.User)
	if err == nil { // 对应的user_count 不存在查找
		if isExist == 0 { // 缓存不存在
			// 分布式锁
			lock, err := models.RedisLock.Obtain(ctx, key, 100*time.Millisecond, nil)
			if err == redislock.ErrNotObtained {
				fmt.Println("Could not obtain lock!")
			} else if err != nil {
				log.Fatalln(err)
			}
			defer lock.Release(ctx)
			// 释放锁
			models.DB.Where("id=?", UserId).First(&user)
			// 将查到的数据存入redis
			user_count := make(map[string]interface{})
			user_count["favorite"] = user.FavoriteCount
			user_count["favorited"] = int64(user.TotalFavorited) + incrNum

			models.RedisDB.HMSet(ctx, key, user_count)
		} else { // 缓存存在
			models.RedisDB.HIncrBy(ctx, key, "favorited", incrNum)
		}
	}
}

type Pipeline struct {
	sync.Mutex
	redis.Pipeliner
}

func (*FavoriteService) deleteFavorite(userId int64, videoId int64) error {
	preKey := "user_favorite"
	key := fmt.Sprintf("%s%d", preKey, userId)
	flag, err := models.RedisDB.SIsMember(models.Ctx, key, videoId).Result()
	if flag { // 如果缓存存在
		// 删除用户喜欢列表的缓存
		models.RedisDB.SRem(models.Ctx, key, videoId)
		// 更新count
		updateCount(models.Ctx, userId, videoId, -1)
	} else {
		var favoriteActionCancel = new(models.Favorite)

		err = models.DB.Transaction(func(tx *gorm.DB) error {
			res := tx.Where("user_id = ? AND video_id= ?", userId, videoId).First(&favoriteActionCancel)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return errors.New("未查询到记录")
			}
			// 记录存在
			if favoriteActionCancel.Status == 1 {
				// 点赞输目-1
				updateCount(models.Ctx, userId, videoId, -1)
				// 可以放到MQ去进一步加速
				tx.Model(favoriteActionCancel).Update("status", 0)
			}
			// status=0 不用管
			return nil
		})
	}

	return err
}

// FindAllFavorite returns a list of Favorite videos.
func (*FavoriteService) FindAllFavorite(userId int64) ([]models.Video, error) {
	var favoriteList []models.Favorite
	videoList := make([]models.Video, 0)
	if err := models.DB.Where("user_id = ? AND status = ?", userId, 1).Find(&favoriteList).Error; err != nil {
		// 找不到记录
		return videoList, nil
	}

	for _, m := range favoriteList {
		var video = models.Video{}
		if err := models.DB.Where("id = ?", m.VideoID).First(&video).Error; err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}

//https://juejin.cn/post/7027347979065360392#heading-20
