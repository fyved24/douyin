package video

import (
	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/services/video"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func PublishVideo(c *gin.Context) {
	log.Printf("PublishVideo")
	req := requests.NewDouyinPublishActionRequest(c)
	userId := "1"
	publishVideosRes, _ := video.PublishVideo(userId, req.Filename)
	c.JSON(http.StatusOK, publishVideosRes)
}
