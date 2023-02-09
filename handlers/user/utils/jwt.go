package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var mySigningKey []byte = []byte("askldfjaqwiopeklasdjqwerfasdfawerfsldfjkalsdfj") //密钥

type MyClaim struct { //jwt验证用
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   string `json:"user_id"`
	jwt.StandardClaims
}

// 签发一个token
func GetUserToken(username string, password string, userID uint) string {
	//创建一个JWT
	myclaim := MyClaim{
		Username: username,
		Password: password,
		UserID:   strconv.Itoa(int(userID)),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60,         //生效时间
			ExpiresAt: time.Now().Unix() + 2000*60*60, //失效时间，先设置为不过期
			Issuer:    "douyin",                       //签发者
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaim) //创建token（未加密），第一个参数是加密方法，第二个参数时自己定义的结构体
	s, err := jwtToken.SignedString(mySigningKey)                  //s就是已签发的token（已加密）
	if err != nil {
		fmt.Println(err)
	}
	return s
}

// 解析Token
func ParseToken(s string) (*MyClaim, error) {
	token, err := jwt.ParseWithClaims(s, &MyClaim{},
		func(t *jwt.Token) (interface{}, error) {
			return mySigningKey, nil
		})
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, errors.New("token is invalid")
	}

	if claims, ok := token.Claims.(*MyClaim); ok {
		return claims, nil
	}

	return nil, errors.New("token is invalid")
}
