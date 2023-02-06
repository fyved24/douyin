package requests

import (
	"github.com/gin-gonic/gin"
)

type DouyinPublishActionRequest struct {
	UserID uint
	Token  string
	Title  string
}

func NewDouyinPublishActionRequest(c *gin.Context) *DouyinPublishActionRequest {
	var req DouyinPublishActionRequest
	req.check(c)
	return &req
}

func (r *DouyinPublishActionRequest) check(c *gin.Context) {

	title := c.PostForm("title")
	r.Title = title
	token := c.PostForm("token")
	r.Token = token
	r.UserID = 13

}
