package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID  uint
	VideoID uint
	status  bool //  1-点赞，0-未点赞
}
