package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断路由是否是/api/*开头的
		if strings.HasPrefix(c.Request.URL.Path, "/api/") || strings.HasPrefix(c.Request.URL.Path, "/page/") {
			c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token,X-Auth-Token")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, HEAD,OPTION")
			// c.Header("Access-Control-Allow-Origin", "*")
			if c.Writer.Header().Get("Access-Control-Allow-Origin") == "" {
				c.Header("Access-Control-Allow-Origin", "*")
			}
			c.Next() // 继续处理请求
		}
		c.Next() // 继续处理请求
	}
}

// 移除不要的header
func RemoveHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许所有的处理函数和其他中间件运行
		c.Next()

		// 在响应被发送前调整头部
		// 这里检查是否存在重复的 Access-Control-Allow-Origin 头，并保留一个
		origins := c.Writer.Header()["Access-Control-Allow-Origin"]
		if len(origins) > 1 {
			// 保留第一个发现的 Access-Control-Allow-Origin 头，移除其他的
			c.Header("Access-Control-Allow-Origin", origins[0])
		}
	}
}
