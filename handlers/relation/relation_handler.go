package relation

import (
	"net/http"
	"strconv"

	"github.com/fyved24/douyin/handlers/user/utils"

	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"github.com/fyved24/douyin/services"
	"github.com/fyved24/douyin/services/comment"
	"github.com/gin-gonic/gin"
)

// ReturnUser 关注表与粉丝表共用的用户数据模型
type ReturnUser struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// UserListResponse 关注表，粉丝表，好友表公用结构体
type UserListResponse struct {
	responses.CommonResponse
	UserList []ReturnUser `json:"user_list"`
}

// // UserListResponse 粉丝表相应结构体
// type UserListResponse struct {
// 	responses.CommonResponse
// 	UserList []ReturnUser `json:"user_list"`
// }

// // UserListResponse 好友表相应结构体
// type FriendListResponse struct {
// 	responses.CommonResponse
// 	UserList []ReturnUser `json:"user_list"`
// }

// RelationAction 关注/取消关注操作
func RelationAction(c *gin.Context) {
	//1.取数据
	//1.1 从token中获取用户id
	token := c.Query("token")

	hostId := utils.GetUserIDFromToken(token)
	//1.2 获取待关注的用户id
	getToUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	guestId := uint(getToUserId)
	//1.3 获取关注操作（关注1，取消关注2）
	getActionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	actionType := uint(getActionType)

	//2.自己关注/取消关注自己不合法
	if hostId == guestId {
		c.JSON(http.StatusOK, responses.CommonResponse{
			StatusCode: 405,
			StatusMsg:  "无法关注自己",
		})
		c.Abort()
		return
	}
	var err error
	//3.service层进行关注/取消关注处理
	err = services.FollowAction(hostId, guestId, actionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.CommonResponse{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, responses.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "关注/取消关注成功！",
		})
		// 修改关注关系后更改本地缓存
		comment.ChangeFollowCacheStates(hostId, guestId, comment.FollowActionEnm(actionType))
	}
}

// FollowList 获取用户关注列表
func FollowList(c *gin.Context) {

	//1.数据预处理
	//1.1获取用户本人id
	token := c.Query("token")
	hostId := utils.GetUserIDFromToken(token)
	//1.2获取其他用户id
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestId := uint(getGuestId)

	//2.判断查询类型，从数据库取用户列表
	var err error
	var userList []models.User
	if guestId == 0 {
		//若其他用户id为0，代表查本人的关注表
		userList, err = services.FollowingList(hostId)
	} else {
		//若其他用户id不为0，代表查对方的关注表
		userList, err = services.FollowingList(guestId)
	}

	//构造返回的数据
	var ReturnUserList = make([]ReturnUser, len(userList))
	for i, m := range userList {
		ReturnUserList[i].Id = m.ID
		ReturnUserList[i].Name = m.Name
		ReturnUserList[i].FollowCount = m.FollowCount
		ReturnUserList[i].FollowerCount = m.FollowerCount
		ReturnUserList[i].IsFollow = services.IsFollowing(hostId, m.ID)
	}

	//3.响应返回
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			CommonResponse: responses.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			CommonResponse: responses.CommonResponse{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: ReturnUserList,
		})
	}
}

// FollowerList 获取用户粉丝列表
func FollowerList(c *gin.Context) {

	//1.数据预处理
	//1.1获取用户本人id
	token := c.Query("token")
	hostId := utils.GetUserIDFromToken(token)
	//1.2获取其他用户id
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestId := uint(getGuestId)
	var err error
	//2.判断查询类型
	var userList []models.User
	if guestId == 0 {
		//查本人的粉丝表
		userList, err = services.FollowerList(hostId)
	} else {
		//查对方的粉丝表
		userList, err = services.FollowerList(guestId)
	}

	//3.判断查询类型，从数据库取用户列表
	var ReturnUserList = make([]ReturnUser, len(userList))
	for i, m := range userList {
		ReturnUserList[i].Id = m.ID
		ReturnUserList[i].Name = m.Name
		ReturnUserList[i].FollowCount = m.FollowCount
		ReturnUserList[i].FollowerCount = m.FollowerCount
		ReturnUserList[i].IsFollow = services.IsFollowing(hostId, m.ID)
	}

	//3.处理
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			CommonResponse: responses.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			CommonResponse: responses.CommonResponse{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: ReturnUserList,
		})
	}
}

// FriendList 获取用户好友列表
func FriendList(c *gin.Context) {

	//1.数据预处理
	//1.1获取用户本人id
	token := c.Query("token")
	hostId := utils.GetUserIDFromToken(token)
	//1.2获取其他用户id
	getGuestId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	guestId := uint(getGuestId)
	var err error
	//2.判断查询类型
	var followerList []models.User
	var followingList []models.User
	var userList []models.User

	if guestId == 0 {
		//查本人的好友表
		followerList, err = services.FollowerList(hostId)
		followingList, err = services.FollowingList(hostId)
	} else {
		//查对方的好友表
		followerList, err = services.FollowerList(guestId)
		followingList, err = services.FollowingList(guestId)
	}
	var mapper = make(map[uint]uint, len(followerList))
	for _, follower := range followerList {
		mapper[uint(follower.ID)] = hostId
	}
	for _, following := range followingList {
		_, ok := mapper[uint(following.ID)]
		if ok {
			userList = append(userList, following)
		}
	}

	//3.判断查询类型，从数据库取用户列表
	var ReturnUserList = make([]ReturnUser, len(userList))
	for i, m := range userList {
		ReturnUserList[i].Id = m.ID
		ReturnUserList[i].Name = m.Name
		ReturnUserList[i].FollowCount = m.FollowCount
		ReturnUserList[i].FollowerCount = m.FollowerCount
		ReturnUserList[i].IsFollow = services.IsFollowing(hostId, m.ID)
	}

	//3.处理
	if err != nil {
		c.JSON(http.StatusBadRequest, UserListResponse{
			CommonResponse: responses.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "查找列表失败！",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, UserListResponse{
			CommonResponse: responses.CommonResponse{
				StatusCode: 0,
				StatusMsg:  "已找到列表！",
			},
			UserList: ReturnUserList,
		})
	}
}
