package models

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/fyved24/douyin/configs"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	RedisDB   *redis.Client
	RedisLock *redislock.Client
)

var Ctx = context.Background()

func InitRedis() {

	redisConfig := configs.Settings.RedisConfigs
	RedisDB = redis.NewClient(&redis.Options{
		// redis服务器地址，ip:port格式，比如：192.168.1.100:6379
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB, // use default DB
	})
	// 全局分布式锁
	RedisLock = redislock.New(RedisDB)
	_, err := RedisDB.Ping(Ctx).Result()

	if err != nil {
		panic("failed to connect redis")
	}

}
func InitAllDB() {
	InitDB()
	InitRedis()
}
func InitDB() {
	var err error
	dsn := "douyin:douyinxiangmu@tcp(101.43.131.38:3306)/douyin?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB.Set("gorm:table_options", "ENGINE=InnoDB")
	if err != nil {
		panic("failed to connect database")
	}
	err = DB.AutoMigrate(&User{}, &Video{}, &Comment{}, &Comment{}, &Follower{}, &Following{}, &Favorite{}, &Message{})
	if err != nil {
		panic("failed to migrate database")
	}
}
