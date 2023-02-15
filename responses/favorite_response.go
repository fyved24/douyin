package responses

import "github.com/fyved24/douyin/models"

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type FavoriteListResponse struct {
	StatusCode int32          `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string         `json:"status_msg"`  // 返回状态描述
	VideoList  []models.Video `json:"video_list"`  // 用户点赞视频列表
}
