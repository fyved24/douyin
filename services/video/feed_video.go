package video

import (
	"context"
	"encoding/json"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"log"
	"time"
)

func FeedVideoList(latestTime time.Time, userID uint) (*responses.DouyinFeedResponse, error) {
	ctx := context.Background()
	redisClient := models.RedisDB
	var videos *[]models.Video
	var cacheVideos []models.Video
	s, err := redisClient.Get(ctx, "feed_list").Result()
	log.Printf("缓存列表 %v", s)

	if err != nil {
		log.Printf("feed list 无缓存 %v", err)
		videos, err = models.QueryFeedVideoListByLatestTime(30, latestTime, userID)
		videosStr, _ := json.Marshal(videos)
		_, err = redisClient.Set(ctx, "feed_list", videosStr, 5*time.Second).Result()
		if err != nil {
			log.Printf("设置缓存feed列表出错 %v", err)
		}
	} else {
		err := json.Unmarshal([]byte(s), &cacheVideos)
		if err != nil {
			log.Printf("反序列化 %v", err)
		}
		log.Printf("获取feed缓存 %v", cacheVideos)
		videos = &cacheVideos
	}

	nextTime := time.Now().UnixNano() / 1e6
	if len(*videos) > 0 {
		nextTime = (*videos)[0].CreatedAt.UnixNano() / 1e6
	}
	log.Printf("next time %v", nextTime)
	return &responses.DouyinFeedResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0, StatusMsg: "success"},
		VideoList:      videos,
		NextTime:       nextTime,
	}, err
}

func QueryUserVideoList(userID uint) (*responses.DouyinFeedResponse, error) {
	videos, err := models.QueryUserVideoList(userID)
	return &responses.DouyinFeedResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0, StatusMsg: "success"},
		VideoList:      videos,
	}, err
}
