package comment

import (
	"net/http"
	"unicode"
	"unicode/utf8"

	"github.com/fyved24/douyin/requests"
	"github.com/fyved24/douyin/responses"
	"github.com/fyved24/douyin/services/comment"
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
	COMMENT_STATUS_COMMENT_CONTENT_ILLEGAL
	COMMENT_STATUS_VIDEO_DONT_EXIST
	COMMENT_STATUS_PARAM_VALID_ERR
)

const (
	STATUS_MSG_SUCCEED                 = "success"
	STATUS_MSG_NOT_LOGINED             = "user not login"
	STATUS_MSG_ILLEGAL_ACTION          = "illegal action"
	STATUS_MSG_COMMENT_CONTENT_ILLEGAL = "illegal comment"
	STATUS_MSG_VIDEO_DONT_EXIST        = "video don't exist"
)

const (
	COMMENT_MAX_LEN      = 100 // 网上查的说抖音的评论区最多100个字的限制
	CODE_POINT_MAX_BYTES = 4
)

// 检查发来的字符串是否符合要求
func validCommentContent(s string) bool {
	// 明显过长过短或是编码不正确
	if len(s) > COMMENT_MAX_LEN*CODE_POINT_MAX_BYTES ||
		len(s) == 0 ||
		!utf8.ValidString(s) {
		return false
	}
	// 转换成Unicode过长
	if cnt := utf8.RuneCountInString(s); cnt > COMMENT_MAX_LEN {
		return false
	}
	// 检查内容是否可打印,是否为全空格类
	blank := true
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
		if !unicode.IsSpace(r) {
			blank = false
		}
	}
	return !blank
}

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
	logined, userID, err := comment.BrowserLogined(&commentActionRequest.Token)
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
	// 检查视频是否真的存在
	if exist, err := comment.VideoExist(commentActionRequest.VideoID); err != nil {
		c.JSON(http.StatusOK, responses.CommentActionResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_VIDEO_DONT_EXIST, StatusMsg: STATUS_MSG_VIDEO_DONT_EXIST},
		})
		return
	} else if !exist {
		c.JSON(http.StatusOK, responses.CommentActionResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_VIDEO_DONT_EXIST, StatusMsg: STATUS_MSG_VIDEO_DONT_EXIST},
		})
		return
	}
	switch commentActionRequest.ActionType {
	case requests.COMMENT_PUBLISH:
		// 检查评论字符串合法性
		if !validCommentContent(commentActionRequest.CommentText) {
			c.JSON(http.StatusOK, responses.CommentActionResponse{
				CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_COMMENT_CONTENT_ILLEGAL, StatusMsg: STATUS_MSG_COMMENT_CONTENT_ILLEGAL},
			})
			return
		}
		// 用户添加评论操作
		respComment, err := comment.AddVideoComment(commentActionRequest.VideoID, userID, commentActionRequest.CommentText)
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
		err := comment.DeleteComment(commentActionRequest.CommentID, userID, commentActionRequest.VideoID)
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
	// 检查视频是否真的存在
	if exist, err := comment.VideoExist(commentListRequest.VideoID); err != nil {
		c.JSON(http.StatusOK, responses.CommentListResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_VIDEO_DONT_EXIST, StatusMsg: STATUS_MSG_VIDEO_DONT_EXIST},
		})
		return
	} else if !exist {
		c.JSON(http.StatusOK, responses.CommentListResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_VIDEO_DONT_EXIST, StatusMsg: STATUS_MSG_VIDEO_DONT_EXIST},
		})
		return
	}
	// 检查浏览用户是否登录
	logined, userID, err := comment.BrowserLogined(&commentListRequest.Token)
	if err != nil {
		c.JSON(http.StatusOK, responses.CommentListResponse{
			CommonResponse: responses.CommonResponse{StatusCode: COMMENT_STATUS_PARSE_JWT_ERR, StatusMsg: err.Error()},
		})
		return
	}
	// 查询视频的所用评论
	videoComments, err := comment.GetVideoComments(commentListRequest.VideoID, logined, userID)
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
