package models

import "gorm.io/gorm"

type Following struct {
	gorm.Model

	HostID   int64
	FollowID int64
}
