package router

import (
	"github.com/fyved24/douyin/handlers/user"
	"github.com/fyved24/douyin/handlers/video"
	"github.com/fyved24/douyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	douyin := app.Group("/douyin")
	douyin.GET("/feed/", video.FeedVideoList)
	douyin.POST("/user/register/", user.Register)
	douyin.POST("/user/login/", user.Login)
	douyin.GET("/user/", user.Info).Use(middleware.JWT())
}
