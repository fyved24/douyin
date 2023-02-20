package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	Model
	VideoID     uint
	UserID      uint
	Content     string
	PublishDate time.Time `gorm:"index"` // 为了按时间顺序排序加入索引
}

// 读出少量用户数据用的
type LiteUser struct {
	Name            string // 用户名称
	FollowCount     uint   // 关注总数
	FollowerCount   uint   // 粉丝总数
	Avatar          string // 用户头像
	BackgroundImage string // 用户个人页顶部大图
	FavoriteCount   uint   // 喜欢数
	Signature       string // 个人简介
	TotalFavorited  uint   // 获赞数量
	WorkCount       uint   // 作品数
}

// 将评论数据加上用户数据
type LiteComment struct {
	ID     uint
	UserID uint
	LiteUser
	Content     string
	PublishDate time.Time
}

// 返回相应视频ID的相应分页的所用评论行，以行的创建时间降序排列
// 感觉这种读取对事务性的要求不高，如果创建连接禁用默认事务可能会好点
// 现在还不知道到底是只查询评论,再去请求获得用户信息还是现在这样直接连接表得到大部分结果
func FindCommentsByVideoID(videoID uint, offset, limit int) (res []LiteComment, err error) {
	err = DB.Model(&Comment{}).
		Select("comments.id, comments.user_id, users.name, users.follow_count, users.follower_count, comments.content, comments.publish_date").
		Joins("left join users on comments.user_id = users.id").
		Where("comments.video_id = ?", videoID).
		Order("comments.publish_date DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

// 去掉表join获得user信息,改为单独获取comment信息之后再搜索用户信息
func FindCommentsByVideoIDWithoutUserInfo(videoID uint, offset, limit int) (res []LiteComment, err error) {
	err = DB.Model(&Comment{}).
		Select("comments.id, comments.user_id, comments.content, comments.publish_date").
		Where("comments.video_id = ?", videoID).
		Order("comments.publish_date DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

// 添加评论时应该已经有了发表评论的用户的信息了
func AddComment(videoID, userID uint, commentText string, publishDate time.Time) (*Comment, error) {
	comment := &Comment{VideoID: videoID, UserID: userID, Content: commentText, PublishDate: publishDate}
	if err := DB.Create(&comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// 删除评论
func DeleteComment(commentID, userID, videoID uint) (err error) {
	err = DB.Where("id = ? and user_id = ? and video_id = ?", commentID, userID, videoID).Delete(&Comment{}).Error
	return
}

// 获得用户的基本信息
func FindUserInfoByID(userID uint) (res *LiteUser, err error) {
	res = &LiteUser{}
	err = DB.Model(&User{}).Where("id = ?", userID).First(res).Error
	return
}

// 增加视频的评论计数
func IncreaseVideoCommentCount(videoID uint, adder int) (err error) {
	video := Video{}
	video.ID = videoID
	err = DB.Model(&video).Update("comment_count", gorm.Expr("comment_count + ?", adder)).Error
	return
}

// 给出userID用户所关注的所有用户的ID
func FindFollowedUsersByUserID(userID uint) (res []uint, err error) {
	err = DB.Model(&Following{}).Select("follow_id").Where("host_id = ?", userID).Find(&res).Error
	return
}

// 读取视频评论数的同时检查视频是否存在
func FindVideoCommentCountByID(videoID uint) (res uint, err error) {
	err = DB.Model(&Video{}).Where("id = ?", videoID).Select("comment_count").Take(&res).Error
	return
}

type LiteUserWithID struct {
	ID uint
	LiteUser
}

func FindUsersInfoByIDs(userID []uint) ([]LiteUserWithID, error) {
	var res []LiteUserWithID
	err := DB.Model(&User{}). // Select("id, name, follow_count, follower_count").
					Where("id in (?)", userID).Scan(&res).Error
	return res, err
}
