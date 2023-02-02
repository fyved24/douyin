package main

import (
	"github.com/fyved24/douyin/models"
	"github.com/fyved24/douyin/router"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	models.InitDB()
	router.InitRouter(app)
	app.Run(":8080")
}
