package utils

import (
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ResolveIPv4 解析域名的 IPv4 地址
func ResolveIPv4(domain string) (string, error) {
	// 解析域名的所有 IP 地址
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", fmt.Errorf("failed to resolve domain: %w", err)
	}

	// 查找第一个 IPv4 地址
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	// 如果没有找到 IPv4 地址
	return "", fmt.Errorf("no IPv4 address found for domain: %s", domain)
}

// 判断一个IP是否是内网IP
func IsPrivateIP(ip net.IP) bool {
	// 定义内网IP段
	privateIPBlocks := []*net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("169.254.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	// 遍历所有内网IP段，检查IP是否在其中
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

// isValidIPOrDomain 验证字符串是否是IP地址（IPv4或IPv6）或者域名
func IsValidIPOrDomain(input string) bool {
	// 检查是否是IP地址
	if net.ParseIP(input) != nil {
		return true
	}

	// 使用正则表达式检查是否是域名
	domainPattern := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
	re := regexp.MustCompile(domainPattern)
	if re.MatchString(input) {
		return true
	}

	// 如果既不是IP也不是域名，返回false
	return false
}

// 生成一个随机字符串
func RandStr(n int) string {
	// 使用本地随机数生成器，避免修改全局随机数生成器的状态
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 定义可用的字符集
	var bytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
	result := make([]byte, n)

	// 生成随机字符串
	for i := 0; i < n; i++ {
		result[i] = bytes[r.Intn(len(bytes))]
	}

	return string(result)
}

// 获取用户浏览器语言
func GetUserLang(c *gin.Context) string {
	acceptLanguage := c.Request.Header.Get("Accept-Language")
	// 将 Accept-Language 转换为小写，方便匹配
	acceptLanguageLower := strings.ToLower(acceptLanguage)

	// 判断是否包含 "zh"
	if strings.Contains(acceptLanguageLower, "zh") {
		return "zh"
	}
	// 默认返回 "en"
	return "en"
}
