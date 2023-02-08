package models

import (
	"time"
)

type Video struct {
	Model
	AuthorID      uint      `json:"-"`
	Author        User      `gorm:"foreignKey:AuthorID"`
	PlayUrl       string    `json:"play_url"`
	CoverUrl      string    `json:"cover_url"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	IsFavorite    bool      `json:"is_favorite"`
	Title         string    `json:"title"`
	Comments      []Comment `json:"comments"`
}

func QueryFeedVideoListByLatestTime(limit int, latestTime time.Time) (*[]Video, error) {
	var videos []Video
	err := DB.Model(&Video{}).Preload("Author").Where("created_at<?", latestTime).Limit(limit).Find(&videos).Error
	return &videos, err
}
func SaveVideo(video *Video) {
	DB.Create(video)
}

func QueryUserVideoList(userID uint) (*[]Video, error) {
	var videos []Video
	err := DB.Model(&Video{}).Where("author_id=", userID).Find(&videos).Error
	return &videos, err
}
