package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fyved24/douyin/middleware"
	"github.com/fyved24/douyin/models"
	"github.com/go-redis/redis/v8"
	"log"
)

//获取聊天记录
func GetChatLogWithCache(userID string, targetID string) (*[]models.Message, error) {
	redisClient := middleware.NewRedisClient("47.93.10.203:6379", "zkrt", 2)
	defer redisClient.Close()

	//messages, err := models.GetMessageByID(userID, targetID)
	//if err != nil {
	//	return nil, err
	//}
	//return messages, nil
	messages, err := getChatLogFromCache(redisClient, userID, targetID)
	if err != nil {
		messages, err = models.GetMessageByID(userID, targetID)
		if err != nil {
			return nil, err
		}
		err = setChatLogToCache(redisClient, userID, targetID, *messages)
		if err != nil {
			log.Printf("Failed to set chat log to cache: %s", err)
		}
	}
	return messages, nil
}

//发送消息存储
func CreateMessage(userID string, targetID string, content string, actionType string) error {
	redisClient := middleware.NewRedisClient("47.93.10.203:6379", "zkrt", 2)
	defer redisClient.Close()
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
	err = setChatLogToCache(redisClient, userID, targetID, []models.Message{*m})
	if err != nil {
		log.Printf("Failed to set chat log to cache: %s", err)
	}
	return nil
}

func getChatLogFromCache(redisClient *redis.Client, userID string, targetID string) (*[]models.Message, error) {
	messageID, err := redisClient.Get(context.Background(), fmt.Sprintf("chat:%s:%s:messageID", userID, targetID)).Result()
	if err != nil {
		return nil, err
	}

	data, err := redisClient.HGetAll(context.Background(), fmt.Sprintf("chat:%s:%s:%s", userID, targetID, messageID)).Result()
	if err != nil {
		return nil, err
	}
	var messages []models.Message
	for _, messageData := range data {
		var message models.Message
		err = json.Unmarshal([]byte(messageData), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return &messages, nil
}

func setChatLogToCache(redisClient *redis.Client, userID string, targetID string, messages []models.Message) error {
	messageID := messages[0].ID
	key := fmt.Sprintf("chat:%s:%s:%s:", userID, targetID, messageID)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return redisClient.HSet(context.Background(), key, data, 0).Err()
}
