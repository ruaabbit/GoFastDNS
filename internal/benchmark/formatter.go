package benchmark

func FormatResults(results []BenchmarkResult) ([]CSVResult, map[string]*ServerStats) {
	domainMap := make(map[string]CSVResult)
	serverStats := make(map[string]*ServerStats)

	// 优化前面的 serverStats 初始化
	for _, result := range results {
		serverStats[result.Server] = &ServerStats{
			MinTime:     -1, // 将初始值设为-1表示未设置
			MaxTime:     0,
			TotalTime:   0,
			SuccessRate: 0,
		}
	}

	// 处理每个域名的结果
	for _, result := range results {
		stats := serverStats[result.Server]

		for _, dr := range result.DomainResults {
			if _, exists := domainMap[dr.Domain]; !exists {
				domainMap[dr.Domain] = CSVResult{
					Domain:      dr.Domain,
					ServerStats: make(map[string]DomainStats),
				}
			}

			responseTime := float64(dr.ResponseTime.Milliseconds())
			domainStat := DomainStats{
				ResponseTime: responseTime,
				Success:      dr.Error == nil,
				RetryCount:   dr.RetryCount,
			}
			if dr.Error != nil {
				domainStat.Error = dr.Error.Error()
			}

			// 更新服务器统计
			stats.TotalQueries++
			if dr.Error == nil && responseTime > 0 { // 添加 responseTime > 0 的判断
				stats.SuccessQueries++
				stats.TotalTime += responseTime

				// 首次成功请求时初始化 MinTime
				if stats.MinTime < 0 || responseTime < stats.MinTime {
					stats.MinTime = responseTime
				}
				if responseTime > stats.MaxTime {
					stats.MaxTime = responseTime
				}
			} else {
				stats.FailedQueries++
			}

			current := domainMap[dr.Domain]
			current.ServerStats[result.Server] = domainStat
			domainMap[dr.Domain] = current
		}

		// 计算服务器统计
		if stats.SuccessQueries > 0 {
			stats.AvgTime = stats.TotalTime / float64(stats.SuccessQueries)
			stats.SuccessRate = float64(stats.SuccessQueries) / float64(stats.TotalQueries) * 100
		}
	}

	// 转换为切片
	csvResults := make([]CSVResult, 0, len(domainMap))
	for _, result := range domainMap {
		csvResults = append(csvResults, result)
	}

	return csvResults, serverStats
}
