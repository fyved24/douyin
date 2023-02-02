package models

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	gorm.Model
	AuthorID      uint `json:"-"`
	Author        User `gorm:"foreignKey:AuthorID"`
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
	Comments      []Comment
}

func QueryFeedVideoListByLatestTime(limit int, latestTime time.Time) (*[]Video, error) {
	var videos []Video
	err := DB.Model(&Video{}).Preload("Author").Where("created_at<?", latestTime).Limit(limit).Find(&videos).Error
	return &videos, err
}
