package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID  int64
	VideoID int64
	Status  int64 //  1-点赞，0-未点赞
}
