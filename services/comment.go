package services

import (
	"errors"
	"time"

	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// 由于用户如果已经登录,给出的token中应该有服务端签发的用户ID,
// 而且系统没有设置,用户修改接口,因此认为服务端签发的用户ID默认合法
// 不进一步进行用户存在判断

// 评论创建日期的格式mm-dd
const CREATE_DATE_FMT = "01-02"

// 没有引入cache获得视频的所有评论并且用连接表的方式获得用户的信息
func getVideoComments(videoID uint, limit, offset int, logined bool, userID uint) (res []responses.Comment, err error) {
	// 根据视频ID得到评论和评论发布者的一些基本用户信息
	cms, err := models.QueryCommentsByVideoID(videoID, offset, limit)
	if err != nil || len(cms) == 0 {
		return
	}
	var userFollowed = map[uint]struct{}{}
	if logined {
		// 如果浏览评论的是已登录的用户需要得到它关注的用户
		followedUsers, err := models.QueryFollowedUsersByUserID(userID)
		if err != nil {
			return nil, err
		}
		for _, usr := range followedUsers {
			userFollowed[usr] = struct{}{}
		}
	}
	res = make([]responses.Comment, len(cms))
	for idx, cm := range cms {
		res[idx].ID = int64(cm.ID)
		res[idx].User.ID = int64(cm.UserID)
		res[idx].User.Name = cm.Name
		res[idx].User.FollowCount = int64(cm.FollowCount)
		res[idx].User.FollowerCount = int64(cm.FollowerCount)
		res[idx].Content = cm.Content
		// 根据评论创建时间生成评论创建日期字符串
		res[idx].CreateDate = cm.PublishDate.Format(CREATE_DATE_FMT)
		if logined {
			// 如果用户登录了且发表评论的用户是浏览者关注的要标注
			_, following := userFollowed[cm.UserID]
			res[idx].User.IsFollow = following
		}
	}
	return
}

// 获得视频的所有评论和评论用户的信息
// 如果浏览用户
func GetVideoComments(videoID uint, logined bool, userID uint) (res []responses.Comment, err error) {
	res, err = getVideoComments(videoID, -1, -1, logined, userID)
	return
}

var MySecretKey = []byte("test_jwt")

type MySimpleUserClaims struct {
	UserID  uint `json:"user_id"`
	Logined bool `json:"logined"`
	jwt.RegisteredClaims
}

func internalTestBrowserLogined(tokenString *string) (logined bool, userID uint, err error) {
	token, err := jwt.ParseWithClaims(*tokenString, &MySimpleUserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return MySecretKey, nil
	})
	if claims, ok := token.Claims.(*MySimpleUserClaims); ok && token.Valid {
		logined = claims.Logined
		userID = claims.UserID
	}
	return
}

// 用户鉴权测试
func BrowserLogined(tokenString *string) (logined bool, userID uint, err error) {
	return internalTestBrowserLogined(tokenString)
}

// 查询评论用户的基本信息
func userBasicInfo(userID uint) (*models.LiteUser, error) {
	res, err := models.QueryUserBasicInfo(userID)
	return res, err
}

func addVideoComment(videoID, userID uint, content string) (*responses.Comment, error) {
	// 评论写数据库
	mr, err := models.AddComment(videoID, userID, content, time.Now())
	if err != nil {
		return nil, err
	}
	// 更新视频评论数应该不太需要原子性
	// TODO: 未来可能可以通过事务保证评论和评论数的原子性更新
	err = models.IncreaseVideoCommentCount(videoID, 1)
	if err != nil {
		return nil, err
	}
	// 发表评论用户的基本信息
	usrInfo, err := userBasicInfo(userID)
	if err != nil {
		return nil, err
	}
	// 返回信息
	var res = responses.Comment{
		ID: int64(mr.ID),
		User: responses.User{
			ID:            int64(userID),
			Name:          usrInfo.Name,
			FollowCount:   int64(usrInfo.FollowCount),
			FollowerCount: int64(usrInfo.FollowerCount),
		},
		Content:    mr.Content,
		CreateDate: mr.PublishDate.Format(CREATE_DATE_FMT),
	}
	return &res, nil
}

// 已登录用户在视频上发表评论
func AddVideoComment(videoID, userID uint, content string) (res *responses.Comment, err error) {
	res, err = addVideoComment(videoID, userID, content)
	return
}

func deldeteComment(commentID, userID, videoID uint) error {
	// TODO: 原子化操作
	err := models.DeleteComment(commentID, userID, videoID)
	if err != nil {
		return err
	}
	err = models.IncreaseVideoCommentCount(videoID, -1)
	return err
}

// 已登录用户删除自己发表的评论
func DeleteComment(commentID, userID, videoID uint) error {
	return deldeteComment(commentID, userID, videoID)
}

func videoExist(videoID uint) (bool, error) {
	_, err := models.QueryVideoCommentCount(videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 检查要操作或查询的视频是否存在
func VideoExist(videoID uint) (bool, error) {
	return videoExist(videoID)
}
