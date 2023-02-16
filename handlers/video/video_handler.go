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

func UserPublishVideoList(c *gin.Context) {
	req := requests.NewDouyinPublishListRequest(c)
	feedVideosRes, _ := video.QueryUserVideoList(req.UserID)
	c.JSON(http.StatusOK, feedVideosRes)
}
