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

// RelationAction 关注/取消关注操作
func RelationAction(c *gin.Context) {
	//1.取数据
	//1.1 从token中获取用户id
	token := c.Query("token")

	hostId := utils.GetUserIDFromToken(token)
	//1.2 获取待关注的用户id
	getToUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	guestId := uint(getToUserId)
	//1.3 获取操作类型（关注1，取消关注2）
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

	//1.获取用户id
	getUserId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	userId := uint(getUserId)

	//2.从数据库取用户的关注列表
	var err error
	var userList []models.User
	userList, err = services.FollowingList(userId)

	//3.构造返回数据
	var ReturnUserList = make([]ReturnUser, len(userList))
	for i, m := range userList {
		ReturnUserList[i].Id = m.ID
		ReturnUserList[i].Name = m.Name
		ReturnUserList[i].FollowCount = m.FollowCount
		ReturnUserList[i].FollowerCount = m.FollowerCount
		ReturnUserList[i].IsFollow = services.IsFollowing(userId, m.ID)
	}

	//4.响应返回
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

	//1.获取用户id
	getUserId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	userId := uint(getUserId)
	var err error
	//2.从数据库取粉丝列表
	var userList []models.User
	userList, err = services.FollowerList(userId)

	//3.构造返回数据
	var ReturnUserList = make([]ReturnUser, len(userList))
	for i, m := range userList {
		ReturnUserList[i].Id = m.ID
		ReturnUserList[i].Name = m.Name
		ReturnUserList[i].FollowCount = m.FollowCount
		ReturnUserList[i].FollowerCount = m.FollowerCount
		ReturnUserList[i].IsFollow = services.IsFollowing(userId, m.ID)
	}

	//4.响应返回
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

	//1.获取用户id
	getUserId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	userId := uint(getUserId)

	var err error
	//2.从数据库取关注、粉丝列表
	var followerList []models.User
	var followingList []models.User
	var friendList []models.User

	followerList, err = services.FollowerList(userId)
	followingList, err = services.FollowingList(userId)

	// 基于关注、粉丝表构造好友表 互相关注即为好友
	var mapper = make(map[uint]uint, len(followerList))
	for _, follower := range followerList {
		mapper[uint(follower.ID)] = userId
	}
	for _, following := range followingList {
		_, ok := mapper[uint(following.ID)]
		if ok {
			friendList = append(friendList, following)
		}
	}

	//3.构造返回数据
	var ReturnUserList = make([]ReturnUser, len(friendList))
	for i, m := range friendList {
		ReturnUserList[i].Id = m.ID
		ReturnUserList[i].Name = m.Name
		ReturnUserList[i].FollowCount = m.FollowCount
		ReturnUserList[i].FollowerCount = m.FollowerCount
		// ReturnUserList[i].IsFollow = services.IsFollowing(userId, m.ID)
		ReturnUserList[i].IsFollow = true
	}

	//4.响应返回
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
