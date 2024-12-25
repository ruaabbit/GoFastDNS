package main

import (
	"GoFastDNS/internal/benchmark"
	"GoFastDNS/internal/config"
	"fmt"
	"log"
)

func main() {
	// 初始化配置
	loadConfig := config.LoadConfig("config.yaml")

	// 运行基准测试
	results := benchmark.RunBenchmark(loadConfig.DNSServers, loadConfig.Domains, loadConfig.Attempts, loadConfig.Timeout)

	// 格式化结果
	formattedResults, serverStats := benchmark.FormatResults(results)

	// 保存结果到Excel文件
	filename := benchmark.SaveResultsToExcel(loadConfig.DNSServers, results, formattedResults, serverStats)
	if filename != "" {
		log.Printf("结果已保存到 %s\n", filename)
	}

	// 等待任意输入
	log.Print("按 任意 键退出...")
	fmt.Scanln()
}
