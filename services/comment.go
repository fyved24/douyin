package services

import (
	"github.com/fyved24/douyin/responses"
)

func getVideoComments(videoID uint, limit, offset int, logined bool, userID uint) (res []responses.Comment, err error) {

	return
}

func GetVideoComments(videoID uint, logined bool, userID uint) (res []responses.Comment, err error) {
	res, err = getVideoComments(videoID, -1, -1, logined, userID)
	return
}
