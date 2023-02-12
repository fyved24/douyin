package models

import "gorm.io/gorm"

type Follower struct {
	gorm.Model
	HostID     uint
	FollowerID uint
}
