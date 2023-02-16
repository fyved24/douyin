package responses

import "github.com/fyved24/douyin/models"

type ChatLogResponse struct {
	CommonResponse
	MessageList *[]models.Message `json:"message_list"`
}
