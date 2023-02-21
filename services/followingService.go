package services

import (
	"errors"

	"github.com/fyved24/douyin/models"
	"gorm.io/gorm"
)

// IsExistByID 根据用户ID判断用户是否存在
func IsExistByID(id uint) bool {
	var user = &models.User{}
	//判断关注是否存在
	if err := models.DB.Model(&models.User{}).
		Where("id=?", id).
		First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		//用户不存在
		return false
	}
	//用户存在
	return true
}

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
	// 过滤被关注用户不存在的情况
	if !IsExistByID(GuestId) {
		return errors.New("关注/取消关注对象不存在")
	}

	//关注操作
	if actionType == 1 {
		//判断关注是否存在
		if IsFollowing(HostId, GuestId) {
			//关注已存在
			return errors.New("关注已存在")
		} else {
			//关注不存在,创建关注(启用事务Transaction)
			err1 := models.DB.Transaction(func(tx *gorm.DB) error {
				// 新增关注 host关注guest
				if err := tx.Model(&models.Following{}).
					Create(&models.Following{
						HostID:   HostId,
						FollowID: GuestId,
					}).Error; err != nil {
					return err
				}

				// 新增粉丝 host是guest粉丝
				if err := tx.Model(&models.Follower{}).
					Create(&models.Follower{
						HostID:     GuestId,
						FollowerID: HostId,
					}).Error; err != nil {
					return err
				}

				//增加host的关注数
				if err := tx.Model(&models.User{}).
					Where("id=?", HostId).
					Update("follow_count", gorm.Expr("follow_count+?", 1)).Error; err != nil {
					print("增加host的关注数", err)
					return err
				}

				//增加guest的粉丝数
				if err := tx.Model(&models.User{}).
					Where("id=?", GuestId).
					Update("follower_count", gorm.Expr("follower_count+?", 1)).Error; err != nil {
					print("增加guest的粉丝数", err)
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
			if err1 := models.DB.Transaction(func(tx *gorm.DB) error {
				if err := tx.Model(&models.Following{}).
					Where("host_id=? AND follow_id=?", HostId, GuestId).Delete(&models.Following{
					HostID:   HostId,
					FollowID: GuestId,
				}).Error; err != nil {
					return err
				}

				if err := tx.Model(&models.Follower{}).
					Where("host_id=? AND follower_id=?", GuestId, HostId).Delete(&models.Follower{
					HostID:     GuestId,
					FollowerID: HostId,
				}).Error; err != nil {
					return err
				}

				//减少host的关注数
				if err := tx.Model(&models.User{}).
					Where("id=?", HostId).
					Update("follow_count", gorm.Expr("follow_count-?", 1)).Error; err != nil {
					return err
				}

				//减少guest的粉丝数
				if err := tx.Model(&models.User{}).
					Where("id=?", GuestId).
					Update("follower_count", gorm.Expr("follower_count-?", 1)).Error; err != nil {
					return err
				}
				return nil
			}); err1 != nil {
				return err1
			}

		} else {
			//关注不存在
			return errors.New("关注不存在")
		}
	}
	if (actionType != 1) && (actionType != 2) {
		return errors.New("actionType参数不合法")
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
