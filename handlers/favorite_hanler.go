package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// Favorite 点赞视频方法
func Favorite(c *gin.Context) {

	// 1. token 验证
	token, _ := c.Get("token")

	// 2. 获得
	user_id := "123"

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseInt(actionTypeStr, 10, 64)

	err := se

}
