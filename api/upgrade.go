package api

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// 升级接口
func Upgrade(c *gin.Context) {
	// 升级包地址
	url := "https://soft.xiaoz.org/UniBin/tcpmk_node/tcpmk_node.tar.gz"
	// 下载这个升级包到/opt/tcpmk_node/目录下
	// 下载这个升级包到/opt/tcpmk_node/目录下
	filePath := "/opt/tcpmk_node/tcpmk_node.tar.gz"
	// 创建目录（如果不存在）
	if err := os.MkdirAll("/opt/tcpmk_node", 0755); err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to create directory",
			"data": "",
		})
		return
	}
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to download file",
			"data": "",
		})
		return
	}
	defer resp.Body.Close()

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to create file",
			"data": "",
		})
		return
	}
	defer out.Close()

	// 将下载的内容写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Failed to write file",
			"data": "",
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "Upgrade successful",
		"data": "",
	})

	// 退出进程
	os.Exit(0)
}
