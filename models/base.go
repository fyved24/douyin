package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id" redis:"id"`
	CreatedAt time.Time      `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" redis:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" redis:"deleted_at"`
}
