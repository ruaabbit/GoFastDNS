package benchmark

import (
	"GoFastDNS/internal/dns"
	"sync"
	"time"
)

func RunBenchmark(servers []string, domains []string, attempts int, timeout time.Duration) []BenchmarkResult {
	var results []BenchmarkResult
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	for _, server := range servers {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			var total time.Duration
			domainResults := make([]DomainResult, 0, len(domains))
			successCount := 0
			totalQueries := len(domains)
			totalRetries := 0

			for _, domain := range domains {
				result := dns.ResolveDNS(s, domain, attempts, timeout)
				domainResult := DomainResult{
					Domain:       domain,
					ResponseTime: result.ResponseTime,
					Error:        result.ResolutionError,
					RetryCount:   result.RetryCount,
				}
				domainResults = append(domainResults, domainResult)
				totalRetries += result.RetryCount // 累加重试次数

				if result.ResolutionError == nil {
					total += result.ResponseTime
					successCount++
				}
			}

			// 计算成功率时将重试次数加入总查询数
			actualTotalQueries := totalQueries + totalRetries
			successRate := float64(successCount) / float64(actualTotalQueries)

			// 避免除零错误
			var avgResponseTime time.Duration
			if successCount > 0 {
				avgResponseTime = total / time.Duration(successCount)
			}

			mu.Lock()
			results = append(results, BenchmarkResult{
				Server:          s,
				AvgResponseTime: avgResponseTime,
				DomainResults:   domainResults,
				SuccessRate:     successRate,
				TotalRetries:    totalRetries,
			})
			mu.Unlock()
		}(server)
	}
	wg.Wait()
	return results
}
