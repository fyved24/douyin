package services

import (
	"errors"

	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/responses"
	"gorm.io/gorm"
)

// 判断HostId是否关注GuestId
func IsFollowing(HostId uint, GuestId uint) bool {
	var relationExist = &models.Following{}
	//判断关注是否存在
	if err := models.DB.Model(&models.Following{}).
		Where("host_id=? AND follow_id=?", HostId, GuestId).
		First(&relationExist).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		//关注不存在
		return false
	}
	//关注存在
	return true
}

// 增加HostId的关注数
func IncreaseFollowCount(HostId uint) error {
	if err := models.DB.Model(&models.User{}).
		Where("id=?", HostId).
		Update("follow_count", gorm.Expr("follow_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 减少HostId的关注数
func DecreaseFollowCount(HostId uint) error {
	if err := models.DB.Model(&models.User{}).
		Where("id=?", HostId).
		Update("follow_count", gorm.Expr("follow_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// 新增关注
func CreateFollowing(HostId uint, GuestId uint) error {
	newFollowing := models.Following{
		HostID:   HostId,
		FollowID: GuestId,
	}

	if err := models.DB.Model(&models.Following{}).Create(&newFollowing).Error; err != nil {
		return err
	}
	return nil
}

// 删除关注
func DeleteFollowing(HostId uint, GuestId uint) error {
	deleteFollowing := models.Following{
		HostID:   HostId,
		FollowID: GuestId,
	}

	if err := models.DB.Model(&models.Following{}).Where("host_id=? AND follow_id=?", HostId, GuestId).Delete(&deleteFollowing).Error; err != nil {
		return err
	}

	return nil
}

// 实现关注操作 host关注guest
func FollowAction(HostId uint, GuestId uint, actionType uint) error {
	//关注操作
	if actionType == 1 {
		//判断关注是否存在
		if IsFollowing(HostId, GuestId) {
			//关注已存在
			return responses.ErrorRelationExit
		} else {
			//关注不存在,创建关注(启用事务Transaction)
			err1 := models.DB.Transaction(func(db *gorm.DB) error {
				// 新增关注 host关注guest
				err := CreateFollowing(HostId, GuestId)
				if err != nil {
					return err
				}
				// 新增粉丝 host是guest粉丝
				err = CreateFollower(GuestId, HostId)
				if err != nil {
					return err
				}
				//增加host的关注数
				err = IncreaseFollowCount(HostId)
				if err != nil {
					return err
				}
				//增加guest的粉丝数
				err = IncreaseFollowerCount(GuestId)
				if err != nil {
					return err
				}
				return nil
			})
			if err1 != nil {
				return err1
			}
		}
	}
	if actionType == 2 {
		//判断关注是否存在
		if IsFollowing(HostId, GuestId) {
			//关注存在,删除关注(启用事务Transaction)
			if err1 := models.DB.Transaction(func(db *gorm.DB) error {
				err := DeleteFollowing(HostId, GuestId)
				if err != nil {
					return err
				}
				err = DeleteFollower(GuestId, HostId)
				if err != nil {
					return err
				}
				//减少host的关注数
				err = DecreaseFollowCount(HostId)
				if err != nil {
					return err
				}
				//减少guest的粉丝数
				err = DecreaseFollowerCount(GuestId)
				if err != nil {
					return err
				}
				return nil
			}); err1 != nil {
				return err1
			}

		} else {
			//关注不存在
			return responses.ErrorRelationNull
		}
	}
	return nil
}

// FollowingList 获取关注表
func FollowingList(HostId uint) ([]models.User, error) {

	var userList []models.User

	// 1.获取关注id列表
	var followingIdList []uint
	if err := models.DB.Model(&models.Following{}).
		Select("follow_id").
		Where("host_id = ?", HostId).
		Scan(&followingIdList).Error; err != nil {
		return userList, nil
	}

	// 2.根据关注id列表，在user表中查询
	if err := models.DB.Model(&models.User{}).
		Where("id IN ?", followingIdList).
		Scan(&userList).Error; err != nil {
		return userList, nil
	}

	return userList, nil
}
