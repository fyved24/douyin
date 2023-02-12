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

// 根据用户名查找是否存在该用户
func HasExistUserByUsername(username string) bool {
	if len(username) == 0 {
		return false
	}
	var user User
	DB.Where("name = ?", username).First(&user)
	if user.ID > 0 {
		return true
	}
	return false
}

// 使用用户名、加密后的密码以及令牌新建一个用户
func AddUser(username string, password string, followCount uint, followerCount uint,
	totalFavorited uint, favoriteCount uint) uint {
	var user User
	user = User{
		Name:           username,
		Password:       password,
		FollowCount:    followCount,
		FollowerCount:  followerCount,
		TotalFavorited: totalFavorited,
		FavoriteCount:  favoriteCount,
	}
	DB.Create(&user)
	return user.ID
}

// 如果能根据用户名和密码找到用户，返回用户ID；否则返回0表示找不到
func SelectIDByUsernameAndPassword(username string, password string) (bool, uint) {
	var user User
	DB.Where("name = ? AND password = ?", username, password).First(&user)
	if user.ID > 0 {
		return true, user.ID
	} else {
		return false, 0
	}
}

// 查找是否拥有token为s的用户
func HasExistUserByToken(s string) bool {
	var user User
	DB.Where("token = ?", s).First(&user)
	return user.ID > 0
}

// 根据用户ID查找用户
func SelectUserByID(id uint) User {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user
}

// 根据用户ID查找用户名
func SelectUsernameByID(id uint) string {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.Name
}

func SelectFollowCountByID(id uint) uint {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.FollowCount
}

func SelectFollowerCountByID(id uint) uint {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.FollowerCount
}

func SelectTotalFavoritedByID(id uint) uint {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.TotalFavorited
}

func SelectFavoriteCountByID(id uint) uint {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.FavoriteCount
}
