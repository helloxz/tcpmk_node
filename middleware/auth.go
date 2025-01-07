package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// AuthMiddleware 是一个 Gin 中间件，用于验证 X-Auth-Token 是否为 admin
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 X-Auth-Token
		token := c.GetHeader("X-Auth-Token")
		if token == "" {
			// 如果 Header 中没有 X-Auth-Token，返回 401 未授权错误
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "X-Auth-Token header is required",
				"data": "",
			})
			c.Abort() // 终止后续处理
			return
		}

		// 获取配置中的token
		myToken := viper.GetString("app.token")

		// 检查 Token 是否为 admin
		if token != myToken {
			// 如果不是 admin，返回 403 禁止访问错误
			c.JSON(200, gin.H{
				"code": 403,
				"msg":  "Access denied. Only admin is allowed.",
				"data": "",
			})
			c.Abort() // 终止后续处理
			return
		}

		// 如果是 admin，继续后续处理
		c.Next()
	}
}
