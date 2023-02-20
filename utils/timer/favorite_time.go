package utils

import (
	"github.com/fyved24/douyin/models"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

const LAYOUT = "2006-01-02 15:04:05"

// RefreshRedisToDB RefreshTime 定时任务的启动方法
func RefreshRedisToDB() {
	log.Println("定时刷新RedisCount2Mysql开始 ：" + time.Now().Format(LAYOUT))
	defer log.Println("定时刷新RedisCount2Mysql结束 ：" + time.Now().Format(LAYOUT))
	RefreshRedisToDBCount("user")
	RefreshRedisToDBCount("video")

}

// RefreshRedisToDBCount 将redis中的count favorite等字段刷新到数据库
func RefreshRedisToDBCount(updateCount string) {
	//user_count+id:{favorite:favoriteNum, favorited:favoritedNum}
	var cursor uint64

	refString := updateCount + "_count*"
	for {
		var keys []string
		var err error
		// 每次扫描100个
		keys, cursor, err = models.RedisDB.Scan(models.Ctx, cursor, refString, 100).Result()
		if err != nil {
			panic(err)
		}
		// 一次更新100
		if updateCount == "user" {
			UpdateUserCount(keys)
		} else {
			UpdateVideoCount(keys)
		}
		// 没有更多key了
		if cursor == 0 {
			break
		}
	}
}

func UpdateUserCount(keys []string) error {
	models.DB.Transaction(func(tx *gorm.DB) error {

		updateVars := make(map[string]interface{})
		for _, key := range keys {
			vars, err := models.RedisDB.HGetAll(models.Ctx, key).Result()
			if err != nil {
				continue
			}
			updateVars["favorite_count"] = vars["favorite"]
			updateVars["total_favorited"] = vars["favorited"]
			id := strings.TrimPrefix(key, "user_count")
			err = tx.Debug().Model(&models.User{}).Where("id = ? ", id).Updates(updateVars).Error
			if err == nil {
				// 删除缓存
				//models.RedisDB.Del(models.Ctx, key)
			}
		}
		return nil
	})

	return nil
}

func UpdateVideoCount(keys []string) error {
	vars, err := models.RedisDB.MGet(models.Ctx, keys...).Result()
	if err != nil {
		return err
	}
	log.Println("异步刷新", keys)
	// 通过key解析id
	videoIds := []string{}
	for _, key := range keys {
		videoIds = append(videoIds, strings.TrimPrefix(key, "video_count"))
	}
	// 批量更新mysql
	models.DB.Transaction(func(tx *gorm.DB) error {
		//var video = new(models.Video)
		for idx, id := range videoIds {
			err := tx.Debug().Model(&models.Video{}).Where("id = ? ", id).Update("favorite_count", vars[idx]).Error
			if err == nil {
				// 删除缓存
				//models.RedisDB.Del(models.Ctx, keys[idx])
			}
		}
		return nil
	})
	return nil
}
