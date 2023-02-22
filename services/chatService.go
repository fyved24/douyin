package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fyved24/douyin/models"
	"github.com/redis/go-redis/v9"
	"github.com/willf/bloom"
	"strconv"
	"time"

	"log"
)

//获取聊天记录
func GetChatLogWithCache(userID int, targetID int) (*[]models.Message, error) {
	//redisClient := middleware.NewRedisClient("47.93.10.203:6379", "zkrt", 2)
	redisClient := models.RedisDB
	//defer redisClient.Close()
	i := strconv.Itoa(userID)
	j := strconv.Itoa(targetID)
	if i == "" || j == "" {
		return &[]models.Message{}, nil
	}
	//messages, err := models.GetMessageByID(userID, targetID)
	//if err != nil {
	//	return nil, err
	//}
	//return messages, nil
	messages, err := getChatLogFromCache(redisClient, userID, targetID)
	if err != nil {
		//if err == redis.Nil {
		//	//未找到key
		//	return &[]models.Message{}, nil
		//}
		messages, err = models.GetMessageByID(userID, targetID)
		if err != nil {
			return &[]models.Message{}, err
		}
		//err = setChatLogToCache(redisClient, userID, targetID, *messages)
		//if err != nil {
		//	log.Printf("Failed to set chat log to cache: %s", err)
		//}
	}

	if len(*messages) == 0 {
		return nil, nil
	}
	return messages, nil
}

//发送消息存储
func CreateMessage(userID int, targetID int, content string, actionType string) error {
	redisClient := models.RedisDB
	i := strconv.Itoa(userID)
	j := strconv.Itoa(targetID)
	if i == "" || j == "" {
		return errors.New("没有消息")
	}
	//defer redisClient.Close()
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

func getChatLogFromCache(redisClient *redis.Client, userID int, targetID int) (*[]models.Message, error) {
	//messageID, err := redisClient.Get(context.Background(), fmt.Sprintf("chat:%s:%s:messageID", userID, targetID)).Result()
	//if err != nil {
	//	return nil, err
	//}
	var bloomFilter *bloom.BloomFilter
	bloomFilter = bloom.New(100000, 5)
	// 查询数据库
	messageID, err := redisClient.Get(context.Background(), fmt.Sprintf("chat:%d:%d:messageID", userID, targetID)).Result()
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
	key := fmt.Sprintf("chat:%d:%d:messageID", userID, targetID)
	if !bloomFilter.TestString(key) {
		return nil, nil
	}
	return &messages, nil
}

func setChatLogToCache(redisClient *redis.Client, userID int, targetID int, messages []models.Message) error {
	messageID := messages[0].ID
	key := fmt.Sprintf("chat:%s:%s:%s:", userID, targetID, messageID)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	if err := redisClient.HSet(context.Background(), key, data, 0).Err(); err != nil {
		return err
	}
	if err := redisClient.Expire(context.Background(), key, time.Hour).Err(); err != nil {
		return err
	}
	return nil
}
