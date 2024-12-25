package benchmark

import "time"

// ServerStats 存储每个DNS服务器的统计信息
type ServerStats struct {
	TotalQueries   int
	SuccessQueries int
	FailedQueries  int
	TotalTime      float64
	AvgTime        float64
	SuccessRate    float64
	MinTime        float64
	MaxTime        float64
}

// CSVResult 包含域名结果和服务器统计信息
type CSVResult struct {
	Domain      string
	ServerStats map[string]DomainStats
}

type DomainStats struct {
	ResponseTime float64
	Success      bool
	RetryCount   int
	Error        string
}

// 修改 DomainResult 结构体
type DomainResult struct {
	Domain       string
	ResponseTime time.Duration
	Error        error
	RetryCount   int // 新增重试次数字段
}

// 修改 BenchmarkResult 结构体
type BenchmarkResult struct {
	Server          string
	AvgResponseTime time.Duration
	DomainResults   []DomainResult
	SuccessRate     float64
	TotalRetries    int // 新增总重试次数字段
}
