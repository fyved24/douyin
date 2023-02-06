package router

import (
	"github.com/fyved24/douyin/handlers/video"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	douyin := app.Group("/douyin")
	douyin.GET("/feed/", video.FeedVideoList)
	douyin.POST("/publish/action/", video.PublishVideoHandler)

	// 文件服务
	app.GET("/file/:filename", video.FileServer)
}
