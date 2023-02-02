package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID  int64
	VideoID int64
	status  bool //  1-点赞，0-未点赞
}
