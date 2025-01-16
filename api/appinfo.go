package api

import "github.com/gin-gonic/gin"

// 声明程序版本
var Version = "1.1.0"

// 声明一个结构体
type appinfo struct {
	Version string `json:"version"`
}

func Appinfo(c *gin.Context) {
	var info appinfo
	info.Version = Version
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": info,
	})
}
