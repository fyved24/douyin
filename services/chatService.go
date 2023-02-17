package services

import "github.com/fyved24/douyin/models"

//获取聊天记录
func GetChatLog(userID string, targetID string) (*[]models.Message, error) {
	messages, err := models.GetMessageByID(userID, targetID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

//发送消息存储
func CreateMessage(userID string, targetID string, content string, actionType string) error {
	m := &models.Message{
		UserID:     userID,
		TargetID:   targetID,
		Content:    content,
		ActionType: actionType,
	}
	err := models.ChatMessageCreat(m)
	if err != nil {
		return err
	}
	return nil
}
