package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string
	Password       string `json:"-"`
	FollowCount    int64
	FollowerCount  int64
	TotalFavorited int64
	FavoriteCount  int64
	Videos         []Video   `gorm:"foreignKey:AuthorID" json:"-"`
	Comments       []Comment `json:"-"`
}
