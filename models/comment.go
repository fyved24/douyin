package models

type Comment struct {
	Model
	VideoID int64  `json:"video_id"`
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}
