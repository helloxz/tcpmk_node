package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// 存储所有CIDR段
var cidrList []*net.IPNet

var (
	ipNets IPNetSet  // 全局变量，存储IPv6地址段
	one    sync.Once // 用于确保加载操作只执行一次
	loaded bool      // 标记是否已加载
)

// IPNetSet 用于存储IPv6地址段
type IPNetSet []*net.IPNet

// LoadIPNetsFromFile 从文件中加载IPv6地址段
func LoadIPNetsFromFile(filename string) (IPNetSet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ipNets IPNetSet
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		_, ipNet, err := net.ParseCIDR(line)
		if err != nil {
			return nil, err
		}
		ipNets = append(ipNets, ipNet)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ipNets, nil
}

// ContainsIP 判断给定的IPv6地址是否在IPNetSet中
func (ipNets IPNetSet) ContainsIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	for _, ipNet := range ipNets {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// IsIPInRanges 判断IP是否在预定义的IPv6地址段内
func IsCNIPv6(ipStr string) bool {
	// 确保只加载一次
	one.Do(func() {
		var err error
		ipNets, err = LoadIPNetsFromFile("data/ipdb/all_cn_ipv6.txt")
		if err != nil {
			fmt.Println("Error loading IP ranges:", err)
			return
		}
		loaded = true
	})

	if !loaded {
		return false
	}

	return ipNets.ContainsIP(ipStr)
}

// 加载CIDR段列表
func loadCNIPList() error {
	filePath := "data/ipdb/all_cn.txt"
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		_, cidr, err := net.ParseCIDR(line)
		if err != nil {
			fmt.Printf("无效的CIDR段: %s\n", line)
			continue
		}
		cidrList = append(cidrList, cidr)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// 判断IP是否在CIDR段内
func IsIPInCN(ip string) bool {
	// 确保CIDR段已加载
	if len(cidrList) == 0 {
		err := loadCNIPList()
		if err != nil {
			fmt.Println("加载CIDR段失败:", err)
			return false
		}
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		fmt.Println("无效的IP地址")
		return false
	}

	for _, cidr := range cidrList {
		if cidr.Contains(parsedIP) {
			// 打印IP段
			// fmt.Printf("IP: %s 在CIDR段: %s 内\n", ip, cidr.String())
			return true
		}
	}
	return false
}
