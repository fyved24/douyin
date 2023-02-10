package test

import (
	"math/rand"
	"testing"
	"time"

	"net/http"

	"github.com/fyved24/douyin/models"
	"github.com/stretchr/testify/assert"
)

var users []models.User
var videos []models.Video
var following map[[2]uint]struct{}

func init() {
	models.InitDB()
	users, videos = makeSomeUsersAndVideos()
	following = makeSomeFollows(users)
}

const (
	COMMENT_ACTION_ADD = 1 + iota
	COMMENT_ACTION_DEL
)

const (
	GEN_COMMENT_COUNT = 1000
)

type commenterTestStruct struct {
	CommentID         uint
	UserID            uint
	UserName          string
	UserFollowCount   int
	UserFollowerCount int
}

func TestCommentHandler(t *testing.T) {
	e := newExpect(t)
	rand.Seed(time.Now().UnixNano())
	videoIdx := rand.Intn(len(videos))
	videoId := videos[videoIdx].ID
	var commentsAdded = make([]commenterTestStruct, 0, GEN_COMMENT_COUNT)
	for i := 0; i < GEN_COMMENT_COUNT; i++ {
		userIdx := rand.Intn(len(users))
		userID := users[userIdx].ID
		userName := users[userIdx].Name
		flwCNt, flwrCnt := users[userIdx].FollowCount, users[userIdx].FollowerCount
		// token := getTestUserToken(userID, true, false)
		token := getTestUserTokenNow(userName, userID, true, false)
		// 添加一个评论的请求
		addCommentResp := e.POST("/douyin/comment/action/").
			WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", COMMENT_ACTION_ADD).WithQuery("comment_text", "测试评论").
			WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", COMMENT_ACTION_ADD).WithFormField("comment_text", "测试评论").
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		addCommentResp.Value("status_code").Number().IsEqual(0)
		addCommentResp.Value("status_msg").String().IsEqual("success")
		addCommentResp.Value("comment").Object().Value("id").Number().Gt(0)
		commentID := uint(addCommentResp.Value("comment").Object().Value("id").Number().Raw())
		addCommentResp.Value("comment").Object().Value("user").Object().Value("id").Number().IsEqual(userID)
		if flwCNt > 0 {
			addCommentResp.Value("comment").Object().Value("user").Object().Value("follow_count").Number().IsEqual(flwCNt)
		}
		if flwrCnt > 0 {
			addCommentResp.Value("comment").Object().Value("user").Object().Value("follower_count").Number().IsEqual(flwrCnt)
		}
		addCommentResp.Value("comment").Object().Value("user").Object().Value("name").String().IsEqual(userName)
		commentsAdded = append(commentsAdded, commenterTestStruct{commentID, userID, userName, int(flwCNt), int(flwrCnt)})
		addCommentResp.Value("comment").Object().Value("content").String().IsEqual("测试评论")
		addCommentResp.Value("comment").Object().Value("create_date").String().IsEqual(time.Now().Format("01-02"))
	}

	// token := getTestUserToken(0, false, false)
	token := getTestUserTokenNow("", 0, false, false)
	// 查询视频评论未登录
	commentListResp := e.GET("/douyin/comment/list/").
		WithQuery("token", token).WithQuery("video_id", videoId).
		WithFormField("token", token).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListResp.Value("status_code").Number().IsEqual(0)
	commentListResp.Value("comment_list").Array().Length().IsEqual(GEN_COMMENT_COUNT)
	containTestComment := false
	for idx, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		cmm := commentsAdded[GEN_COMMENT_COUNT-idx-1]
		comment.Value("id").Number().IsEqual(cmm.CommentID)
		comment.Value("user").Object().Value("id").Number().IsEqual(cmm.UserID)
		comment.Value("user").Object().Value("name").String().IsEqual(cmm.UserName)
		if cmm.UserFollowCount > 0 {
			comment.Value("user").Object().Value("follow_count").Number().IsEqual(cmm.UserFollowCount)
		} else {
			comment.Value("user").Object().NotContainsKey("follow_count")
		}
		if cmm.UserFollowerCount > 0 {
			comment.Value("user").Object().Value("follower_count").Number().IsEqual(cmm.UserFollowerCount)
		} else {
			comment.Value("user").Object().NotContainsKey("follower_count")
		}
		comment.Value("content").String().NotEmpty().IsEqual("测试评论")
		comment.Value("create_date").String().NotEmpty().IsEqual(time.Now().Format("01-02"))
		containTestComment = true
	}

	assert.True(t, containTestComment, "Can't find test comment in list")

	// 查询视频评论已登录
	loginedUserIdx := rand.Intn(len(users))
	loginedUserID := users[loginedUserIdx].ID
	loginedUserName := users[loginedUserIdx].Name
	// loginedToken := getTestUserToken(loginedUserID, true, false)
	loginedToken := getTestUserTokenNow(loginedUserName, loginedUserID, true, false)
	commentListFavResp := e.GET("/douyin/comment/list/").
		WithQuery("token", loginedToken).WithQuery("video_id", videoId).
		WithFormField("token", loginedToken).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListFavResp.Value("status_code").Number().IsEqual(0)
	commentListFavResp.Value("comment_list").Array().Length().IsEqual(GEN_COMMENT_COUNT)
	containTestComment = false
	for idx, element := range commentListFavResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		cmm := commentsAdded[GEN_COMMENT_COUNT-idx-1]
		comment.Value("id").Number().IsEqual(cmm.CommentID)
		comment.Value("user").Object().Value("id").Number().IsEqual(cmm.UserID)
		comment.Value("user").Object().Value("name").String().IsEqual(cmm.UserName)
		if cmm.UserFollowCount > 0 {
			comment.Value("user").Object().Value("follow_count").Number().IsEqual(cmm.UserFollowCount)
		}
		if cmm.UserFollowerCount > 0 {
			comment.Value("user").Object().Value("follower_count").Number().IsEqual(cmm.UserFollowerCount)
		}
		key := [2]uint{loginedUserID, cmm.UserID}
		if _, has := following[key]; has {
			comment.Value("user").Object().Value("is_follow").Boolean().True()
		} else {
			comment.Value("user").Object().NotContainsKey("is_follow")
		}
		comment.Value("content").String().NotEmpty().IsEqual("测试评论")
		comment.Value("create_date").String().NotEmpty().IsEqual(time.Now().Format("01-02"))
		containTestComment = true
	}
	assert.True(t, containTestComment, "Can't find test comment in list")

	// 测试删除评论
	for idx, cmm := range commentsAdded {
		curUserID := cmm.UserID
		curUserName := cmm.UserName
		// curUsrToken := getTestUserToken(curUserID, true, false)
		curUsrToken := getTestUserTokenNow(curUserName, curUserID, true, false)
		// 删除一个评论
		delCommentResp := e.POST("/douyin/comment/action/").
			WithQuery("token", curUsrToken).WithQuery("video_id", videoId).WithQuery("action_type", 2).WithQuery("comment_id", cmm.CommentID).
			WithFormField("token", curUsrToken).WithFormField("video_id", videoId).WithFormField("action_type", 2).WithFormField("comment_id", cmm.CommentID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		delCommentResp.Value("status_code").Number().IsEqual(0)
		// 查看剩余评论数量
		commentListFavResp := e.GET("/douyin/comment/list/").
			WithQuery("token", loginedToken).WithQuery("video_id", videoId).
			WithFormField("token", loginedToken).WithFormField("video_id", videoId).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		commentListFavResp.Value("status_code").Number().IsEqual(0)
		if GEN_COMMENT_COUNT-1-idx > 0 {
			commentListFavResp.Value("comment_list").Array().Length().IsEqual(GEN_COMMENT_COUNT - 1 - idx)
		} else {
			commentListFavResp.NotContainsKey("comment_list")
		}
	}

}

func BenchmarkComment(b *testing.B) {
	e := newBenchExpect(b)
	rand.Seed(time.Now().UnixNano())
	videoIdx := rand.Intn(len(videos))
	videoId := videos[videoIdx].ID
	var commentsAdded = make([]commenterTestStruct, 0, GEN_COMMENT_COUNT)
	for i := 0; i < GEN_COMMENT_COUNT; i++ {
		userIdx := rand.Intn(len(users))
		userID := users[userIdx].ID
		userName := users[userIdx].Name
		flwCNt, flwrCnt := users[userIdx].FollowCount, users[userIdx].FollowerCount
		// token := getTestUserToken(userID, true, false)
		token := getTestUserTokenNow(userName, userID, true, false)
		// 添加一个评论的请求
		addCommentResp := e.POST("/douyin/comment/action/").
			WithQuery("token", token).WithQuery("video_id", videoId).WithQuery("action_type", COMMENT_ACTION_ADD).WithQuery("comment_text", "测试评论").
			WithFormField("token", token).WithFormField("video_id", videoId).WithFormField("action_type", COMMENT_ACTION_ADD).WithFormField("comment_text", "测试评论").
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		addCommentResp.Value("status_code").Number().IsEqual(0)
		addCommentResp.Value("status_msg").String().IsEqual("success")
		addCommentResp.Value("comment").Object().Value("id").Number().Gt(0)
		commentID := uint(addCommentResp.Value("comment").Object().Value("id").Number().Raw())
		addCommentResp.Value("comment").Object().Value("user").Object().Value("id").Number().IsEqual(userID)
		if flwCNt > 0 {
			addCommentResp.Value("comment").Object().Value("user").Object().Value("follow_count").Number().IsEqual(flwCNt)
		}
		if flwrCnt > 0 {
			addCommentResp.Value("comment").Object().Value("user").Object().Value("follower_count").Number().IsEqual(flwrCnt)
		}
		addCommentResp.Value("comment").Object().Value("user").Object().Value("name").String().IsEqual(userName)
		commentsAdded = append(commentsAdded, commenterTestStruct{commentID, userID, userName, int(flwCNt), int(flwrCnt)})
		addCommentResp.Value("comment").Object().Value("content").String().IsEqual("测试评论")
		addCommentResp.Value("comment").Object().Value("create_date").String().IsEqual(time.Now().Format("01-02"))
	}

	// token := getTestUserToken(0, false, false)
	token := getTestUserTokenNow("", 0, false, false)
	// 查询视频评论未登录
	commentListResp := e.GET("/douyin/comment/list/").
		WithQuery("token", token).WithQuery("video_id", videoId).
		WithFormField("token", token).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListResp.Value("status_code").Number().IsEqual(0)
	commentListResp.Value("comment_list").Array().Length().IsEqual(GEN_COMMENT_COUNT)
	containTestComment := false
	for idx, element := range commentListResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		cmm := commentsAdded[GEN_COMMENT_COUNT-idx-1]
		comment.Value("id").Number().IsEqual(cmm.CommentID)
		comment.Value("user").Object().Value("id").Number().IsEqual(cmm.UserID)
		comment.Value("user").Object().Value("name").String().IsEqual(cmm.UserName)
		if cmm.UserFollowCount > 0 {
			comment.Value("user").Object().Value("follow_count").Number().IsEqual(cmm.UserFollowCount)
		} else {
			comment.Value("user").Object().NotContainsKey("follow_count")
		}
		if cmm.UserFollowerCount > 0 {
			comment.Value("user").Object().Value("follower_count").Number().IsEqual(cmm.UserFollowerCount)
		} else {
			comment.Value("user").Object().NotContainsKey("follower_count")
		}
		comment.Value("content").String().NotEmpty().IsEqual("测试评论")
		comment.Value("create_date").String().NotEmpty().IsEqual(time.Now().Format("01-02"))
		containTestComment = true
	}

	assert.True(b, containTestComment, "Can't find test comment in list")

	// 查询视频评论已登录
	loginedUserIdx := rand.Intn(len(users))
	loginedUserID := users[loginedUserIdx].ID
	loginedUserName := users[loginedUserIdx].Name
	// loginedToken := getTestUserToken(loginedUserID, true, false)
	loginedToken := getTestUserTokenNow(loginedUserName, loginedUserID, true, false)
	commentListFavResp := e.GET("/douyin/comment/list/").
		WithQuery("token", loginedToken).WithQuery("video_id", videoId).
		WithFormField("token", loginedToken).WithFormField("video_id", videoId).
		Expect().
		Status(http.StatusOK).
		JSON().Object()
	commentListFavResp.Value("status_code").Number().IsEqual(0)
	commentListFavResp.Value("comment_list").Array().Length().IsEqual(GEN_COMMENT_COUNT)
	containTestComment = false
	for idx, element := range commentListFavResp.Value("comment_list").Array().Iter() {
		comment := element.Object()
		cmm := commentsAdded[GEN_COMMENT_COUNT-idx-1]
		comment.Value("id").Number().IsEqual(cmm.CommentID)
		comment.Value("user").Object().Value("id").Number().IsEqual(cmm.UserID)
		comment.Value("user").Object().Value("name").String().IsEqual(cmm.UserName)
		if cmm.UserFollowCount > 0 {
			comment.Value("user").Object().Value("follow_count").Number().IsEqual(cmm.UserFollowCount)
		}
		if cmm.UserFollowerCount > 0 {
			comment.Value("user").Object().Value("follower_count").Number().IsEqual(cmm.UserFollowerCount)
		}
		key := [2]uint{loginedUserID, cmm.UserID}
		if _, has := following[key]; has {
			comment.Value("user").Object().Value("is_follow").Boolean().True()
		} else {
			comment.Value("user").Object().NotContainsKey("is_follow")
		}
		comment.Value("content").String().NotEmpty().IsEqual("测试评论")
		comment.Value("create_date").String().NotEmpty().IsEqual(time.Now().Format("01-02"))
		containTestComment = true
	}
	assert.True(b, containTestComment, "Can't find test comment in list")

	// 测试删除评论
	for idx, cmm := range commentsAdded {
		curUserID := cmm.UserID
		curUserName := cmm.UserName
		// curUsrToken := getTestUserToken(curUserID, true, false)
		curUsrToken := getTestUserTokenNow(curUserName, curUserID, true, false)
		// 删除一个评论
		delCommentResp := e.POST("/douyin/comment/action/").
			WithQuery("token", curUsrToken).WithQuery("video_id", videoId).WithQuery("action_type", 2).WithQuery("comment_id", cmm.CommentID).
			WithFormField("token", curUsrToken).WithFormField("video_id", videoId).WithFormField("action_type", 2).WithFormField("comment_id", cmm.CommentID).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		delCommentResp.Value("status_code").Number().IsEqual(0)
		// 查看剩余评论数量
		commentListFavResp := e.GET("/douyin/comment/list/").
			WithQuery("token", loginedToken).WithQuery("video_id", videoId).
			WithFormField("token", loginedToken).WithFormField("video_id", videoId).
			Expect().
			Status(http.StatusOK).
			JSON().Object()
		commentListFavResp.Value("status_code").Number().IsEqual(0)
		if GEN_COMMENT_COUNT-1-idx > 0 {
			commentListFavResp.Value("comment_list").Array().Length().IsEqual(GEN_COMMENT_COUNT - 1 - idx)
		} else {
			commentListFavResp.NotContainsKey("comment_list")
		}
	}
}
