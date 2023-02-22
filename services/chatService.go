package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fyved24/douyin/models"
	"github.com/redis/go-redis/v9"
	"github.com/willf/bloom"
	"log"
	"strconv"
	"time"
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
	fmt.Print(messages)
	if messages == nil {
		//messages, err = models.GetMessageByID(userID, targetID)
		//if messages == nil {
		//	return &[]models.Message{}, err
		//}
		//err = setChatLogToCache(redisClient, userID, targetID, *messages)
		//if err != nil {
		//	log.Printf("Failed to set chat log to cache: %s", err)
		//}
		return nil, err
	}

	//if len(*messages) == 0 {
	//	return nil, nil
	//}
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
	err = setChatLogToCache(redisClient, userID, targetID, *m)
	if err != nil {
		log.Printf("Failed to set chat log to cache: %s", err)
	}
	return nil
}

func getChatLogFromCache(redisClient *redis.Client, userID int, targetID int) (*[]models.Message, error) {

	var bloomFilter *bloom.BloomFilter
	bloomFilter = bloom.New(100000, 5)
	key := fmt.Sprintf("chat:%d:%d", userID, targetID)
	//var messages []models.Message
	result, err := redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]models.Message, len(result))
	i := 0
	for _, value := range result {
		var message models.Message
		//fmt.Print(value)
		err = json.Unmarshal([]byte(value), &message)
		if err != nil {
			return nil, err
		}
		//mID, _ := strconv.Atoi(messageID)
		//message.ID = uint(mID)
		messages[i] = message
		i++
	}

	if messages == nil && !bloomFilter.TestString(key) {
		return nil, nil
	}
	fmt.Print(messages)
	return &messages, nil
}

func setChatLogToCache(redisClient *redis.Client, userID int, targetID int, messages models.Message) error {
	messageID := strconv.Itoa(int(messages.ID))
	var key string
	key = fmt.Sprintf("chat:%d:%d", userID, targetID)
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	if err := redisClient.HSet(context.Background(), key, messageID, data).Err(); err != nil {
		return err
	}
	if err := redisClient.Expire(context.Background(), key, time.Hour).Err(); err != nil {
		return err
	}
	return nil
}
