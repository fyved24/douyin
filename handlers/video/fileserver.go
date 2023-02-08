package video

import "github.com/gin-gonic/gin"

func FileServer(c *gin.Context) {
	filename := c.Param("filename")
	c.File("local_storage/" + filename)
}
