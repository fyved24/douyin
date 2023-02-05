package handlers

import (
	"net/http"

	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/services"
	"github.com/gin-gonic/gin"
)

func FeedVideoList(c *gin.Context) {

	req := requests.NewDouyinFeedRequest(c)
	feedVideosRes, _ := services.FeedVideoList(req.LatestTime)
	c.JSON(http.StatusOK, feedVideosRes)
}
