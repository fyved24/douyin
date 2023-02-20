package video

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/responses"
	"github.com/fyved24/douyin/services/comment"
	videoService "github.com/fyved24/douyin/services/video"
	"github.com/fyved24/douyin/utils"
	"github.com/gin-gonic/gin"
)

var (
	LocalStorage     = "local_storage/"
	FileServerPrefix = "http://192.168.31.80:8080/file/"
)

func PublishVideoAction(c *gin.Context) {
	req := requests.NewDouyinPublishActionRequest(c)
	file, _ := c.FormFile("data")
	ext := filepath.Ext(file.Filename)
	name := utils.NewFileName(req.UserID)
	videoFilename := name + ext
	coverFilename := name + ".jpg"

	videoFilePath := LocalStorage + videoFilename
	videoFileURL := FileServerPrefix + videoFilename

	coverFilePath := LocalStorage + coverFilename
	coverFileURL := FileServerPrefix + coverFilename
	err := c.SaveUploadedFile(file, videoFilePath)
	err = utils.CutFirstFrameOfVideo(coverFilePath, videoFilePath)

	if err != nil {
		log.Printf("err %v", err)
		c.JSON(http.StatusInternalServerError, responses.CommonResponse{
			StatusCode: 1,
			StatusMsg:  "获取封面失败",
		})
		return
	}
	video := &models.Video{
		AuthorID: req.UserID,
		PlayUrl:  videoFileURL,
		CoverUrl: coverFileURL,
		Title:    req.Title,
	}
	publishVideosRes, _ := videoService.SavePublishVideo(video)
	c.JSON(http.StatusOK, publishVideosRes)
	// 发布视频成功时更新本地缓存中用户信息的作品数
	comment.ChangeUserCacheWorkCount(req.UserID)
}
