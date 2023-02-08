package models

type Comment struct {
	Model
	VideoID uint   `json:"video_id"`
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}
