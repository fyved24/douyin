package requests

import (
	"github.com/fyved24/douyin/utils"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
)

type DouyinPublishActionRequest struct {
	Token string
	Filename  string
	Title string
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
	userId := "1"
	file, _ := c.FormFile("data")
	ext := filepath.Ext(file.Filename)
	filename := utils.NewFileName(userId) + ext
	err := c.SaveUploadedFile(file, filename)
	if err != nil {
		return
	}

	r.Filename = filename
	log.Printf("上传文件名: %s\n", filename)

}
