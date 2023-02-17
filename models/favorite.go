package models

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserID  int64
	VideoID int64
	Status  int64 //  1-点赞，0-未点赞
}

// SelectFavoriteCountByID 根据用户ID查找某个用户点赞过的视频数量
func SelectFavoriteCountByID(id uint) uint {
	return 0
}
