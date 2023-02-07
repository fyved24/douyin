package models_test

import (
	"testing"
	"time"

	"github.com/fyved24/douyin/models"
	"github.com/stretchr/testify/assert"
)

var exampleUserA = models.User{
	Name:        "exampleUserA",
	FollowCount: 1,
}

// A follows B
var exampleUserB = models.User{
	Name:          "exampleUserB",
	FollowerCount: 2,
}

// C follows B
var exampleUserC = models.User{
	Name:        "exampleUserC",
	FollowCount: 1,
}

var exampleFollowingRelCB, exampleFollowingRelAB models.Following
var exampleFollowerRelCB, exampleFollowerRelAB models.Follower

var exampleVideo = models.Video{
	Title: "exampleVideo",
}

var exampleCommentA = models.Comment{
	Content: "comment by exampleUserA.",
}

var exampleCommentB = models.Comment{
	Content: "comment by exampleUserB.",
}

var exampleCommentC = models.Comment{
	Content: "comment by exampleUserC.",
}

func init() {
	models.InitDB()
	models.DB.Create(&exampleUserA)
	models.DB.Create(&exampleUserB)
	models.DB.Create(&exampleUserC)
	// A follows B
	exampleFollowerRelAB = models.Follower{HostID: int64(exampleUserB.ID), FollowerID: int64(exampleUserA.ID)}
	exampleFollowingRelAB = models.Following{HostID: int64(exampleUserA.ID), FollowID: int64(exampleUserB.ID)}
	// C follows B
	exampleFollowerRelCB = models.Follower{HostID: int64(exampleUserB.ID), FollowerID: int64(exampleUserC.ID)}
	exampleFollowingRelCB = models.Following{HostID: int64(exampleUserC.ID), FollowID: int64(exampleUserB.ID)}
	models.DB.Create(&exampleFollowerRelAB)
	models.DB.Create(&exampleFollowerRelCB)
	models.DB.Create(&exampleFollowingRelAB)
	models.DB.Create(&exampleFollowingRelCB)
	exampleVideo.Author = exampleUserB
	exampleVideo.AuthorID = exampleUserB.ID
	models.DB.Create(&exampleVideo)
}

func assertComment(t *testing.T, expect, actrual *models.Comment) {
	assert.Equal(t, expect.ID, actrual.ID)
	assert.Equal(t, expect.Content, actrual.Content)
	assert.Equal(t, expect.UserID, actrual.UserID)
	assert.Equal(t, expect.VideoID, actrual.VideoID)
}

func assertLiteComment(t *testing.T, expectC *models.Comment, expectU *models.User, actrual *models.LiteComment) {
	assert.Equal(t, expectU.ID, actrual.UserID)
	assert.Equal(t, expectU.Name, actrual.LiteUser.Name)
	assert.Equal(t, expectU.FollowCount, actrual.LiteUser.FollowCount)
	assert.Equal(t, expectU.FollowerCount, actrual.LiteUser.FollowerCount)
	assert.Equal(t, expectC.ID, actrual.ID)
	assert.Equal(t, expectC.Content, actrual.Content)
}

func TestComment(t *testing.T) {

	comments := []*models.Comment{&exampleCommentA, &exampleCommentB, &exampleCommentC}
	users := []*models.User{&exampleUserA, &exampleUserB, &exampleUserC}
	for idx, comment := range comments {
		// Test add comments
		res, err := models.AddComment(exampleVideo.ID, users[idx].ID, comment.Content, time.Now())
		if err != nil {
			t.Error(err)
		}
		*comment = *res
		tmpComment := models.Comment{}
		err = models.DB.Find(&tmpComment, comment.ID).Error
		if err != nil {
			t.Error(err)
		}
		assertComment(t, comment, &tmpComment)

		// Test increase video comment count
		err = models.IncreaseVideoCommentCount(comment.VideoID, 1)
		if err != nil {
			t.Error(err)
		}
		err = models.DB.Find(&exampleVideo, exampleVideo.ID).Error
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, int64(1+idx), exampleVideo.CommentCount)
	}

	// Test query comments
	qr, err := models.QueryCommentsByVideoID(exampleVideo.ID, -1, -1)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, qr, len(comments))
	for idx, lc := range qr {
		assertLiteComment(t, comments[len(comments)-idx-1], users[len(comments)-idx-1], &lc)
	}

	// Test delete comment
	for idx, comment := range comments {
		err := models.DeleteComment(comment.ID, users[idx].ID, comment.VideoID)
		if err != nil {
			t.Error(err)
		}
		tmpComment := models.Comment{}
		err = models.DB.Take(&tmpComment, comment.ID).Error
		assert.NotNil(t, err)
		// Test decrease video comment count
		err = models.IncreaseVideoCommentCount(comment.VideoID, -1)
		if err != nil {
			t.Error(err)
		}
		err = models.DB.Find(&exampleVideo, exampleVideo.ID).Error
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, int64(len(comments)-1-idx), exampleVideo.CommentCount)
	}

	// Test query comments
	qr, err = models.QueryCommentsByVideoID(exampleVideo.ID, -1, -1)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, qr, 0)
}

func assertUserBasicInfo(t *testing.T, expect *models.User, actrual *models.LiteUser) {
	assert.Equal(t, expect.Name, actrual.Name)
	assert.Equal(t, expect.FollowCount, actrual.FollowCount)
	assert.Equal(t, expect.FollowerCount, actrual.FollowerCount)
}

func TestUserBasicInfo(t *testing.T) {
	users := []models.User{exampleUserA, exampleUserB, exampleUserC}
	for _, user := range users {
		ub, err := models.QueryUserBasicInfo(user.ID)
		if err != nil {
			t.Error(err)
		}
		err = models.DB.Find(&user).Error
		if err != nil {
			t.Error(err)
		}
		assertUserBasicInfo(t, &user, ub)
	}
}

func TestFollow(t *testing.T) {
	// A follows B
	followed, err := models.QueryFollowedUsersByUserID(exampleUserA.ID)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, followed, 1)
	assert.Equal(t, followed[0], exampleUserB.ID)

	// C follows B
	followed, err = models.QueryFollowedUsersByUserID(exampleUserC.ID)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, followed, 1)
	assert.Equal(t, followed[0], exampleUserB.ID)

	followed, err = models.QueryFollowedUsersByUserID(exampleUserB.ID)
	if err != nil {
		t.Error(err)
	}
	assert.Len(t, followed, 0)
}
