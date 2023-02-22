package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/willf/bloom"
)

func BloomFilter() gin.HandlerFunc {
	filter := bloom.New(100000, 5) // 创建布隆过滤器
	return func(c *gin.Context) {
		c.Set("bloom", filter) // 存储布隆过滤器
		c.Next()
	}
}
