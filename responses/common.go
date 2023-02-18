package responses

import "errors"

type CommonResponse struct {
	StatusCode int32  `json:"status_code"`          // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg,omitempty"` // 返回状态描述
}

// 从simple-demo中抄来的json格式不确定正确性
type Comment struct {
	ID         int64  `json:"id,omitempty"`          // 评论id
	User       User   `json:"user"`                  // 评论用户信息
	Content    string `json:"content,omitempty"`     // 评论内容
	CreateDate string `json:"create_date,omitempty"` // 评论发布日期，格式 mm-dd
}

type User struct {
	ID            int64  `json:"id,omitempty"`             // 用户id
	Name          string `json:"name,omitempty"`           // 用户名称
	FollowCount   int64  `json:"follow_count,omitempty"`   // 关注总数
	FollowerCount int64  `json:"follower_count,omitempty"` // 粉丝总数
	IsFollow      bool   `json:"is_follow,omitempty"`      // true-已关注，false-未关注
	// 新添加的回复字段
	Avatar          string `json:"avatar,omitempty"`           // 用户头像
	BackgroundImage string `json:"background_image,omitempty"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count,omitempty"`   // 喜欢数
	Signature       string `json:"signature,omitempty"`        // 个人简介
	TotalFavorited  string `json:"total_favorited,omitempty"`  // 获赞数量
	WorkCount       int64  `json:"work_count,omitempty"`       // 作品数
}

var (
	ErrorUserNameNull    = errors.New("用户名为空")
	ErrorUserNameExtend  = errors.New("用户名长度不符合规范")
	ErrorPasswordNull    = errors.New("密码为空")
	ErrorPasswordLength  = errors.New("密码长度不符合规范")
	ErrorUserExit        = errors.New("用户已存在")
	ErrorFullPossibility = errors.New("用户不存在，账号或密码出错")
	ErrorNullPointer     = errors.New("空指针异常")
	ErrorPasswordFalse   = errors.New("密码错误")
	ErrorRelationExit    = errors.New("关注已存在")
	ErrorRelationNull    = errors.New("关注不存在")
)
