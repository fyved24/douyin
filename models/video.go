package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	AuthorID      uint `json:"-"`
	Author        User `gorm:"foreignKey:AuthorID"`
	PlayUrl       string
	CoverUrl      string
	FavoriteCount uint
	CommentCount  uint
	IsFavorite    bool
	Title         string
	Comments      []Comment
}

func QueryFeedVideoListByLatestTime(limit int, latestTime time.Time) (*[]Video, error) {
	var videos []Video
	err := DB.Model(&Video{}).Preload("Author").Where("created_at<?", latestTime).Limit(limit).Find(&videos).Error
	return &videos, err
}
