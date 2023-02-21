package video

import (
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/services/comment"
	videoService "github.com/fyved24/douyin/services/video"
	"github.com/fyved24/douyin/utils"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
	"path/filepath"
)

var (
	videoBucket  = "videos"
	imagesBucket = "images"
)

func PublishVideoAction(c *gin.Context) {
	req := requests.NewDouyinPublishActionRequest(c)
	ctx := c.Copy()
	fileHeader, _ := ctx.FormFile("data")
	//copyContext := c.Copy()

	ext := filepath.Ext(fileHeader.Filename)
	name := utils.NewFileName(req.UserID)
	videoFilename := name + ext
	coverFilename := name + ".jpeg"

	endpointURL := models.MinIOClient.EndpointURL().String()
	// 根据文件所在 Bucket 生成文件链接
	videoFileURL := endpointURL + "/" + videoBucket + "/" + videoFilename
	coverFileURL := endpointURL + "/" + imagesBucket + "/" + coverFilename
	go func() {
		log.Printf("视频上传")
		defer log.Printf("全部上传成功")
		file, _ := fileHeader.Open()
		defer file.Close()
		object, err := models.MinIOClient.PutObject(ctx, videoBucket, videoFilename, file, fileHeader.Size, minio.PutObjectOptions{ContentType: "video/mp4"})
		if err != nil {
			log.Printf("视频上传出错 %v", err)
			return
		}
		log.Printf("视频上传成功 %v", object)

		imageBuff := utils.CutFirstFrameOfVideo(videoFileURL)

		object, err = models.MinIOClient.PutObject(ctx, imagesBucket, coverFilename, imageBuff, int64(imageBuff.Len()), minio.PutObjectOptions{ContentType: "image/jpeg"})
		if err != nil {
			log.Printf("视频封面上传出错 %v", err)
			return
		}
	}()

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

func Upload() {

}
