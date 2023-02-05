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
}

func NewDouyinFeedRequest(c *gin.Context) *DouyinFeedRequest {
	var feedRequest DouyinFeedRequest
	feedRequest.check(c)
	return &feedRequest
}

// 数据校验
func (r *DouyinFeedRequest) check(c *gin.Context) {
	timestamp := c.Query("latest_time")
	token := c.Query("token")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.CommonResponse{StatusCode: 1})
	}
	latestTime = time.Unix(0, intTime*1e6)
	log.Printf("token: %s, latestTime: %v", token, latestTime )
	r.LatestTime = latestTime
	r.Token = token

}
