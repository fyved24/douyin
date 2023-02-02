package models

import "gorm.io/gorm"

type Follower struct {
	gorm.Model
	HostID     int64
	FollowerID int64
}
