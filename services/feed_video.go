package services

import (
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"time"
)

func FeedVideoList(latestTime time.Time) (*responses.DouyinFeedResponse, error) {
	videos, err := models.QueryFeedVideoListByLatestTime(10, latestTime)
	nextTime := time.Now().Unix() / 1e6
	return &responses.DouyinFeedResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0},
		VideoList:      videos,
		NextTime:       nextTime,
	}, err
}
