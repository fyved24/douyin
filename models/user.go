package models

type User struct {
	Model
	Name           string    `json:"name"`
	Password       string    `json:"-"`
	FollowCount    int64     `json:"follow_count"`
	FollowerCount  int64     `json:"follower_count"`
	TotalFavorited int64     `json:"total_favorited"`
	FavoriteCount  int64     `json:"favorite_count"`
	Videos         []Video   `gorm:"foreignKey:AuthorID" json:"-"`
	Comments       []Comment `json:"-"`
}
