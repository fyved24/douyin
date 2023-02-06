package comment

import (
	"net/http"

	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/responses"
	"github.com/gin-gonic/gin"
)

var DemoUser = responses.User{
	ID:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      false,
}

var DemoComments = []responses.Comment{
	{
		ID:         1,
		User:       DemoUser,
		Content:    "Test Comment",
		CreateDate: "05-01",
	},
}

var DemoToken = "zhangleidouyin"

var UsersLoginInfo = map[string]responses.User{
	DemoToken: {
		ID:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var commentMap = map[uint][]responses.Comment{}

// 评论操作的controller
func CommentAction(c *gin.Context) {

	commentActionRequest, err := requests.ReadCommentActionRequest(c)

	if err != nil {
		c.JSON(http.StatusOK, responses.CommonResponse{StatusCode: -1, StatusMsg: err.Error()})
		return
	}

	if user, exist := UsersLoginInfo[commentActionRequest.Token]; exist {
		switch commentActionRequest.ActionType {
		case requests.COMMENT_PUBLISH:
			text := commentActionRequest.CommentText
			var commentResponse = responses.Comment{
				ID:         1,
				User:       user,
				Content:    text,
				CreateDate: "05-01",
			}
			commentMap[commentActionRequest.VideoID] = append(commentMap[commentActionRequest.VideoID], commentResponse)
			c.JSON(http.StatusOK, responses.CommentActionResponse{CommonResponse: responses.CommonResponse{StatusCode: 0},
				Comment: commentResponse})
			return

		case requests.COMMENT_DELETE:
			var delIdx int = -1
			for idx, ele := range commentMap[commentActionRequest.VideoID] {
				if ele.ID == int64(commentActionRequest.CommentID) {
					delIdx = idx
					break
				}
			}
			if delIdx >= 0 {
				commentMap[commentActionRequest.VideoID] = append(commentMap[commentActionRequest.VideoID][:delIdx], commentMap[commentActionRequest.VideoID][delIdx+1:]...)
				c.JSON(http.StatusOK, responses.CommentActionResponse{CommonResponse: responses.CommonResponse{StatusCode: 0, StatusMsg: "success"}})
			} else {
				c.JSON(http.StatusOK, responses.CommentActionResponse{CommonResponse: responses.CommonResponse{StatusCode: 1, StatusMsg: "no such comment"}})
			}

		}
	} else {
		c.JSON(http.StatusOK, responses.CommonResponse{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// 评论列表的controller
func CommentList(c *gin.Context) {
	commentListRequest, err := requests.ReadCommentListRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommonResponse{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	var videoComments []responses.Comment = commentMap[commentListRequest.VideoID]
	c.JSON(http.StatusOK, responses.CommentListResponse{
		CommonResponse: responses.CommonResponse{StatusCode: 0, StatusMsg: "success"},
		CommentList:    videoComments,
	})
}
