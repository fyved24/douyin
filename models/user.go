package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name           string
	Password       string `json:"-"`
	FollowCount    uint
	FollowerCount  uint
	TotalFavorited uint
	FavoriteCount  uint
	Videos         []Video   `gorm:"foreignKey:AuthorID" json:"-"`
	Comments       []Comment `json:"-"`
}
