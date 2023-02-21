package router

import (
	"github.com/fyved24/douyin/handlers/chat"
	"github.com/fyved24/douyin/handlers/user"

	"github.com/fyved24/douyin/handlers/comment"
	"github.com/fyved24/douyin/handlers/favorite"
	"github.com/fyved24/douyin/handlers/relation"
	"github.com/fyved24/douyin/handlers/video"
	"github.com/fyved24/douyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {

	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		// feed
		douyinGroup.GET("/feed/", middleware.OptionalTokenJWT(), video.FeedVideoList)

		publishGroup := douyinGroup.Group("publish")
		publishGroup.Use(middleware.JWT())
		{
			publishGroup.POST("/action/", video.PublishVideoAction)
			publishGroup.GET("/list/", video.UserPublishVideoList)
		}

		// relation路由组
		relationGroup := douyinGroup.Group("relation")
		relationGroup.Use(middleware.JWT())
		{
			relationGroup.POST("/action/", relation.RelationAction)
			relationGroup.GET("/follow/list/", relation.FollowList)
			relationGroup.GET("/follower/list/", relation.FollowerList)
			relationGroup.GET("/friend/list/", relation.FriendList)
		}

		commentGroup := douyinGroup.Group("comment")
		{
			commentGroup.GET("/list/", comment.CommentList)
			commentGroup.POST("/action/", comment.CommentAction)
		}

		userGroup := douyinGroup.Group("user")
		{
			userGroup.POST("/register/", user.Register)
			userGroup.POST("/login/", user.Login)
			userGroup.GET("/", user.Info).Use(middleware.JWT())
		}

		favoriteGroup := douyinGroup.Group("favorite")
		favoriteGroup.Use(middleware.JWT())
		{
			favoriteGroup.POST("/action/", middleware.DefaultLimit(), favorite.Favorite)
			favoriteGroup.GET("/list/", favorite.FavoriteList)
		}
		//chat路由组
		chatGroup := douyinGroup.Group("message")
		{
			chatGroup.POST("/action/", chat.CreateMessage)
			chatGroup.GET("/chat/", chat.UserMessageLog)
		}

	}

}
