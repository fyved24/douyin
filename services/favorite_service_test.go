package services

import (
	"testing"

	"github.com/fyved24/douyin/models"
	"github.com/stretchr/testify/assert"
)

func TestFavoriteAction1(t *testing.T) {
	models.InitDB()
	err := FavoriteAction(1, 1, 1)
	if err != nil {
		t.Log(err)
	}
}

func TestFavoriteAction2(t *testing.T) {
	models.InitDB()
	err := FavoriteAction(1, 2, 2)
	if err != nil {
		t.Log(err)
	}
}

func TestFindAllFavorite(t *testing.T) {
	models.InitDB()
	res, err := FindAllFavorite(1)

	t.Log(res)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(res), 1)

}
