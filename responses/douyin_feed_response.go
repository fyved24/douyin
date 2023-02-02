package responses

import (
	"github.com/fyved24/douyin/models"
)

type DouyinFeedResponse struct {
	CommonResponse
	VideoList *[]models.Video
	NextTime  int64
}
