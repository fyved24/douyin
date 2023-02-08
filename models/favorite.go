package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID  int64
	VideoID int64
	Status  int64 //  1-点赞，0-未点赞
}

type FavoriteActionResponse struct {
	StatusCode int32   `json:"status_code,omitempty"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg,omitempty"`  // 返回状态描述
}

type FavoriteListResponse struct {
	StatusCode int32   `json:"status_code,omitempty"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg,omitempty"`  // 返回状态描述
	VideoList  []Video `json:"video_list,omitempty"`  // 用户点赞视频列表
}
