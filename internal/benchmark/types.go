package benchmark

import "time"

type DomainResult struct {
	Domain       string
	ResponseTime time.Duration
	Error        error
	RetryCount   int
}

type BenchmarkResult struct {
	Server          string
	AvgResponseTime time.Duration
	DomainResults   []DomainResult
	SuccessRate     float64
	TotalRetries    int
}
