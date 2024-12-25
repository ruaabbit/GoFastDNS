package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DNSServers []string      `yaml:"dns_servers"`
	Domains    []string      `yaml:"domains"`
	Attempts   int           `yaml:"attempts"` // 最大重试次数
	Timeout    time.Duration `yaml:"timeout"`  // DNS 查询超时时间
}

func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("无法打开配置文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("无法关闭配置文件: %v", err)
		}
	}(file)

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}

	// 设置默认值
	if config.Attempts <= 0 {
		config.Attempts = 1
	}
	if config.Timeout <= 0 {
		config.Timeout = 2 * time.Second
	}

	return config
}
