package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	VideoID uint
	UserID  uint
	Content string
}
