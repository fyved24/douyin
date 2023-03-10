package models

type User struct {
	Model
	Name            string `json:"name"`
	Password        string `json:"-"`
	FollowCount     uint   `json:"follow_count"`
	FollowerCount   uint   `json:"follower_count"`
	TotalFavorited  uint   `json:"total_favorited"`
	FavoriteCount   uint   `json:"favorite_count"`
	Avatar          string
	BackgroundImage string
	WorkCount       uint
	Signature       string
	Videos          []Video   `gorm:"foreignKey:AuthorID" json:"-"`
	Comments        []Comment `json:"-"`
}

// HasExistUserByUsername 根据用户名查找是否存在该用户
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

// AddUser 使用用户名、加密后的密码以及令牌新建一个用户
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

// SelectIDByUsernameAndPassword 如果能根据用户名和密码找到用户，返回用户ID；否则返回0表示找不到
func SelectIDByUsernameAndPassword(username string, password string) (bool, uint) {
	var user User
	DB.Where("name = ? AND password = ?", username, password).First(&user)
	if user.ID > 0 {
		return true, user.ID
	} else {
		return false, 0
	}
}

// SelectUsernameByID 根据用户ID查找用户名
func SelectUsernameByID(id uint) string {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.Name
}

// SelectFollowCountByID 根据用户ID查找用户关注数量
func SelectFollowCountByID(id uint) uint {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.FollowCount
}

// SelectFollowerCountByID 根据用户ID查找用户粉丝数量
func SelectFollowerCountByID(id uint) uint {
	var user User
	DB.Where("id = ?", id).First(&user)
	return user.FollowerCount
}

// SelectWorkCountByID 根据用户ID查找某个用户的视频数量
func SelectWorkCountByID(userID uint) uint {
	var user User
	DB.Model(&User{}).Where("id = ?", userID).First(&user)
	return user.WorkCount
}

// SelectFavoriteCountByID 根据用户ID查找某个用户点赞过的视频数量
func SelectFavoriteCountByID(id uint) uint {
	var user User
	DB.Model(&User{}).Where("id = ?", id).First(&user)
	return user.FavoriteCount
}

// SelectTotalFavoritedByID 根据用户ID查找某个用户被点赞的个数
func SelectTotalFavoritedByID(id uint) uint {
	var user User
	DB.Model(&User{}).Where("id = ?", id).First(&user)
	return user.TotalFavorited
}
