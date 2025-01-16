package api

import (
	"github.com/ipipdotnet/ipdb-go"
)

// 新版解析纯真IP数据库，使用的ipip数据库格式
func ParseQQwryNew(ip string) (IPInfo, error) {
	db, err := ipdb.NewCity("data/ipdb/qqwry.ipdb")
	if err != nil {
		return IPInfo{}, err
	}

	// fmt.Println(db.FindMap(ip, "CN")) // return map[string]string
	result, err := db.FindMap(ip, "CN")
	// 如果IP查找失败
	if err != nil {
		return IPInfo{}, err
	}

	var ipInfo = IPInfo{
		IP:      ip,
		Country: result["country_name"],
		Region:  result["region_name"],
		City:    result["city_name"],
		County:  "",
		ISP:     result["isp_domain"],
		Area:    "",
	}

	return ipInfo, nil
}
