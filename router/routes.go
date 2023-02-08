package router

import (
	comment "github.com/fyved24/douyin/handlers/comment"
	relation "github.com/fyved24/douyin/handlers/relation"
	video "github.com/fyved24/douyin/handlers/video"
	"github.com/fyved24/douyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// feed
		douyinGroup.GET("/feed/", video.FeedVideoList)

		publishGroup := douyinGroup.Group("publish")
		{
			publishGroup.GET("/action/", video.PublishVideoAction)
			publishGroup.GET("/list/", video.UserPublishVideoList)
		}
		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action/", middleware.JwtMiddleware(), relation.RelationAction)
			relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), relation.FollowList)
			relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), relation.FollowerList)
		}

		commentGroup := douyinGroup.Group("comment")
		{
			commentGroup.GET("/list/", comment.CommentList)
			commentGroup.GET("/action/", comment.CommentAction)
		}

	}

	// 文件服务
	r.GET("/file/:filename", video.FileServer)

}
