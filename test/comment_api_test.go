package test

import (
	"net/http"
	"testing"

	"github.com/fyved24/douyin/handlers/comment"
	"github.com/stretchr/testify/assert"
	// "github.com/fyved24/douyin/models"
)

func TestCommentHandler(t *testing.T) {
	e := newExpect(t)

	const (
		videoId   = 1
		commentId = 1
	)

	// 添加一个评论的请求
	addCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", comment.DemoToken).WithQuery("video_id", videoId).WithQuery("action_type", 1).WithQuery("comment_text", "测试评论").
		WithFormField("token", comment.DemoToken).WithFormField("video_id", videoId).WithFormField("action_type", 1).WithFormField("comment_text", "测试评论").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	addCommentResp.Value("status_code").Number().IsEqual(0)
	addCommentResp.Value("comment").Object().Value("id").Number().Gt(0)
	addCommentResp.Value("comment").Object().Value("id").Number().IsEqual(commentId)

	// 查询视频评论
	commentListResp := e.GET("/douyin/comment/list/").
		WithQuery("token", comment.DemoToken).WithQuery("video_id", videoId).
		WithFormField("token", comment.DemoToken).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListResp.Value("status_code").Number().IsEqual(0)
	containTestComment := false
	for _, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		comment.ContainsKey("id")
		// t.Logf("the response id is:%#v", comment.Value("id"))
		comment.ContainsKey("user")
		comment.Value("content").String().NotEmpty()
		comment.Value("create_date").String().NotEmpty()
		if int(comment.Value("id").Number().Raw()) == commentId {
			containTestComment = true
		}
	}

	assert.True(t, containTestComment, "Can't find test comment in list")

	// 删除一个评论
	delCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", comment.DemoToken).WithQuery("video_id", videoId).WithQuery("action_type", 2).WithQuery("comment_id", commentId).
		WithFormField("token", comment.DemoToken).WithFormField("video_id", videoId).WithFormField("action_type", 2).WithFormField("comment_id", commentId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	delCommentResp.Value("status_code").Number().Equal(0)

	// 再次查看,看是否已经删除
	commentListResp = e.GET("/douyin/comment/list/").
		WithQuery("token", comment.DemoToken).WithQuery("video_id", videoId).
		WithFormField("token", comment.DemoToken).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListResp.Value("status_code").Number().IsEqual(0)

	containTestComment = false
	for _, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		comment.ContainsKey("id")
		// t.Logf("the response id is:%#v", comment.Value("id"))
		comment.ContainsKey("user")
		comment.Value("content").String().NotEmpty()
		comment.Value("create_date").String().NotEmpty()
		if int(comment.Value("id").Number().Raw()) == commentId {
			containTestComment = true
		}
	}

	assert.False(t, containTestComment, "test comment in list not deleted")
}

func TestComment(t *testing.T) {
	e := newExpect(t)

	feedResp := e.GET("/douyin/feed/").Expect().Status(http.StatusOK).JSON().Object()
	feedResp.Value("status_code").Number().Equal(0)
	feedResp.Value("video_list").Array().Length().Gt(0)
	firstVideo := feedResp.Value("video_list").Array().First().Object()
	videoId := firstVideo.Value("id").Number().Raw()

	_, token := getTestUserToken(testUserA, e)

	addCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 1).WithQuery("comment_text", "测试评论").
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 1).WithFormField("comment_text", "测试评论").
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	addCommentResp.Value("status_code").Number().Equal(0)
	addCommentResp.Value("comment").Object().Value("id").Number().Gt(0)
	commentId := int(addCommentResp.Value("comment").Object().Value("id").Number().Raw())

	commentListResp := e.GET("/douyin/comment/list/").
		WithQuery("token", token).WithQuery("video_id", videoId).
		WithFormField("token", token).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListResp.Value("status_code").Number().Equal(0)
	containTestComment := false
	for _, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		comment.ContainsKey("id")
		comment.ContainsKey("user")
		comment.Value("content").String().NotEmpty()
		comment.Value("create_date").String().NotEmpty()
		if int(comment.Value("id").Number().Raw()) == commentId {
			containTestComment = true
		}
	}

	assert.True(t, containTestComment, "Can't find test comment in list")

	delCommentResp := e.POST("/douyin/comment/action/").
		WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", 2).WithQuery("comment_id", commentId).
		WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", 2).WithFormField("comment_id", commentId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	delCommentResp.Value("status_code").Number().Equal(0)
}
