package router

import (
	"tcpmk_node/api"
	"tcpmk_node/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Start() {
	//gin运行模式
	RunMode := viper.GetString("server.mode")
	// 设置运行模式
	gin.SetMode(RunMode)
	//运行gin
	r := gin.Default()

	//使用全局跨域中间件
	r.Use(middleware.CorsMiddleware())
	r.POST("/api/traceroute", middleware.Auth(), api.Traceroute)
	r.POST("/api/icmp/ping", middleware.Auth(), api.IcmpPing)

	// r.GET("/test", api.Test)

	//获取服务端配置
	port := ":" + viper.GetString("server.port")
	// 运行服务
	r.Run(port)

}
