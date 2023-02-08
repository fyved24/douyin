package models

type User struct {
	Model
	Name           string    `json:"name"`
	Password       string    `json:"-"`
	FollowCount    uint      `json:"follow_count"`
	FollowerCount  uint      `json:"follower_count"`
	TotalFavorited uint      `json:"total_favorited"`
	FavoriteCount  uint      `json:"favorite_count"`
	Videos         []Video   `gorm:"foreignKey:AuthorID" json:"-"`
	Comments       []Comment `json:"-"`
}
