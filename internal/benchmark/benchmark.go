package benchmark

import (
	"GoFastDNS/internal/dns"
	"GoFastDNS/internal/ping"
	"log"
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
				// DNS 解析
				result := dns.ResolveDNS(s, domain, attempts, timeout)

				// 执行 Ping 测试
				dnsPingResult := ping.PingDNSResult(result)

				domainResult := DomainResult{
					Domain:         domain,
					ResponseTime:   result.ResponseTime,
					Error:          result.ResolutionError,
					RetryCount:     result.RetryCount,
					DnsPingResults: dnsPingResult, // 添加 Ping 结果
				}

				domainResults = append(domainResults, domainResult)
				totalRetries += result.RetryCount

				if result.ResolutionError == nil {
					total += result.ResponseTime
					successCount++

					// 记录 Ping 结果
					if len(dnsPingResult.PingResults) > 0 {
						log.Printf("DNS服务器：%s，域名：%s，解析IP：%v，平均延迟：%v\n",
							s, domain, result.Answers, dnsPingResult.AvgRTT)
					}
				}
			}

			actualTotalQueries := totalQueries + totalRetries
			successRate := float64(successCount) / float64(actualTotalQueries)

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
			log.Printf("DNS服务器：%s，平均响应时间：%s，成功率：%.2f%%，总重试次数：%d\n",
				s, avgResponseTime, successRate*100, totalRetries)
		}(server)
	}
	wg.Wait()
	return results
}
