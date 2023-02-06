package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"` // 为了按时间顺序排序加入索引
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	VideoID   uint
	UserID    uint
	Content   string
}

type LiteUser struct {
	Name          string
	FollowCount   int64
	FollowerCount int64
}

type LiteComment struct {
	ID        uint
	CreatedAt time.Time
	UserID    uint
	LiteUser
	Content string
}

// 返回相应视频ID的相应分页的所用评论行，以行的创建时间降序排列
// 感觉这种读取对事务性的要求不高，如果创建连接禁用默认事务可能会好点
// 现在还不知道到底是只查询评论,再去请求获得用户信息还是现在这样直接连接表得到大部分结果
func QueryCommentsByVideoID(videoID uint, offset, limit int) (res []LiteComment, err error) {
	err = DB.Model(&Comment{}).
		Select("comments.id, comments.created_at, comments.user_id, users.name, users.follow_count, users.follower_count, comments.content").
		Joins("left join users on comments.user_id = users.id").
		Where("comments.video_id = ?", videoID).
		Order("comments.created_at DESC").Limit(limit).Offset(offset).Find(&res).Error
	return
}

// 添加评论时应该已经有了发表评论的用户的信息了
func AddComment(videoID, userID uint, commentText string) (*Comment, error) {
	comment := &Comment{VideoID: videoID, UserID: userID, Content: commentText}
	if err := DB.Create(&comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// 删除评论
func DeleteComment(commentID uint) (err error) {
	err = DB.Delete(&Comment{}, commentID).Error
	return
}
