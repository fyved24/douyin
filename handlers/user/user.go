package user

import (
	"fmt"
	"github.com/fyved24/douyin/handlers/user/utils"
	"github.com/fyved24/douyin/models"
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
	token := utils.GetUserToken(username, password, userID) //获取JWT令牌
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
	token := utils.GetUserToken(username, password, userID) //获取JWT令牌
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
	tokenStruct, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		})
		return
	}
	myUserID := tokenStruct.UserID     //访问者的userID
	myUsername := tokenStruct.Username //访问者的userName

	fmt.Println("host: " + hostUsername)
	fmt.Println("me: " + myUserID)
	fmt.Println("me: " + myUsername)

	//if username1 == "" && username2 == "" {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status_code": 1,
	//		"status_msg":  "无法根据ID找到用户名，也无法根据token找到用户名",
	//	})
	//	return
	//}
	//if username1 == "" {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status_code": 1,
	//		"status_msg":  "无法根据ID找到用户名",
	//	})
	//	return
	//}
	//if username2 == "" {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status_code": 1,
	//		"status_msg":  "无法根据token找到用户名",
	//	})
	//	return
	//}

}
