package api

import (
	"strconv"
	"tcpmk_node/utils"
	"time"

	"github.com/gin-gonic/gin"
	probing "github.com/prometheus-community/pro-bing"
)

// 声明一个结构体，用来返回数据
type icmpinfo struct {
	PacketLoss int    `json:"packet_loss"`
	MaxRtt     int64  `json:"max_rtt"`
	MinRtt     int64  `json:"min_rtt"`
	AvgRtt     int64  `json:"avg_rtt"`
	PingCount  int    `json:"ping_count"`
	Host       string `json:"host"`
	IP         string `json:"ip"`
}

func IcmpPing(c *gin.Context) {
	// 获取用户输入的IP
	host := c.DefaultPostForm("host", "")
	// 获取ping的次数
	countStr := c.DefaultPostForm("count", "3")
	// 将字符串转换为整数
	count, err := strconv.Atoi(countStr)

	if err != nil {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid count",
			"data": "",
		})
		return
	}
	// 如果次数不在1-10之间，则返回错误
	if count < 1 || count > 10 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Count must be between 1 and 10",
			"data": "",
		})
		return
	}

	// 验证是否是合法的host
	if !utils.IsValidIPOrDomain(host) {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "Invalid host",
			"data": "",
		})
		return
	}

	// 将主机名解析为IPv4地址
	ip, err := utils.ResolveIPv4(host)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Domain name resolution failed",
			"data": "",
		})
		return
	}

	pinger, err := probing.NewPinger(ip)

	if err != nil {
		// fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Session creation failed!",
			"data": "",
		})
		return
	}
	// ping 3次
	pinger.Count = count
	// 这个是总的超时时间，单个包默认是1s，由于最大ping次数是10，所以这里设置为11s
	pinger.Timeout = 11 * time.Second
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		// fmt.Println(err)
		c.JSON(200, gin.H{
			"code": 500,
			"msg":  "Running failed!",
			"data": "",
		})
		return
	}

	// 获取丢包率
	stats := pinger.Statistics()
	var info icmpinfo
	//stats.PacketLoss 转为int，不需要保留小数点
	info.PacketLoss = int(stats.PacketLoss)

	// 参考这个获取时间：https://stackoverflow.com/questions/41503758/conversion-of-time-duration-type-microseconds-value-to-milliseconds
	info.MaxRtt = int64(stats.MaxRtt / time.Millisecond)
	info.MinRtt = int64(stats.MinRtt / time.Millisecond)
	info.AvgRtt = int64(stats.AvgRtt / time.Millisecond)
	info.PingCount = pinger.Count
	info.Host = host
	info.IP = ip

	// 获取丢包率
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "success",
		"data": info,
	})
}
