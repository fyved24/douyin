package services

import (
	"errors"
	"github.com/fyved24/douyin/models"
	"gorm.io/gorm"
)

// FavoriteAction 点赞操作  1-点赞，2-取消点赞
func FavoriteAction(userId int64, videoId int64, actionType int64) (err error) {
	// 1 点赞

	switch actionType {
	case 1:
		err := addFavorite(userId, videoId)
		if err != nil {
			return err
		}
	case 2:
		err := deleteFavorite(userId, videoId)
		if err != nil {
			return err
		}
	default:
		return errors.New("操作异常")
	}

	return nil
}

func addFavorite(userId int64, videoId int64) error {
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		// 1.新增点赞操作
		video := new(models.Video)

		if err := tx.First(&video, videoId).Error; err != nil {
			return err
		}
		var favoriteExit = &models.Favorite{}
		res := tx.Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteExit)
		if res.Error != nil { // 不存在
			res = tx.Create(&models.Favorite{UserID: userId, VideoID: videoId, Status: 1})
			if res.Error != nil {
				return res.Error
			}
			// 2.改变 video 表中的 favorite count
			res = tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
			if res.Error != nil {
				return res.Error
			}
		} else { // 存在
			if favoriteExit.Status == 0 { // 修改状态， count +1
				tx.Model(favoriteExit).Update("status", 1)

				// 2.改变 video 表中的 favorite count
				res := tx.Model(video).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
				if res.Error != nil {
					return res.Error
				}
			}
		}

		return nil
	})
	return err
}

func deleteFavorite(userId int64, videoId int64) error {
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		//var favoriteCancel = &models.Favorite{}
		favoriteActionCancel := models.Favorite{
			UserID:  userId,
			VideoID: videoId,
			Status:  0, //0-未点赞
		}

		res := tx.Model(favoriteActionCancel).Where("user_id = ? AND video_id= ?", userId, videoId).Update("status", 0)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return errors.New("未查询到记录")
		}
		video := new(models.Video)
		// 2.改变 video 表中的 favorite count
		res = tx.Model(video).Where("id = ? AND favorite_count >0", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}
		return nil
	})
	return err
}

// FindAllFavorite returns a list of Favorite videos.
func FindAllFavorite(userId int64) ([]models.Video, error) {
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
