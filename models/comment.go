package models

import (
	"gorm.io/gorm"
	"time"
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
	Name          string
	FollowCount   int64
	FollowerCount int64
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
func QueryCommentsByVideoID(videoID uint, offset, limit int) (res []LiteComment, err error) {
	err = DB.Model(&Comment{}).
		Select("comments.id, comments.user_id, users.name, users.follow_count, users.follower_count, comments.content, comments.publish_date").
		Joins("left join users on comments.user_id = users.id").
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
func QueryUserBasicInfo(userID uint) (res *LiteUser, err error) {
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
func QueryFollowedUsersByUserID(userID uint) (res []uint, err error) {
	err = DB.Model(&Following{}).Select("follow_id").Where("host_id = ?", userID).Find(&res).Error
	return
}
