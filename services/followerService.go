package services

import (
	"errors"

	"github.com/fyved24/douyin/models"
	"gorm.io/gorm"
)

// 判断HostId是否有GuestId这个粉丝
func IsFollower(HostId uint, GuestId uint) bool {
	//1.数据模型准备
	var relationExist = &models.Follower{}
	//2.查询粉丝表中粉丝是否存在
	if err := models.DB.Model(&models.Follower{}).
		Where("host_id=? AND follower_id=?", HostId, GuestId).
		First(&relationExist).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		//粉丝不存在
		return false
	}
	//粉丝存在
	return true
}

// 增加HostId的粉丝数
func IncreaseFollowerCount(HostId uint) error {
	if err := models.DB.Model(&models.User{}).
		Where("id=?", HostId).
		Update("follower_count", gorm.Expr("follower_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 减少HostId的粉丝数
func DecreaseFollowerCount(HostId uint) error {
	if err := models.DB.Model(&models.User{}).
		Where("id=?", HostId).
		Update("follower_count", gorm.Expr("follower_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 新增粉丝
func CreateFollower(HostId uint, GuestId uint) error {

	//1.Follower数据准备
	newFollower := models.Follower{
		HostID:     HostId,
		FollowerID: GuestId,
	}

	//2.新建follower
	if err := models.DB.Model(&models.Follower{}).
		Create(&newFollower).Error; err != nil {
		return err
	}
	return nil
}

// 删除粉丝
func DeleteFollower(HostId uint, GuestId uint) error {
	//1.Follower数据准备
	newFollower := models.Follower{
		HostID:     HostId,
		FollowerID: GuestId,
	}

	//2.删除follower
	if err := models.DB.Model(&models.Follower{}).
		Where("host_id=? AND follower_id=?", HostId, GuestId).
		Delete(&newFollower).Error; err != nil {
		return err
	}

	return nil
}

// 获取粉丝表
func FollowerList(HostId uint) ([]models.User, error) {

	var userList []models.User

	// 1.获取粉丝id列表
	var followerIdList []uint
	if err := models.DB.Model(&models.Follower{}).
		Select("follower_id").
		Where("host_id = ?", HostId).
		Scan(&followerIdList).Error; err != nil {
		return userList, nil
	}

	// 2.根据粉丝id列表，在user表中查询
	if err := models.DB.Model(&models.User{}).
		Where("id IN ?", followerIdList).
		Scan(&userList).Error; err != nil {
		return userList, nil
	}

	return userList, nil
}
