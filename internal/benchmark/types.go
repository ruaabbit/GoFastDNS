package benchmark

import (
	"GoFastDNS/internal/ping"
	"time"
)

type DomainResult struct {
	Domain       string
	ResponseTime time.Duration
	Error        error
	RetryCount   int
	DnsPingResults  ping.DNSPingResult // 新增字段
}

type BenchmarkResult struct {
	Server          string
	AvgResponseTime time.Duration
	DomainResults   []DomainResult
	SuccessRate     float64
	TotalRetries    int
}
