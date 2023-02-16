package chat

import (
	jwtutils "github.com/fyved24/douyin/handlers/user/utils"
	"github.com/fyved24/douyin/responses"
	"github.com/fyved24/douyin/services"
	"github.com/gin-gonic/gin"
)

func UserMessageLog(ctx *gin.Context) {
	token := ctx.Query("token")
	tokenUserID, err := jwtutils.ParseToken(token)
	if err != nil {
		ctx.JSON(500, responses.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "token error",
		})
		return
	}
	userID := tokenUserID.UserID
	targetID := ctx.Query("to_user_id")
	messages, err := services.GetChatLog(userID, targetID)
	if err != nil {
		ctx.JSON(500, responses.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "error",
		})
		return
	}

	ctx.JSON(200, responses.ChatLogResponse{
		CommonResponse: responses.CommonResponse{},
		MessageList:    messages,
	})
}

func CreateMessage(ctx *gin.Context) {

	token := ctx.Query("token")
	tokenUserID, err := jwtutils.ParseToken(token)
	if err != nil {
		ctx.JSON(500, responses.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "token error",
		})
		return
	}
	userID := tokenUserID.UserID
	targetID := ctx.Query("to_user_id")
	content := ctx.Query("content")
	actionType := ctx.Query("action_type")
	err = services.CreateMessage(userID, targetID, content, actionType)
	if err != nil {
		ctx.JSON(500, responses.CommonResponse{
			StatusCode: 1,
			StatusMsg:  "ok",
		})
		return
	}

	ctx.JSON(200, responses.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "ok",
	})
}
