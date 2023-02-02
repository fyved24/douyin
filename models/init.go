package models

import "gorm.io/gorm"
import "gorm.io/driver/mysql"

var (
	DB *gorm.DB
)

func InitDB() {
	var err error
	dsn := "root:123456@tcp(127.0.0.1:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = DB.AutoMigrate(&User{}, &Video{}, &Comment{}, &Comment{}, &Follower{}, &Following{}, &Favorite{})
	if err != nil {
		panic("failed to migrate database")
	}
}
