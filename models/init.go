package models

import (
	"context"
	"log"

	"github.com/bsm/redislock"
	"github.com/fyved24/douyin/configs"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	RedisDB     *redis.Client
	RedisLock   *redislock.Client
	MinIOClient *minio.Client
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

func InitMinIO() {
	ctx := context.Background()
	minioConfig := configs.Settings.MinIOConfigs
	endpoint := minioConfig.Endpoint
	accessKeyID := minioConfig.AccessKeyID
	secretAccessKey := minioConfig.SecretAccessKey
	useSSL := minioConfig.UseSSL
	var err error
	// Initialize minio client object.
	MinIOClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	videoBucketName := "videos"
	imageBucketName := "images"
	location := "local"

	err = MinIOClient.MakeBucket(ctx, videoBucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinIOClient.BucketExists(ctx, videoBucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", videoBucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Video Bucket Successfully created %s\n", videoBucketName)
	}
	err = MinIOClient.MakeBucket(ctx, imageBucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := MinIOClient.BucketExists(ctx, imageBucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", imageBucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Image Bucket Successfully created %s\n", imageBucketName)
	}
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

func InitAllDB() {
	InitDB()
	InitRedis()
	InitMinIO()
}
