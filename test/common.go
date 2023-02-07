package test

import (
	"math/rand"
	"net/http"
	"testing"

	"time"

	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/services"
	"github.com/gavv/httpexpect/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-uuid"
)

var serverAddr = "http://localhost:8080"

func newExpect(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

// 本地生成一些伪造的jwt
func getTestUserToken(userID uint, logined bool, expired bool) string {
	claims := services.MySimpleUserClaims{
		UserID:  userID,
		Logined: logined,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}
	if expired {
		claims.ExpiresAt = claims.IssuedAt
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	ss, _ := token.SignedString(services.MySecretKey)
	return ss
}

// 给数据库生成一些用户和视频
func makeSomeUsersAndVideos() (users []models.User, videos []models.Video) {
	rand.Seed(time.Now().UnixNano())
	count := rand.Int31n(100000)
	users = make([]models.User, count)
	videos = make([]models.Video, count)
	// 生成一些用户
	for idx := range users {
		name, _ := uuid.GenerateUUID()
		flwCnt, flwrCnt := rand.Int31(), rand.Int31()
		users[idx] = models.User{
			Name:          name,
			FollowCount:   int64(flwCnt),
			FollowerCount: int64(flwrCnt),
		}
	}
	// 生成一些视频
	models.DB.Create(&users)
	for idx := range videos {
		vName, _ := uuid.GenerateUUID()
		userIdx := rand.Int31n(count)
		videos[idx] = models.Video{
			AuthorID: users[userIdx].ID,
			Author:   users[userIdx],
			Title:    vName,
		}
	}
	models.DB.Create(&videos)
	return
}
