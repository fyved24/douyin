package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddComment(t *testing.T) {
	InitDB()
	DB.Exec("DELETE FROM comments")
	DB.Exec("DELETE FROM videos")
	DB.Exec("DELETE FROM users")
	user := User{Name: "1号"}
	DB.Create(&user)
	insertVideos := []Video{{AuthorID: user.ID, PlayUrl: "1"}, {AuthorID: user.ID, PlayUrl: "2"}}
	DB.Create(&insertVideos)
	for _, video := range insertVideos {
		if cp, err := AddComment(video.ID, user.ID, "test test"); err != nil {
			t.Error(err)
		} else {
			t.Logf("inserted comment: %#v", *cp)
		}
	}
}

func TestDeleteComment(t *testing.T) {
	InitDB()
	DB.Exec("DELETE FROM comments")
	DB.Exec("DELETE FROM videos")
	DB.Exec("DELETE FROM users")
	user := User{Name: "1号"}
	DB.Create(&user)
	insertVideos := []Video{{AuthorID: user.ID, PlayUrl: "1"}, {AuthorID: user.ID, PlayUrl: "2"}}
	DB.Create(&insertVideos)
	for _, video := range insertVideos {
		if cp, err := AddComment(video.ID, user.ID, "test test"); err != nil {
			t.Error(err)
		} else {
			t.Logf("inserted comment: %#v", *cp)
			if err := DeleteComment(cp.ID); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestQueryCommentsByVideoID(t *testing.T) {
	InitDB()
	DB.Exec("DELETE FROM comments")
	DB.Exec("DELETE FROM videos")
	DB.Exec("DELETE FROM users")
	user := User{Name: "1号"}
	DB.Create(&user)
	insertVideos := []Video{{AuthorID: user.ID, PlayUrl: "1"}, {AuthorID: user.ID, PlayUrl: "2"}}
	DB.Create(&insertVideos)
	var cps = map[uint][]*Comment{}
	for _, video := range insertVideos {
		if cp, err := AddComment(video.ID, user.ID, "test test"); err != nil {
			t.Error(err)
		} else {
			cps[video.ID] = append(cps[video.ID], cp)
		}
	}

	for _, video := range insertVideos {
		if liteComments, err := QueryCommentsByVideoID(video.ID, -1, -1); err != nil {
			t.Error(err)
		} else {
			for idx, liteComment := range liteComments {
				assert.Equal(t, cps[video.ID][idx].ID, liteComment.ID)
				if err := DeleteComment(liteComment.ID); err != nil {
					t.Error()
				}
			}
			if liteComments, err := QueryCommentsByVideoID(video.ID, -1, -1); err != nil {
				t.Error(err)
			} else {
				assert.Len(t, liteComments, 0)
			}
		}
	}
}
