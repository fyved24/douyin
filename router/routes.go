package router

import (
	"github.com/fyved24/douyin/handlers"
	video2 "github.com/fyved24/douyin/handlers/video"
	"github.com/fyved24/douyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// feed
		douyinGroup.GET("/feed/", handlers.FeedVideoList)

		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action/", middleware.JwtMiddleware(), handlers.RelationAction)
			relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), handlers.FollowList)
			relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), handlers.FollowerList)
		}
	} // 文件服务
	r.GET("/file/:filename", video2.FileServer)

}
