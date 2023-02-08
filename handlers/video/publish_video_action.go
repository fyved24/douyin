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

var (
	LocalStorage     = "local_storage/"
	FileServerPrefix = "http://192.168.31.80:8080/file/"
)

func PublishVideoAction(c *gin.Context) {
	req := requests.NewDouyinPublishActionRequest(c)

	file, _ := c.FormFile("data")
	ext := filepath.Ext(file.Filename)
	filename := utils.NewFileName(req.UserID) + ext
	videoFilePath := LocalStorage + filename
	videoFileURL := FileServerPrefix + filename
	err := c.SaveUploadedFile(file, videoFilePath)
	if err != nil {
		return
	}
	video := &models.Video{
		AuthorID: req.UserID,
		PlayUrl:  videoFileURL,
		Title:    req.Title,
	}
	publishVideosRes, _ := videoService.SavePublishVideo(video)
	c.JSON(http.StatusOK, publishVideosRes)
}
