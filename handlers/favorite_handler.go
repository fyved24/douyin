package handlers

import (
	"fmt"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Favorite 点赞视频方法
func Favorite(c *gin.Context) {

	// 1. token 验证
	token, _ := c.Get("token")
	fmt.Println(token)
	// 2. 获得
	//userId := utils.GetUserIDByToken(token)
	userId, _ := strconv.ParseInt("1", 10, 64)

	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseInt(videoIdStr, 10, 64)

	actionTypeStr := c.Query("action_type")
	actionType, _ := strconv.ParseInt(actionTypeStr, 10, 64)

	err := services.FavoriteAction(userId, videoId, actionType)

	msg := "success"

	if err != nil {
		msg = err.Error()
		c.JSON(200, models.FavoriteActionResponse{
			500, &msg,
		})
	} else {
		c.JSON(200, models.FavoriteActionResponse{
			0, &msg,
		})
	}

}

func FavoriteList(c *gin.Context) {
	// 1. token 验证
	token, _ := c.Get("token")
	fmt.Println(token)
	// 2. 获得
	//userId := utils.GetUserIDByToken(token)
	userId, _ := strconv.ParseInt("1", 10, 64)

	res, err := services.FindAllFavorite(userId)

	msg := "success"
	if err != nil {
		msg = "查询失败"
		c.JSON(200, models.FavoriteListResponse{
			1, &msg, nil,
		})
	} else {
		c.JSON(200, models.FavoriteListResponse{
			StatusCode: 0, StatusMsg: &msg, VideoList: res,
		})
	}

}
