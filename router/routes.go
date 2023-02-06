package router

import (
	"github.com/fyved24/douyin/handlers/comment"
	"github.com/fyved24/douyin/handlers/video"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	douyin := app.Group("/douyin")
	douyin.GET("/feed/", video.FeedVideoList)
	douyin.GET("/comment/list/", comment.CommentList)
	douyin.POST("/comment/action/", comment.CommentAction)
}
