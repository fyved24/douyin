package user

import (
	"github.com/fyved24/douyin/handlers/user/utils"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	hasExistUser := models.HasExistUserByUsername(username) //查找有没有重复用户名

	//排除格式错误情况
	if hasExistUser {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "注册失败，该用户已存在",
			"user_id":     0,
			"token":       "NoToken",
		})
		return
	}
	if len(username) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "注册失败，用户名不能为空",
			"user_id":     0,
			"token":       "NoToken",
		})
		return
	}
	if len(username) > 32 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "注册失败，用户名最长为32个字符",
			"user_id":     0,
			"token":       "NoToken",
		})
		return
	}
	if len(password) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "注册失败，密码不能为空",
			"user_id":     0,
			"token":       "NoToken",
		})
		return
	}
	if len(password) > 32 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "注册失败，密码最长为32个字符",
			"user_id":     0,
			"token":       "NoToken",
		})
		return
	}

	password = utils.EncodePassword(password) //加密存储密码

	userID := models.AddUser(username, password, 0, 0,
		0, 0) //注册成功，更新用户信息
	token := utils.GetUserToken(username, password, userID, false) //获取JWT令牌
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "注册成功",
		"user_id":     userID,
		"token":       token,
	})
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	password = utils.EncodePassword(password) //加密存储密码
	hasExistUser, userID := models.SelectIDByUsernameAndPassword(username, password)
	if hasExistUser == false {
		//未找到该用户
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "用户名或密码错误",
			"user_id":     0,
			"token":       "NoToken",
		})
		return
	}
	//用户名和密码都正确，成功登录
	token := utils.GetUserToken(username, password, userID, true) //获取JWT令牌
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "登录成功",
		"user_id":     userID,
		"token":       token,
	})
}

func Info(c *gin.Context) {
	hostUserID := utils.StringToUint(c.Query("user_id"))  //空间页面host的userid
	hostUsername := models.SelectUsernameByID(hostUserID) //空间页面主人的username

	//根据token获取myUserId
	token := c.Query("token")
	myUserID := utils.GetUserIDFromToken(token)     //访问者的userID
	myUsername := utils.GetUsernameFromToken(token) //访问者的userName

	if myUsername == "" && hostUsername == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "无法得知访问者与被访问者",
		})
		return
	}
	if hostUsername == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "无法得知被访问者",
		})
		return
	}
	if myUsername == "" {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "无法得知访问者",
		})
		return
	}

	id := hostUserID                                                                                      //个人主页id
	name := hostUsername                                                                                  //个人主页username
	followCount := models.SelectFollowCountByID(id)                                                       //个人主页关注数量
	followerCount := models.SelectFollowerCountByID(id)                                                   //个人主页粉丝数量
	isFollow := services.IsFollower(hostUserID, myUserID)                                                 //是否关注该个人主页
	WorkCount := models.SelectWorkCountByID(id)                                                           //个人主页视频数量
	FavoriteCount := models.SelectFavoriteCountByID(id)                                                   //个人主页赞数量
	TotalFavorited := models.SelectTotalFavoritedByID(id)                                                 //个人主页获赞数量
	Signature := "这个用户很懒，没有留下任何信息……"                                                                      //个人主页简介
	Avatar := "https://adguycn990-typoraimage.oss-cn-hangzhou.aliyuncs.com/202211231910960.webp"          //个人主页头像
	BackgroundImage := "https://adguycn990-typoraimage.oss-cn-hangzhou.aliyuncs.com/202211231914310.webp" //个人主页背景图
	type AutoGenerated struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
		User       struct {
			ID              uint   `json:"id"`
			Name            string `json:"name"`
			FollowCount     uint   `json:"follow_count"`
			FollowerCount   uint   `json:"follower_count"`
			WorkCount       uint   `json:"work_count"`
			FavoriteCount   uint   `json:"favorite_count"`
			TotalFavorited  uint   `json:"total_favorited"`
			IsFollow        bool   `json:"is_follow"`
			Avatar          string `json:"avatar"`
			BackgroundImage string `json:"background_image"`
			Signature       string `json:"signature"`
		} `json:"user"`
	}

	resp := AutoGenerated{
		StatusCode: 0,
		StatusMsg:  "访问个人主页成功",
		User: struct {
			ID              uint   `json:"id"`
			Name            string `json:"name"`
			FollowCount     uint   `json:"follow_count"`
			FollowerCount   uint   `json:"follower_count"`
			WorkCount       uint   `json:"work_count"`
			FavoriteCount   uint   `json:"favorite_count"`
			TotalFavorited  uint   `json:"total_favorited"`
			IsFollow        bool   `json:"is_follow"`
			Avatar          string `json:"avatar"`
			BackgroundImage string `json:"background_image"`
			Signature       string `json:"signature"`
		}{
			ID:              id,
			Name:            name,
			FollowCount:     followCount,
			FollowerCount:   followerCount,
			WorkCount:       WorkCount,
			FavoriteCount:   FavoriteCount,
			TotalFavorited:  TotalFavorited,
			IsFollow:        isFollow,
			Avatar:          Avatar,
			BackgroundImage: BackgroundImage,
			Signature:       Signature,
		},
	}
	c.JSON(http.StatusOK, resp)
}
