package requests

import "github.com/gin-gonic/gin"

type DouyinPublishListRequest struct {
	UserID uint
	Token  string
	Title  string
}

func NewDouyinPublishListRequest(c *gin.Context) *DouyinPublishListRequest {
	var req DouyinPublishListRequest
	req.check(c)
	return &req
}

func (r *DouyinPublishListRequest) check(c *gin.Context) {

	userID := c.GetString("user_id ")
	r.Title = userID
	token := c.Param("token")
	r.Token = token

}
