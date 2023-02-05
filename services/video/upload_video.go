package video

import (
	"github.com/fyved24/douyin/responses"
)

func PublishVideo(userID, filename string) (*responses.DouyinPublishActionResponse, error) {
	return &responses.DouyinPublishActionResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0},
	}, nil

}
