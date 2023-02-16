package requests

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type DouyinPublishListRequest struct {
	UserID uint
	Token  string
}

func NewDouyinPublishListRequest(c *gin.Context) *DouyinPublishListRequest {
	var req DouyinPublishListRequest
	req.check(c)
	return &req
}

func (r *DouyinPublishListRequest) check(c *gin.Context) {

	userID, _ := strconv.ParseUint(c.GetString("user_id"), 10, 64)
	r.UserID = uint(userID)
	token := c.Param("token")
	r.Token = token

}
