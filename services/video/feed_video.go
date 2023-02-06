package video

import (
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"log"
	"time"
)

func FeedVideoList(latestTime time.Time) (*responses.DouyinFeedResponse, error) {
	videos, err := models.QueryFeedVideoListByLatestTime(30, latestTime)
	nextTime := time.Now().UnixNano() / 1e6
	if len(*videos) > 0 {
		nextTime = (*videos)[0].CreatedAt.UnixNano() / 1e6
	}
	log.Printf("next time %v", nextTime)
	return &responses.DouyinFeedResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0},
		VideoList:      videos,
		NextTime:       nextTime,
	}, err
}
