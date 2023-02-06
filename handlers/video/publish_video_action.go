package video

import (
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/requests"
	videoService "github.com/fyved24/douyin/services/video"
	"github.com/fyved24/douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func PublishVideoHandler(c *gin.Context) {
	req := requests.NewDouyinPublishActionRequest(c)

	file, _ := c.FormFile("data")
	ext := filepath.Ext(file.Filename)
	videoFilePath := utils.NewFileName(req.UserID) + ext
	err := c.SaveUploadedFile(file, videoFilePath)
	if err != nil {
		return
	}
	video := &models.Video{
		AuthorID: req.UserID,
		PlayUrl:  videoFilePath,
		Title:    req.Title,
	}
	publishVideosRes, _ := videoService.SavePublishVideo(video)
	c.JSON(http.StatusOK, publishVideosRes)
}
