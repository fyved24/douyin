package router

import (
	"github.com/fyved24/douyin/handlers"
	"github.com/fyved24/douyin/handlers/video"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *gin.Engine) {
	douyin := app.Group("/douyin")
	douyin.GET("/feed/", video.FeedVideoList)

	douyin.POST("/favorite/action/", handlers.Favorite)
	douyin.GET("/favorite/list/", handlers.FavoriteList)
}
