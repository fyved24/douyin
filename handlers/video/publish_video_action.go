package video

import (
	"context"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/responses"
	"github.com/fyved24/douyin/services/comment"
	videoService "github.com/fyved24/douyin/services/video"
	"github.com/fyved24/douyin/utils"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var (
	videoBucket  = "videos"
	imagesBucket = "images"
)

func PublishVideoAction(c *gin.Context) {
	req := requests.NewDouyinPublishActionRequest(c)
	fileHeader, _ := c.FormFile("data")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	file, _ := fileHeader.Open()
	defer file.Close()
	ext := filepath.Ext(fileHeader.Filename)
	name := utils.NewFileName(req.UserID)
	videoFilename := name + ext
	coverFilename := name + ".jpeg"
	object, err := models.MinIOClient.PutObject(ctx, videoBucket, videoFilename, file, fileHeader.Size, minio.PutObjectOptions{ContentType: "video/mp4"})
	if err != nil {
		return
	}
	log.Printf("视频上传成功 %v", object)

	endpointURL := models.MinIOClient.EndpointURL().String()
	videoFileURL := endpointURL + "/" + videoBucket + "/" + videoFilename

	imageBuff := utils.CutFirstFrameOfVideo(videoFileURL)

	object, err = models.MinIOClient.PutObject(ctx, imagesBucket, coverFilename, imageBuff, int64(imageBuff.Len()), minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		return
	}
	log.Printf("视频上传成功 %v", object)

	coverFileURL := endpointURL + "/" + imagesBucket + "/" + coverFilename

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
