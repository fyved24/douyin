package models

import "gorm.io/gorm"

type Following struct {
	gorm.Model

	HostID   uint
	FollowID uint
}
