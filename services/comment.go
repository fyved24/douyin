package services

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	jwtutils "github.com/fyved24/douyin/handlers/user/utils"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 由于用户如果已经登录,给出的token中应该有服务端签发的用户ID,
// 而且系统没有设置,用户修改接口,因此认为服务端签发的用户ID默认合法
// 不进一步进行用户存在判断

// 评论创建日期的格式mm-dd
const CREATE_DATE_FMT = "01-02"

var ErrCommentFetchFailed = errors.New("get video's comments failed")
var ErrUserFetchFailed = errors.New("can't find user")
var ErrFollowingFetchFailed = errors.New("user subscribes find failed")

var localCacheLock sync.Mutex
var localCache *ristretto.Cache
var cacheInitOnce sync.Once

func cacheInit() {
	var err error
	localCache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		KeyToHash: func(key interface{}) (uint64, uint64) {
			k := key.(uint)
			return uint64(k), 0
		},
	})
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
}

// 没有引入cache获得视频的所有评论并且用连接表的方式获得用户的信息
func getVideoComments(videoID uint, limit, offset int, logined bool, userID uint) (res []responses.Comment, err error) {
	// 根据视频ID得到评论和评论发布者的一些基本用户信息
	cms, err := models.QueryCommentsByVideoID(videoID, offset, limit)
	if err != nil {
		logrus.Error(err)
		err = ErrCommentFetchFailed
		return
	}
	if len(cms) == 0 {
		return
	}
	var userFollowed = map[uint]struct{}{}
	if logined {
		// 如果浏览评论的是已登录的用户需要得到它关注的用户
		followedUsers, err := models.QueryFollowedUsersByUserID(userID)
		if err != nil {
			logrus.Error(err)
			err = ErrFollowingFetchFailed
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
	// res, err = getVideoComments(videoID, -1, -1, logined, userID)
	res, err = getVideoCommentsWithCache(videoID, userID, logined)
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

var ErrUserAuthFailed = errors.New("user authentication failed")

func currentBrowserLogined(tokenString *string) (logined bool, userID uint, err error) {
	if len(*tokenString) == 0 {
		return false, 0, nil
	}
	clm, err := jwtutils.ParseToken(*tokenString)
	if err != nil {
		logrus.Error(err)
		err = ErrUserAuthFailed
		return
	}
	logined = clm.IsLogin
	id, err := strconv.ParseUint(clm.UserID, 10, 64)
	if err != nil {
		logrus.Error(err)
		err = ErrUserAuthFailed
		return
	}
	userID = uint(id)
	return
}

// 用户鉴权测试
func BrowserLogined(tokenString *string) (logined bool, userID uint, err error) {
	// return internalTestBrowserLogined(tokenString)
	return currentBrowserLogined(tokenString)
}

// 查询评论用户的基本信息
func userBasicInfo(userID uint) (*models.LiteUser, error) {
	res, err := models.QueryUserBasicInfo(userID)
	if err != nil {
		return nil, err
	}
	return res, err
}

var ErrAddCommentFailed = errors.New("publish video comment failed")
var ErrModifyVideoStatsFailed = errors.New("change video comment count failed")

func addVideoComment(videoID, userID uint, content string) (*responses.Comment, error) {
	// 评论写数据库
	mr, err := models.AddComment(videoID, userID, content, time.Now())
	if err != nil {
		logrus.Error(err)
		err = ErrAddCommentFailed
		return nil, err
	}
	// 更新视频评论数应该不太需要原子性
	// TODO: 未来可能可以通过事务保证评论和评论数的原子性更新
	err = models.IncreaseVideoCommentCount(videoID, 1)
	if err != nil {
		logrus.Error(err)
		err = ErrModifyVideoStatsFailed
		return nil, err
	}
	// 发表评论用户的基本信息
	usrInfo, err := userBasicInfo(userID)
	if err != nil {
		logrus.Error(err)
		err = ErrUserFetchFailed
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
	// res, err = addVideoComment(videoID, userID, content)
	// return
	return addVideoCommentWithCache(videoID, userID, content)
}

var ErrDeleteCommentFailed = errors.New("delete video comment failed")

func deldeteComment(commentID, userID, videoID uint) error {
	// TODO: 原子化操作
	err := models.DeleteComment(commentID, userID, videoID)
	if err != nil {
		logrus.Error(err)
		err = ErrDeleteCommentFailed
		return err
	}
	err = models.IncreaseVideoCommentCount(videoID, -1)
	if err != nil {
		logrus.Error(err)
		err = ErrModifyVideoStatsFailed
		return err
	}
	return err
}

// 已登录用户删除自己发表的评论
func DeleteComment(commentID, userID, videoID uint) error {
	// return deldeteComment(commentID, userID, videoID)
	return deldeteCommentWithCache(commentID, userID, videoID)
}

var ErrFindVideoFailed = errors.New("find video failed")

func videoExist(videoID uint) (bool, error) {
	_, err := models.QueryVideoCommentCount(videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		logrus.Error(err)
		err = ErrFindVideoFailed
		return false, err
	}
	return true, nil
}

// 检查要操作或查询的视频是否存在
func VideoExist(videoID uint) (bool, error) {
	return videoExist(videoID)
}

const (
	INIT_COMMENT_BUFF_SIZE = 1000
	INIT_USER_BUFF_SIZE    = 1000
)

const BASE_CACHE_TTL_MINUTES = 5

// 带缓存的评论读取
func getVideoCommentsWithCache(videoID, userID uint, logined bool) (res []responses.Comment, err error) {
	//如果没有初始化过缓存本地缓存,初始化本地缓存
	cacheInitOnce.Do(cacheInit)
	// 本地缓存操作
	localCacheLock.Lock()
	// 查询缓存
	cachedObj, ok := localCache.Get(videoID)
	cachedComments, ok := cachedObj.([]responses.Comment)
	if !ok {
		// 如果缓存中没有，访问数据库得到
		localCacheLock.Unlock()
		res, err = getVideoComments(videoID, -1, -1, logined, userID)
		if err != nil {
			return
		}
		localCacheLock.Lock()
		localCache.SetWithTTL(videoID, res, int64(len(res)), time.Minute*BASE_CACHE_TTL_MINUTES)
		localCacheLock.Unlock()
		return
	}
	res = make([]responses.Comment, len(cachedComments))
	copy(res, cachedComments)
	localCacheLock.Unlock()
	//  如果视频没有评论或是浏览者未登录,无需进一步修改
	if len(res) == 0 || !logined {
		return
	}
	var userFollowed = map[uint]struct{}{}
	// 如果浏览评论的是已登录的用户需要得到它关注的用户
	followedUsers, err := models.QueryFollowedUsersByUserID(userID)
	if err != nil {
		logrus.Error(err)
		err = ErrFollowingFetchFailed
		return nil, err
	}
	for _, usr := range followedUsers {
		userFollowed[usr] = struct{}{}
	}
	for idx := range res {
		// 如果用户登录了且发表评论的用户是浏览者关注的要标注
		_, following := userFollowed[uint(res[idx].User.ID)]
		res[idx].User.IsFollow = following
	}
	return
}

func deldeteCommentWithCache(commentID, userID, videoID uint) error {
	// 先更新数据库再更新缓存
	err := deldeteComment(commentID, userID, videoID)
	if err != nil {
		return err
	}
	// 查找本地缓存
	cacheInitOnce.Do(cacheInit)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	cachedObj, ok := localCache.Get(videoID)
	cachedComments, ok := cachedObj.([]responses.Comment)
	// 没有缓存则不用更新缓存
	if !ok {
		return nil
	}
	// 如果找到缓存更新缓存,这部分代码整个本地缓存仍在被锁定,保证在并发中缓存更新内容不被破坏
	// 删掉缓存中的对应评论
	deleteFound := false
	for idx, ele := range cachedComments {
		if ele.ID == int64(commentID) {
			deleteFound = true
			for i := idx + 1; i < len(cachedComments); i++ {
				cachedComments[i-1] = cachedComments[i]
			}
			break
		}
	}
	if deleteFound {
		cachedComments = cachedComments[:len(cachedComments)-1]
	}
	// 将更新后的缓存存回去
	localCache.SetWithTTL(videoID, cachedComments, int64(len(cachedComments)), time.Minute*BASE_CACHE_TTL_MINUTES)
	return nil
}

func addVideoCommentWithCache(videoID, userID uint, content string) (*responses.Comment, error) {
	// 先更新数据库再更新缓存
	resp, err := addVideoComment(videoID, userID, content)
	if err != nil {
		return nil, err
	}
	// 操作本地缓存
	cacheInitOnce.Do(cacheInit)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	cachedObj, ok := localCache.Get(videoID)
	cachedComments, ok := cachedObj.([]responses.Comment)
	// 如果没有缓存不用操作
	if !ok {
		return resp, nil
	}
	// 如果找到缓存更新缓存,这部分代码整个本地缓存仍在被锁定,保证在并发中缓存更新内容不被破坏
	// 向缓存中插入新的评论
	cachedComments = append(cachedComments, responses.Comment{})
	for idx, ele := range cachedComments {
		if ele.CreateDate <= resp.CreateDate {
			for i := len(cachedComments) - 1; i > idx; i-- {
				cachedComments[i] = cachedComments[i-1]
			}
			cachedComments[idx] = *resp
			break
		}
	}
	// 将更新后的缓存存回去
	localCache.SetWithTTL(videoID, cachedComments, int64(len(cachedComments)), time.Minute*BASE_CACHE_TTL_MINUTES)
	return resp, nil
}

// 由于评论中带有用户的关注和粉丝数,因此每次关注关系更新直接清空缓存
func ClearLocalCache() {
	cacheInitOnce.Do(cacheInit)
	localCacheLock.Lock()
	localCache.Clear()
	localCacheLock.Unlock()
}
