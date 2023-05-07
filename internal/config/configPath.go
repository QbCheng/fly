package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	EnvConfigFilePath     = "G1_CONFIG_PATH"
	DefaultConfigFilePath = "config.yaml"
)

// FilePath 优先级: 命令行 > 环境变量 > 默认值
func FilePath() string {
	var config string
	flag.StringVar(&config, "c", "", "choose config file.")
	flag.Parse()
	if config != "" {
		// 命令行
		fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		return config
	}
	if configEnv := os.Getenv(EnvConfigFilePath); configEnv != "" {
		// 判断 EnvConfig 常量存储的环境变量是否为空
		fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", EnvConfigFilePath, config)
		return config
	}
	fmt.Printf("默认的配置文件路径, config的路径为%s\n", DefaultConfigFilePath)
	return DefaultConfigFilePath
}
