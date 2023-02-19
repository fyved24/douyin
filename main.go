package main

import (
	"fmt"
	"github.com/fyved24/douyin/configs"
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/router"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {

	configs.InitConfig()
	app := gin.Default()
	models.InitAllDB()
	router.InitRouter(app)
	pprof.Register(app)
	app.Run(":8080")
	app.Run(fmt.Sprintf(":%d", configs.Settings.Port))
}
