package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestQueryVideoListByLatestTime(t *testing.T) {
	InitDB()
	DB.Exec("DELETE FROM videos")
	DB.Exec("DELETE FROM users")
	user := User{Name: "1号"}
	DB.Create(&user)
	insertVideos := []Video{{AuthorID: user.ID, PlayUrl: "1"}, {AuthorID: user.ID, PlayUrl: "2"}}
	DB.Create(&insertVideos)
	findVideos, _ := QueryFeedVideoListByLatestTime(10, time.Now())
	for _, video := range *findVideos {
		b, _ := json.MarshalIndent(video, "", "  ")
		t.Log(string(b))
	}
	assert.Equal(t, len(*findVideos), 2)
	assert.Equal(t, user.ID, (*findVideos)[0].Author.ID)
}
