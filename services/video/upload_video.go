package video

import (
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
)

func SavePublishVideo(video *models.Video) (*responses.DouyinPublishActionResponse, error) {

	models.SaveVideo(video)
	return &responses.DouyinPublishActionResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0},
	}, nil

}
