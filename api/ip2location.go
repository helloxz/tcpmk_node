package api

import (
	"errors"
	"net"

	"github.com/ip2location/ip2location-go/v9"
)

func ParseIP2Location(ip string) (IPInfo, error) {
	// IP数据库位置
	ipPath := "data/ipdb/"

	parsedIP := net.ParseIP(ip) // 解析 IP 字符串
	// 判断是否为 IPv4
	if parsedIP.To4() != nil {
		ipPath += "IP2LOCATION-LITE-DB3.BIN"
	} else {
		return IPInfo{}, errors.New("Unsupported IPv6 address!")
	}

	// 打开数据库
	db, err := ip2location.OpenDB(ipPath)
	if err != nil {
		return IPInfo{}, err
	}
	results, err := db.Get_all(ip)
	if err != nil {
		return IPInfo{}, err
	}
	// 关闭数据库
	defer db.Close()

	var ipInfo = IPInfo{
		IP:       ip,
		Country:  results.Country_long,
		City:     results.City,
		ISP:      "", // 免费数据库不支持results.Isp
		Region:   results.Region,
		County:   "", // 免费数据库不支持results.District
		Area:     "",
		Zip:      "",
		ASN:      "", // 免费数据库不支持results.Asn
		Timezone: "",
	}
	return ipInfo, nil
}
