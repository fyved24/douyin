package responses

import (
	"github.com/fyved24/douyin/models"
)

type DouyinFeedResponse struct {
	CommonResponse
	VideoList *[]models.Video `json:"video_list"`
	NextTime  int64           `json:"next_time,omitempty"`
}
