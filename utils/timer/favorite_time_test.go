package utils

import (
	"fmt"
	"github.com/fyved24/douyin/configs"
	"github.com/fyved24/douyin/models"
	"testing"
)

func TestUpdateVideoId(t *testing.T) {
	fmt.Println("test ")
	configs.InitConfig()
	models.InitAllDB()

	//var keys []string =[...]string{"video_count", "video_count"}
	keys := []string{"video_count38", "video_count41"}
	err := UpdateVideoCount(keys)
	fmt.Println(err)

}
