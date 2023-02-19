package test

import (
	"math/rand"
	"net/http"
	"testing"

	"time"

	jwtutils "github.com/fyved24/douyin/handlers/user/utils"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/services/comment"
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

func newBenchExpect(t *testing.B) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		Client:   http.DefaultClient,
		BaseURL:  serverAddr,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func getTestUserTokenBefore(userID uint, logined bool, expired bool) string {
	claims := comment.MySimpleUserClaims{
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(comment.MySecretKey)
	if err != nil {
		panic(err)
	}
	return ss
}

func getTestUserTokenNow(userName string, userID uint, logined bool, expired bool) string {
	return jwtutils.GetUserToken(userName, "", userID, logined)
}

// 本地生成一些伪造的jwt
func getTestUserToken(userID uint, logined bool, expired bool) string {
	return getTestUserTokenBefore(userID, logined, expired)
}

const BIG_UINT = 1000

// 给数据库生成一些用户和视频
func makeSomeUsersAndVideos() (users []models.User, videos []models.Video) {
	rand.Seed(time.Now().UnixNano())
	count := rand.Int31n(BIG_UINT) + 1
	users = make([]models.User, count)
	videos = make([]models.Video, count)
	// 生成一些用户
	for idx := range users {
		name, _ := uuid.GenerateUUID()
		// flwCnt, flwrCnt := rand.Int31(), rand.Int31()
		flwCnt, flwrCnt := 0, 0
		users[idx] = models.User{
			Name:          name,
			FollowCount:   uint(flwCnt),
			FollowerCount: uint(flwrCnt),
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

func makeSomeFollows(users []models.User) map[[2]uint]struct{} {
	rand.Seed(time.Now().UnixNano())
	following := make(map[[2]int]struct{})
	res := make(map[[2]uint]struct{})
	for i := 0; i < BIG_UINT*2; i++ {
		host, fl := rand.Intn(len(users)), rand.Intn(len(users))
		key := [2]int{host, fl}
		resKey := [2]uint{users[host].ID, users[fl].ID}
		if _, visited := following[key]; host == fl || visited {
			continue
		}
		following[key] = struct{}{}
		res[resKey] = struct{}{}
		models.DB.Create(&models.Following{HostID: users[host].ID, FollowID: users[fl].ID})
		models.DB.Create(&models.Follower{HostID: users[fl].ID, FollowerID: users[host].ID})
		users[host].FollowCount++
		users[fl].FollowerCount++
		models.DB.Model(&users[host]).Update("follow_count", users[host].FollowCount)
		models.DB.Model(&users[fl]).Update("follower_count", users[fl].FollowerCount)
	}
	return res
}
