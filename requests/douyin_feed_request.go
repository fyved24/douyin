package requests

import (
	"github.com/fyved24/douyin/responses"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type DouyinFeedRequest struct {
	LatestTime time.Time
	Token      string
	UserID     uint
}

func NewDouyinFeedRequest(c *gin.Context) *DouyinFeedRequest {
	var feedRequest DouyinFeedRequest
	feedRequest.check(c)
	return &feedRequest
}

// 数据校验
func (r *DouyinFeedRequest) check(c *gin.Context) {
	var intTime int64
	var err error
	timestamp := c.Query("latest_time")
	latestTime := time.Now()
	if timestamp != "" {
		intTime, err = strconv.ParseInt(timestamp, 10, 64)
		latestTime = time.Unix(0, intTime*1e6)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.CommonResponse{StatusCode: 1})
	}
	r.LatestTime = latestTime

	token := c.Query("token")
	r.Token = token
	if userid := c.GetString("user_id"); userid != "" {
		userID, _ := strconv.ParseUint(userid, 10, 64)
		r.UserID = uint(userID)
	}
	log.Printf("token: %s, latestTime: %v", token, latestTime)

}
