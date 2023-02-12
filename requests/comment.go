package requests

import (
	"github.com/gin-gonic/gin"
)

type CommentListRequest struct {
	Token   string `json:"token" xml:"token" form:"token" `                            // 用户鉴权token
	VideoID uint   `json:"video_id" xml:"video_id" form:"video_id" binding:"required"` // 视频id
}

type CommentActionType uint

const (
	COMMENT_PUBLISH = 1 + iota
	COMMENT_DELETE
)

type CommentActionRequest struct {
	Token       string            `json:"token" xml:"token" form:"token" binding:"required"`                                      // 用户鉴权token
	VideoID     uint              `json:"video_id" xml:"video_id" form:"video_id" binding:"required"`                             // 视频id
	ActionType  CommentActionType `json:"action_type" xml:"action_type" form:"action_type" binding:"required,oneof=1 2"`          // 1-发布评论，2-删除评论
	CommentText string            `json:"comment_text" xml:"comment_text" form:"comment_text" binding:"required_if=ActionType 1"` // 用户填写的评论内容，在action_type=1的时候使用
	CommentID   uint              `json:"comment_id" xml:"comment_id" form:"comment_id" binding:"required_if=ActionType 2"`       // 要删除的评论id，在action_type=2的时候使用
}

func ReadCommentListRequest(c *gin.Context) (*CommentListRequest, error) {
	var commentListRequest CommentListRequest
	err := commentListRequest.readAndCheck(c)
	if err != nil {
		return nil, err
	}
	return &commentListRequest, nil
}

// 数据校验
func (r *CommentListRequest) readAndCheck(c *gin.Context) error {
	err := c.ShouldBind(r)
	if err != nil {
		return err
	}
	return nil
}

func ReadCommentActionRequest(c *gin.Context) (*CommentActionRequest, error) {
	var commentActionRequest CommentActionRequest
	err := commentActionRequest.readAndCheck(c)
	if err != nil {
		return nil, err
	}
	return &commentActionRequest, nil
}

// 数据校验
func (r *CommentActionRequest) readAndCheck(c *gin.Context) error {
	err := c.ShouldBind(r)
	if err != nil {
		return err
	}
	return nil
}
