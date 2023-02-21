package comment

import (
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/dgraph-io/ristretto/z"
	jwtutils "github.com/fyved24/douyin/handlers/user/utils"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-uuid"
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
			switch k := key.(type) {
			case uint:
				return uint64(k), 0
			default:
				return z.KeyToHash(key)
			}

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
	cms, err := models.FindCommentsByVideoID(videoID, offset, limit)
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
		followedUsers, err := models.FindFollowedUsersByUserID(userID)
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

		// 新添加的用户信息内容
		res[idx].User.Avatar = cm.Avatar
		res[idx].User.BackgroundImage = cm.BackgroundImage
		res[idx].User.FavoriteCount = int64(cm.FavoriteCount)
		res[idx].User.Signature = cm.Signature
		res[idx].User.TotalFavorited = int64(cm.TotalFavorited)
		res[idx].User.WorkCount = int64(cm.WorkCount)

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

func getVideoCommentsWithoutUserInfo(videoID uint, limit, offset int) (res []responses.Comment, err error) {
	// 根据视频ID得到评论
	cms, err := models.FindCommentsByVideoIDWithoutUserInfo(videoID, offset, limit)
	if err != nil {
		logrus.Error(err)
		err = ErrCommentFetchFailed
		return
	}
	if len(cms) == 0 {
		return
	}
	res = make([]responses.Comment, len(cms))
	for idx, cm := range cms {
		res[idx].ID = int64(cm.ID)
		res[idx].User.ID = int64(cm.UserID)
		res[idx].Content = cm.Content
		// 根据评论创建时间生成评论创建日期字符串
		res[idx].CreateDate = cm.PublishDate.Format(CREATE_DATE_FMT)
	}
	return
}

// 获得视频的所有评论和评论用户的信息
// 如果浏览用户
func GetVideoComments(videoID uint, logined bool, userID uint) (res []responses.Comment, err error) {
	// res, err = getVideoComments(videoID, -1, -1, logined, userID)
	// res, err = getVideoCommentsWithCache(videoID, userID, logined)
	res, err = getVideoCommentsWithSeperateCache(videoID, userID, logined)
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
	res, err := models.FindUserInfoByID(userID)
	if err != nil {
		return nil, err
	}
	return res, err
}

func userBasicInfoWithCache(userID uint) (*models.LiteUser, error) {
	cacheInitOnce.Do(cacheInit)
	// 如果缓存中有用户信息直接取出使用
	localCacheLock.Lock()
	key := genCacheKey(USER_INFOS, userID)
	userInfoObj, _ := localCache.Get(key)
	userInfo, ok := userInfoObj.(models.LiteUser)
	if ok {
		defer localCacheLock.Unlock()
		cp := userInfo
		return &cp, nil
	}
	localCacheLock.Unlock()
	// 缓存未命中时从数据库读取用户信息
	res, err := userBasicInfo(userID)
	if err != nil {
		return nil, err
	}
	// 存入缓存
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	cp := *res
	localCache.Set(key, cp, 1)
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
	// usrInfo, err := userBasicInfo(userID)
	usrInfo, err := userBasicInfoWithCache(userID)
	if err != nil {
		logrus.Error(err)
		err = ErrUserFetchFailed
		return nil, err
	}
	// 返回信息
	var res = responses.Comment{
		ID:         int64(mr.ID),
		Content:    mr.Content,
		User:       responses.User{ID: int64(mr.UserID)},
		CreateDate: mr.PublishDate.Format(CREATE_DATE_FMT),
	}
	fillACommentUserInfo(&res, *usrInfo)
	return &res, nil
}

// 已登录用户在视频上发表评论
func AddVideoComment(videoID, userID uint, content string) (res *responses.Comment, err error) {
	// res, err = addVideoComment(videoID, userID, content)
	// return
	// return addVideoCommentWithCache(videoID, userID, content)
	return addVideoCommentWithSeperateCache(videoID, userID, content)
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
	// return deldeteCommentWithCache(commentID, userID, videoID)
	return deldeteCommentWithSeperateCache(commentID, userID, videoID)
}

var ErrFindVideoFailed = errors.New("find video failed")

func videoExist(videoID uint) (bool, error) {
	_, err := models.FindVideoCommentCountByID(videoID)
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

func videoExistWithCache(videoID uint) (bool, error) {
	// 可能的缓存初始化
	cacheInitOnce.Do(cacheInit)
	key := genCacheKey(VIDEO_COMMENTS, videoID)
	// 从本地缓存中读取出键对应的值
	valObj, hit := localCache.Get(key)
	if hit {
		// 如果存的确实是评论列表,认为视频是存在的
		// 否则认为是标识非法视频的标识
		switch val := valObj.(type) {
		case []responses.Comment:
			return true, nil
		case bool:
			return val, nil
		default:
			// 这里加入位置类型的报错可能比较好
			return false, nil
		}
	}
	// 本地缓存中没有视频的情况下,向数据库查找是否存在相应视频
	exist, err := videoExist(videoID)
	if err != nil {
		return false, err
	}
	// 加锁是为了防止复制产生破坏性效果
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	if !exist {
		// 如果数据库中不存在视频那么需要一个标识来表示其为不合法的视频ID
		localCache.SetWithTTL(key, false, 1, time.Second*time.Duration(rand.Intn(CACHE_TTL_RAND_SECONDS)+1))
		return false, nil
	} else {
		// 这种设置一般不会出现,因为对视频操作之前一般都会先访问该视频的评论列表
		localCache.Set(key, true, 1)
		return true, nil
	}
}

// 检查要操作或查询的视频是否存在
func VideoExist(videoID uint) (bool, error) {
	// return videoExist(videoID)
	return videoExistWithCache(videoID)
}

const (
	BASE_CACHE_TTL_MINUTES = 5
	CACHE_TTL_RAND_SECONDS = 60
)

// 带缓存的评论读取
func getVideoCommentsWithCache(videoID, userID uint, logined bool) (res []responses.Comment, err error) {
	//如果没有初始化过缓存本地缓存,初始化本地缓存
	cacheInitOnce.Do(cacheInit)
	// 本地缓存操作
	localCacheLock.Lock()
	// 查询缓存
	cachedObj, _ := localCache.Get(videoID)
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
	followedUsers, err := models.FindFollowedUsersByUserID(userID)
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
	cachedObj, _ := localCache.Get(videoID)
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
	cachedObj, _ := localCache.Get(videoID)
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

const (
	VIDEO_COMMENTS = "vc:"
	USER_INFOS     = "ui:"
	USER_FOLLOWS   = "uf:"
	VIDEO_AUTHOR   = "va:"
	COMMENT_LOCK   = "cl:"
)

// 用来生成本地缓存的key
func genCacheKey(prefix string, id uint) string {
	return prefix + strconv.FormatUint(uint64(id), 10)
}
func getUserFollowsWithCache(userID uint) (map[uint]struct{}, error) {
	// 接下来访问已登录用户的关注关系缓存
	var userFollowed, userFollowedCopy map[uint]struct{}
	key := genCacheKey(USER_FOLLOWS, userID)
	localCacheLock.Lock()
	userFollowedObj, _ := localCache.Get(key)
	userFollowed, ok := userFollowedObj.(map[uint]struct{})
	// 如果在本地缓存中查找到了用户的关注关系那复制一份出来
	if ok {
		userFollowedCopy = make(map[uint]struct{}, len(userFollowed))
		for k, _ := range userFollowed {
			userFollowedCopy[k] = struct{}{}
		}
	}
	localCacheLock.Unlock()
	// 如果本地缓存中没有找到该用户的关注关系,从数据库中请求到相应数据放到缓存中
	if !ok {
		followedUsers, err := models.FindFollowedUsersByUserID(userID)
		if err != nil {
			logrus.Error(err)
			err = ErrFollowingFetchFailed
			return nil, err
		}
		userFollowedCopy = make(map[uint]struct{}, len(followedUsers))
		userFollowed = make(map[uint]struct{}, len(followedUsers))
		for _, usr := range followedUsers {
			userFollowedCopy[usr] = struct{}{}
			userFollowed[usr] = struct{}{}
		}
		localCacheLock.Lock()
		localCache.Set(key, userFollowed, int64(len(userFollowed)))
		localCacheLock.Unlock()
	}
	return userFollowedCopy, nil
}
func getVideoCommentsRemoteAndFillCache(videoID uint) (res []responses.Comment, err error) {
	videoCacheKey := genCacheKey(VIDEO_COMMENTS, videoID)
	// 这里先不处理登录用户的关注关系
	// 因为我想缓存用户的关注关系
	// 尝试获取相应视频的评论缓存赋值锁
	lockID, locked, err := commentCacheInitLock(videoID)
	if err != nil {
		logrus.Error(err)
	}
	if !locked {
		// 如果已经被他人占用锁那么等待一段时间
		<-time.After(time.Millisecond * LOCK_POSSESS_DURATION_MILISEC)
		// 如果等待后缓存已经被设置,直接将缓存返回即可
		localCacheLock.Lock()
		defer localCacheLock.Unlock()
		cacheValObj, _ := localCache.Get(videoCacheKey)
		cacheVal, ok := cacheValObj.([]responses.Comment)
		if !ok {
			// 如果等待了也没有获得本地缓存中的数据
			// 认为可能有什么问题先返回这次请求
			return nil, ErrCommentFetchFailed
		}
		res = make([]responses.Comment, len(cacheVal))
		copy(res, cacheVal)
		return res, nil
	}
	// 如果是自己占有了本地缓存赋值锁,那么就需要在过程结束后释放锁
	defer func() {
		UnlockErr := commentCacheInitUnlock(videoID, lockID)
		if err != nil {
			logrus.Error(UnlockErr)
		}
	}()
	// 接下来是自己占用锁时的进行数据库访问和本地缓存赋值的过程
	res, err = getVideoCommentsWithoutUserInfo(videoID, -1, -1)
	if err != nil {
		return nil, err
	}

	localCacheLock.Lock()
	// 缓存评论区
	resCopy := make([]responses.Comment, len(res))
	copy(resCopy, res)
	localCache.Set(videoCacheKey, resCopy, int64(len(res)))
	localCacheLock.Unlock()
	return res, nil
}

func fillACommentUserInfo(needFill *responses.Comment, userInfo models.LiteUser) {
	needFill.User.FollowCount = int64(userInfo.FollowCount)
	needFill.User.FollowerCount = int64(userInfo.FollowerCount)
	needFill.User.Name = userInfo.Name
	needFill.User.Avatar = userInfo.Avatar
	needFill.User.BackgroundImage = userInfo.BackgroundImage
	needFill.User.Signature = userInfo.Signature
	needFill.User.FavoriteCount = int64(userInfo.FavoriteCount)
	needFill.User.TotalFavorited = int64(userInfo.TotalFavorited)
	needFill.User.WorkCount = int64(userInfo.WorkCount)
}

func fillCommentUsersInfoWithCache(needFill []responses.Comment) error {
	var needQuery []uint // 我们假定这种需要二次查找用户信息的情况不常见
	var needRefill []int
	localCacheLock.Lock()
	// 为每条评论填充评论用户的信息
	for idx := range needFill {
		userInfoObj, _ := localCache.Get(genCacheKey(USER_INFOS, uint(needFill[idx].User.ID)))
		userInfo, ok := userInfoObj.(models.LiteUser)
		if !ok {
			// 如果在本地缓存中没找到相应用户信息,记录要查找的用户信息
			needRefill = append(needRefill, idx)
			needQuery = append(needQuery, uint(needFill[idx].User.ID))
		} else {
			// 如果在存储中找到了对应用户的信息,直接赋值即可
			fillACommentUserInfo(&needFill[idx], userInfo)
		}
	}
	localCacheLock.Unlock()
	if needQuery != nil {
		res, err := models.FindUsersInfoByIDs(needQuery)
		if err != nil {
			logrus.Error(err)
			err = ErrUserFetchFailed
			return err
		}
		localCacheLock.Lock()
		defer localCacheLock.Unlock()
		for _, ele := range res {
			localCache.Set(genCacheKey(USER_INFOS, ele.ID), models.LiteUser{Name: ele.Name, FollowCount: ele.FollowCount, FollowerCount: ele.FollowerCount}, 1)
		}
		for _, idx := range needRefill {
			userInfoObj, _ := localCache.Get(genCacheKey(USER_INFOS, uint(needFill[idx].User.ID)))
			userInfo, ok := userInfoObj.(models.LiteUser)
			// 这些需要填写的用户信息是新写到缓存的不应该不存在
			if !ok {
				return ErrUserFetchFailed
			}
			fillACommentUserInfo(&needFill[idx], userInfo)
		}
	}
	return nil
}

func getVideoCommentsWithSeperateCache(videoID, userID uint, logined bool) (res []responses.Comment, err error) {
	//如果没有初始化过缓存本地缓存,初始化本地缓存
	cacheInitOnce.Do(cacheInit)
	// 本地缓存操作
	localCacheLock.Lock()
	// 查询缓存
	cachedObj, _ := localCache.Get(genCacheKey(VIDEO_COMMENTS, videoID))
	cachedComments, hit := cachedObj.([]responses.Comment)
	if !hit {
		// 如果缓存中没有，访问数据库得到
		localCacheLock.Unlock()
		// 从数据库读取评论并加载到本地缓存
		res, err = getVideoCommentsRemoteAndFillCache(videoID)
		if err != nil {
			return nil, err
		}
		// 填充所有评论的用户的信息
		if err = fillCommentUsersInfoWithCache(res); err != nil {
			return nil, err
		}
		// 这里也不返回结果了接下来要继续处理用户登录的情况
	} else {
		res = make([]responses.Comment, len(cachedComments))
		copy(res, cachedComments)
		localCacheLock.Unlock()
		// 填充所有评论的用户的信息
		if err = fillCommentUsersInfoWithCache(res); err != nil {
			return nil, err
		}
	}
	//  如果视频没有评论或是浏览者未登录,无需进一步修改
	if len(res) == 0 || !logined {
		return
	}
	// 从缓存中或从数据库中获得用户的关注关系
	userFollowedCopy, err := getUserFollowsWithCache(userID)
	// 将评论中用户的用户
	for idx := range res {
		// 如果用户登录了且发表评论的用户是浏览者关注的要标注
		_, following := userFollowedCopy[uint(res[idx].User.ID)]
		res[idx].User.IsFollow = following
	}
	return
}

func deldeteCommentWithSeperateCache(commentID, userID, videoID uint) error {
	// 先更新数据库再更新缓存
	err := deldeteComment(commentID, userID, videoID)
	if err != nil {
		return err
	}
	// 查找本地缓存
	cacheInitOnce.Do(cacheInit)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	// 与之前的带缓存删除之间只差一个查询key不同
	key := genCacheKey(VIDEO_COMMENTS, videoID)
	cachedObj, _ := localCache.Get(key)
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
	localCache.Set(key, cachedComments, int64(len(cachedComments)))
	return nil
}

func addVideoCommentWithSeperateCache(videoID, userID uint, content string) (*responses.Comment, error) {
	// 先更新数据库再更新缓存
	resp, err := addVideoComment(videoID, userID, content)
	if err != nil {
		return nil, err
	}
	// 操作本地缓存
	cacheInitOnce.Do(cacheInit)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	// 与之前的带缓存添加之间只差一个查询key不同
	key := genCacheKey(VIDEO_COMMENTS, videoID)
	cachedObj, _ := localCache.Get(key)
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
	localCache.Set(key, cachedComments, int64(len(cachedComments)))
	return resp, nil
}

type FollowActionEnm int

const (
	FOLLOW_ACTION_FOLLOW FollowActionEnm = 1 + iota
	FOLLOW_ACTION_UNFOLLOW
)

// 用户的关注状态变化时将本地缓存的关注状态更新
func ChangeFollowCacheStates(hostId, guestId uint, actionType FollowActionEnm) {
	cacheInitOnce.Do(cacheInit)
	// host (un)following guest
	// host is (not) a guest's follower
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	// 更新用户的关注列表
	hostFLKey := genCacheKey(USER_FOLLOWS, hostId)
	if flObj, ok := localCache.Get(hostFLKey); ok {
		if fl, ok := flObj.(map[uint]struct{}); ok {
			switch actionType {
			case FOLLOW_ACTION_FOLLOW:
				fl[guestId] = struct{}{}
			case FOLLOW_ACTION_UNFOLLOW:
				delete(fl, guestId)
			}
		}
	}
	// 更新用户的关注数
	hostUIKey := genCacheKey(USER_INFOS, hostId)
	if userInfoObj, ok := localCache.Get(hostUIKey); ok {
		if userInfo, ok := userInfoObj.(models.LiteUser); ok {
			switch actionType {
			case FOLLOW_ACTION_FOLLOW:
				userInfo.FollowCount++
			case FOLLOW_ACTION_UNFOLLOW:
				if userInfo.FollowCount > 0 {
					userInfo.FollowCount--
				}
			}
			// 缓存回存
			localCache.Set(hostUIKey, userInfo, 1)
		}
	}
	// 更新被关注用户的粉丝数
	guestUIKey := genCacheKey(USER_INFOS, guestId)
	if userInfoObj, ok := localCache.Get(guestUIKey); ok {
		if userInfo, ok := userInfoObj.(models.LiteUser); ok {
			switch actionType {
			case FOLLOW_ACTION_FOLLOW:
				userInfo.FollowerCount++
			case FOLLOW_ACTION_UNFOLLOW:
				if userInfo.FollowerCount > 0 {
					userInfo.FollowerCount--
				}
			}
			// 缓存回存
			localCache.Set(guestUIKey, userInfo, 1)
		}
	}
}

type FavoriteActionEnm int

const (
	FAVORITE_ACTION_FAVORITE FavoriteActionEnm = 1 + iota
	FAVORITE_ACTION_UNFAVORITE
)

func ChangeUserCacheFavoriteState(userID, videoID uint, actionType FavoriteActionEnm) {
	cacheInitOnce.Do(cacheInit)
	userKey := genCacheKey(USER_INFOS, userID)
	// 先给用户修改喜欢数量
	localCacheLock.Lock()
	userInfoObj, _ := localCache.Get(userKey)
	userInfo, ok := userInfoObj.(models.LiteUser)
	if ok {
		switch actionType {
		case FAVORITE_ACTION_FAVORITE:
			userInfo.FavoriteCount++
		case FAVORITE_ACTION_UNFAVORITE:
			userInfo.FavoriteCount--
		}
		localCache.Set(userKey, userInfo, 1)
	}
	localCacheLock.Unlock()
	// 再给视频作者更新被赞数
	// 但是首先要读取到视频的作者
	vaKey := genCacheKey(VIDEO_AUTHOR, videoID)
	localCacheLock.Lock()
	authorIDObj, _ := localCache.Get(vaKey)
	authorID, ok := authorIDObj.(uint)
	localCacheLock.Unlock()
	if !ok {
		var err error
		// 如果缓存中没有保存视频的作者是谁需要先从数据库中找到视频作者
		authorID, err = models.FindVideoAuthorByVideoID(videoID)
		if err != nil {
			// 这一步一般来说不应该出错,但是如果出错了可以直接返回,牺牲一致性保证服务
			logrus.Error(err)
			return
		}
		localCacheLock.Lock()
		defer localCacheLock.Unlock()
		// 缓存视频作者
		localCache.Set(vaKey, authorID, 1)
	}
	// 缓存中缓存了视频作者,那么直接更新作者的信息即可
	authorKey := genCacheKey(USER_INFOS, authorID)
	authorInfoObj, _ := localCache.Get(authorKey)
	authorInfo, inCache := authorInfoObj.(models.LiteUser)
	if inCache {
		switch actionType {
		case FAVORITE_ACTION_FAVORITE:
			authorInfo.TotalFavorited++
		case FAVORITE_ACTION_UNFAVORITE:
			authorInfo.TotalFavorited--
		}
		localCache.Set(authorKey, authorInfo, 1)
	}

}

func ChangeUserCacheWorkCount(userID uint) {
	cacheInitOnce.Do(cacheInit)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	key := genCacheKey(USER_INFOS, userID)
	userInfoObj, _ := localCache.Get(key)
	userInfo, ok := userInfoObj.(models.LiteUser)
	if ok {
		userInfo.WorkCount++
		localCache.Set(key, userInfo, 1)
	}
}

var ErrCommentCacheInitFailed = errors.New("comment cache init lock failed")

const (
	LOCK_POSSESS_DURATION_MILISEC = 50 // 最多占有锁50毫秒
)

// 只在评论列表的缓存初始化上加锁一个原因是,只有这个操作的数据量理论上是最大的
func commentCacheInitLock(videoID uint) (keyID string, lockPossessed bool, err error) {
	cacheInitOnce.Do(cacheInit)
	key := genCacheKey(COMMENT_LOCK, videoID)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	// 生成一个标识,来标识当前处理的内容
	keyID, err = uuid.GenerateUUID()
	if err != nil {
		logrus.Error(err)
		err = ErrCommentCacheInitFailed
		return "", false, err
	}
	// 如果已经有线程占有了锁,那么就不能再获取锁了
	_, ok := localCache.Get(key)
	if ok {
		return "", false, nil
	}
	// 否则自己申请占有锁并给出占有时间
	localCache.SetWithTTL(key, keyID, 1, time.Millisecond*LOCK_POSSESS_DURATION_MILISEC)
	return keyID, true, nil
}

func commentCacheInitUnlock(videoID uint, keyID string) error {
	cacheInitOnce.Do(cacheInit)
	key := genCacheKey(COMMENT_LOCK, videoID)
	localCacheLock.Lock()
	defer localCacheLock.Unlock()
	// 从缓存中读锁
	valObj, ok := localCache.Get(key)
	if !ok {
		// 可能锁已经超时了
		return nil
	}
	val, ok := valObj.(string)
	if !ok || val != keyID {
		// 可能锁被别人占有了
		return nil
	}
	// 如果是自己占有的锁,就要释放锁
	localCache.Del(key)
	return nil
}
