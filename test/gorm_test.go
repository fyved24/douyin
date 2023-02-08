package test

import (
	"testing"

	"github.com/fyved24/douyin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func TestGorm(t *testing.T) {
	var err error
	dsn := "root:root123@tcp(localhost:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = DB.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{}, &models.Comment{}, &models.Follower{}, &models.Following{}, &models.Favorite{})
	if err != nil {
		panic("failed to migrate database")
	}
	u1 := models.User{Name: "admin", Password: "admin", FollowCount: 0, FollowerCount: 0, FavoriteCount: 0, TotalFavorited: 0}
	DB.Create(&u1)
}
