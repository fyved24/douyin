package models

import (
	"errors"
	"time"
)

type Video struct {
	Model
	AuthorID      uint      `json:"-"`
	Author        User      `gorm:"foreignKey:AuthorID" json:"author"`
	PlayUrl       string    `json:"play_url"`
	CoverUrl      string    `json:"cover_url"`
	FavoriteCount int64     `json:"favorite_count"`
	CommentCount  int64     `json:"comment_count"`
	IsFavorite    bool      `json:"is_favorite"`
	Title         string    `json:"title"`
	Comments      []Comment `json:"comments"`
}

func QueryFeedVideoListByLatestTime(limit int, latestTime time.Time, userID uint) (*[]Video, error) {
	var videos []Video
	err := DB.Model(&Video{}).Preload("Author").Where("created_at<?", latestTime).Limit(limit).Find(&videos).Error
	if userID == 0 {
		return &videos, err

	}
	for i := 0; i < len(videos); i++ {
		var favorite Favorite
		if DB.Model(&Favorite{}).Where("user_id = ? AND video_id = ?", userID, videos[i].ID).First(&favorite).Error == nil {
			if favorite.UserID != 0 && favorite.Status == 1 {
				videos[i].IsFavorite = true
			}
		}

	}
	return &videos, err
}
func SaveVideo(video *Video) error {
	var user User
	err := DB.Model(&User{}).Where("id = ?", video.AuthorID).First(&user).Error
	if err != nil {
		return err
	}
	user.WorkCount = user.WorkCount + 1
	err = DB.Create(video).Error
	if err != nil {
		return errors.New("视频保存出错")
	}
	err = DB.Model(&user).Update("work_count", user.WorkCount).Error
	if err != nil {
		return errors.New("修改作品数量出错")
	}
	return nil
}

func QueryUserVideoList(userID uint) (*[]Video, error) {
	var videos []Video
	err := DB.Model(&Video{}).Where("author_id=?", userID).Find(&videos).Error
	for i := 0; i < len(videos); i++ {
		var favorite Favorite
		if DB.Model(&Favorite{}).Where("user_id = ? AND video_id = ?", userID, videos[i].ID).First(&favorite).Error == nil {
			if favorite.UserID != 0 && favorite.Status == 1 {
				videos[i].IsFavorite = true
			}
		}

	}
	return &videos, err
}
