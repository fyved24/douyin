package handlers

import (
	"fmt"
	"github.com/fyved24/douyin/responses"
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

	if err != nil {

		c.JSON(200, responses.FavoriteActionResponse{
			500, err.Error(),
		})
	} else {
		c.JSON(200, responses.FavoriteActionResponse{
			200, "success",
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

	//msg := "success"
	if err != nil {

		c.JSON(200, responses.FavoriteListResponse{
			1, err.Error(), nil,
		})
	} else {
		c.JSON(200, responses.FavoriteListResponse{
			StatusCode: 0, StatusMsg: "success", VideoList: res,
		})
	}

}
