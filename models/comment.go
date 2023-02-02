package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	VideoID int64
	UserID  uint
	Content string
}
