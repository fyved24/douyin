package comment

import (
	"net/http"

	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/responses"
	"github.com/fyved24/douyin/services"
	"github.com/gin-gonic/gin"
)

const (
	COMMENT_STATUS_CODE_SUCCESS int32 = iota
	COMMENT_STATUS_PARSE_ACTION_REQUEST_ERR
	COMMENT_STATUS_PARSE_LIST_REQUEST_ERR
	COMMENT_STATUS_PARSE_JWT_ERR
	COMMENT_STATUS_ACTION_NOT_LOGINED
	COMMENT_STATUS_PUBLISH_FAILED
	COMMENT_STATUS_DELETE_FAILED
	COMMENT_STATUS_ILLEGAL_ACTION
	COMMENT_STATUS_GET_VIDEO_COMM_ERR
)

const (
	STATUS_MSG_SUCCEED        = "success"
	STATUS_MSG_NOT_LOGINED    = "user not login"
	STATUS_MSG_ILLEGAL_ACTION = "illegal action"
)

// 评论操作的controller
func CommentAction(c *gin.Context) {
	// 读取参数
	commentActionRequest, err := requests.ReadCommentActionRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommentActionResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_PARSE_ACTION_REQUEST_ERR, StatusMsg: err.Error()},
		})
		return
	}
	// 评论操作都需要用户已登录
	logined, userID, err := services.BrowserLogined(&commentActionRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommentActionResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_PARSE_JWT_ERR, StatusMsg: err.Error()},
		})
		return
	}
	if !logined {
		c.JSON(http.StatusOK, responses.CommentActionResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_ACTION_NOT_LOGINED, StatusMsg: STATUS_MSG_NOT_LOGINED},
		})
		return
	}
	switch commentActionRequest.ActionType {
	case requests.COMMENT_PUBLISH:
		// 用户添加评论操作
		respComment, err := services.AddVideoComment(commentActionRequest.VideoID, userID, commentActionRequest.CommentText)
		if err != nil {
			c.JSON(http.StatusOK, responses.CommentActionResponse{
				CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_PUBLISH_FAILED, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, responses.CommentActionResponse{CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_CODE_SUCCESS, StatusMsg: STATUS_MSG_SUCCEED},
			Comment: *respComment})
	case requests.COMMENT_DELETE:
		// 用户删除评论操作
		err := services.DeleteComment(commentActionRequest.CommentID, userID, commentActionRequest.VideoID)
		if err != nil {
			c.JSON(http.StatusOK, responses.CommentActionResponse{
				CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_DELETE_FAILED, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, responses.CommentActionResponse{CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_CODE_SUCCESS, StatusMsg: STATUS_MSG_SUCCEED}})
	default:
		// 非法操作类型
		c.JSON(http.StatusOK, responses.CommentActionResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_ILLEGAL_ACTION, StatusMsg: STATUS_MSG_ILLEGAL_ACTION},
		})
	}

}

// 评论列表的controller
func CommentList(c *gin.Context) {
	// 读取参数
	commentListRequest, err := requests.ReadCommentListRequest(c)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommentListResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_PARSE_LIST_REQUEST_ERR, StatusMsg: err.Error()},
		})
		return
	}
	// 检查浏览用户是否登录
	logined, userID, err := services.BrowserLogined(&commentListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommentListResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_PARSE_JWT_ERR, StatusMsg: err.Error()},
		})
		return
	}
	// 查询视频的所用评论
	videoComments, err := services.GetVideoComments(commentListRequest.VideoID, logined, userID)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommentListResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_GET_VIDEO_COMM_ERR, StatusMsg: err.Error()},
		})
		return
	}
	// 返回评论
	c.JSON(http.StatusOK, responses.CommentListResponse{
		CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_CODE_SUCCESS, StatusMsg: STATUS_MSG_SUCCEED},
		CommentList:    videoComments,
	})
}
