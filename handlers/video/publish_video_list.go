package video

import (
	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/services/video"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserPublishVideoList(c *gin.Context) {
	req := requests.NewDouyinPublishListRequest(c)
	feedVideosRes, _ := video.QueryUserVideoList(req.UserID)
	c.JSON(http.StatusOK, feedVideosRes)
}
