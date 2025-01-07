package api

import (
	"fmt"
	"net"
	"tcpmk_node/utils"

	"github.com/gin-gonic/gin"
)

func Traceroute(c *gin.Context) {
	// 设置适当的响应头以支持流式传输
	// 设置响应header头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 获取用户输入的主机名
	host := c.DefaultPostForm("host", "")
	// 判断主机名是否合法
	if !utils.IsValidIPOrDomain(host) {
		c.SSEvent("close", "Invalid host")
		c.Writer.Flush() // 添加这一行来刷新缓冲区
		return
	}
	// 解析为iPv4地址
	ip, _ := utils.ResolveIPv4(host)

	// 目标 IP 地址，可以替换为您要跟踪的地址
	// 目标 IP 地址，可以替换为您要跟踪的地址
	targetIP := net.ParseIP(ip)
	if targetIP == nil {
		c.SSEvent("close", "Invalid IP address")
		c.Writer.Flush() // 添加这一行来刷新缓冲区
		return
	}

	// 发起路由跟踪
	results, err := utils.TraceRoute(targetIP)
	if err != nil {
		c.SSEvent("close", "Failed to trace route")
		c.Writer.Flush() // 添加这一行来刷新缓冲区
		return
	}

	// 逐行读取并打印结果
	for result := range results {
		if result.ErrorMsg != "" {
			// fmt.Println("dsdsd")
			// fmt.Println(result.ErrorMsg)
			// 打印错误消息
			c.SSEvent("data", result.ErrorMsg)
			c.Writer.Flush() // 添加这一行来刷新缓冲区
			continue
		}
		if result.IP == nil {
			// fmt.Println("dsdsd")
			// fmt.Printf("%d\t*\n", result.Hop)
			continue
		}
		// 如果是内网IP则跳过
		if utils.IsPrivateIP(result.IP) {
			continue
		}
		// fmt.Printf("%d\t%s\t", result.Hop, result.IP)
		// 解析IP得到归属地
		address, _ := ParseIP2Location(result.IP.String())
		// 获取时间
		// 拼接为字符串
		var formattedString string
		for _, rtt := range result.RTTs {
			rttSeconds := rtt.Seconds() * 1000
			// 得到毫秒
			formattedString = fmt.Sprintf("%.0fms", rttSeconds)
		}
		// 拼接归属地
		output := result.IP.String() + "," + address.Country + " " + address.Region + " " + address.City + "," + formattedString
		// 可选：发送完成标志
		c.SSEvent("data", output)
		c.Writer.Flush() // 添加这一行来刷新缓冲区
	}

	// 可选：发送完成标志
	c.SSEvent("close", "Stream finished")
	c.Writer.Flush() // 添加这一行来刷新缓冲区
}
