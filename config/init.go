package config

import (
	"fmt"
	"os"
	"sync"
	"tcpmk_node/utils"

	"github.com/spf13/viper"
)

var (
	once sync.Once
)

func InitConfig() {
	once.Do(func() {
		//默认配置文件
		config_file := "data/config/config.toml"
		// 检查配置文件是否存在，不存在则拷贝一个
		if _, err := os.Stat(config_file); os.IsNotExist(err) {
			// fmt.Println("Config file not found, copying from config.example.toml")
			// 拷贝配置文件
			err := CopyFile("./default.toml", config_file)
			if err != nil {
				fmt.Println("Failed to copy config file:", err)
				os.Exit(1)
			}
		}

		viper.SetConfigFile(config_file) // 指定配置文件路径
		//指定ini类型的文件
		viper.SetConfigType("toml")
		err := viper.ReadInConfig() // 读取配置信息
		if err != nil {             // 读取配置信息失败
			// 写入日志
			fmt.Println("Failed to read config:", err)
			os.Exit(1)
		}
		// 读取一个配置
		token := viper.GetString("app.token")

		// 如果token为空，则随机设置一个
		if token == "" {
			// 生成一个随机数
			token = utils.RandStr(16)
			viper.Set("app.token", token)
			err := viper.WriteConfig()

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("dsds" + token)
		}
	})
}

// CopyFile 复制文件
func CopyFile(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	return
}
