package api

import "github.com/gin-gonic/gin"

func Pong(c *gin.Context) {
	// 返回pong 文本
	c.String(200, "pong")
}
