package services

import (
	"github.com/fyved24/douyin/models"
)

// FavoriteAction 点赞操作  1-点赞，2-取消点赞
func FavoriteAction(userId int64, videoId int64, actionType int64) (err error) {
	// 1 点赞
	if actionType == 1 {
		favorite := models.Favorite{
			UserID:  userId,
			VideoID: videoId,
			Status:  1,
		}

		var favoriteExist = &models.Favorite{}

		result := models.DB.Model(&models.Favorite{}).
			Where("user_id=? AND video_id=?", userId, videoId).
			First(&favoriteExist)

		if result.Error != nil { // 不存在
			models.DB.Model(&models.Favorite{}).
				Create(&favorite)

			// videoId对应的userId的TotalFavorited增加
			videoUserId, err := GetVideoAuthor(videoId)
			if err != nil {
				return err
			}

			if err := AddTotalFavorited(videoUserId); err != nil {
				return err
			}
		} else { // 喜欢记录存在
			if favoriteExist.Status == 0 { // 更改喜欢记录
				// 视频的喜欢+1
				// 视频的作者被喜欢+1
				// 喜欢记录的状态更新为1

			}

			//status == 1  video的favorite_count不变
		}

	}

}

func AddTotalFavorited(id interface{}) interface{} {

	return nil
}

func GetVideoAuthor(id int64) (interface{}, interface{}) {
	return nil, nil
}
