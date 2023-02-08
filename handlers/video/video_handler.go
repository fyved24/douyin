package video

import (
	"net/http"

	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/services/video"
	"github.com/gin-gonic/gin"
)

func FeedVideoList(c *gin.Context) {

	req := requests.NewDouyinFeedRequest(c)
	feedVideosRes, _ := video.FeedVideoList(req.LatestTime)
	c.JSON(http.StatusOK, feedVideosRes)
}
