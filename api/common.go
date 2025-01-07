package api

// 通用结构体文件
type IPInfo struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`  // 国家
	Region   string `json:"region"`   // 区域，一般代表省份
	City     string `json:"city"`     // 城市
	County   string `json:"county"`   // 区县
	ISP      string `json:"isp"`      // 云营商
	Area     string `json:"area"`     // 综合区域
	Zip      string `json:"zip"`      // 邮编
	ASN      string `json:"asn"`      // 自治域号
	Timezone string `json:"timezone"` // 时区
}
