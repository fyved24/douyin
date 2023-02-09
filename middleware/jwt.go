package middleware

import (
	"github.com/fyved24/douyin/handlers/user/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWT验证拦截，若token解析失败直接http返回，使用方法：路由组.Use()。
// 根据token获取userID和userName可参考user/user.go中Info函数里面的代码
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			token = c.PostForm("token")
		}

		//并未获取到token，拦截
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": 1,
				"status_msg":  "invalidToken",
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)

		//token解析失败
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": 1,
				"status_msg":  err.Error(),
			})
			c.Abort()
			return
		}

		userID := claims.UserID
		username := claims.Username
		c.Set("user_id", userID)    // 保存userID到Context的key中，可以通过Get()取
		c.Set("username", username) //保存username到Context的key中

		c.Next()
	}
}
