package ping

import (
	"time"

	"GoFastDNS/internal/dns"

	probing "github.com/prometheus-community/pro-bing"
)

const (
	defaultCount    = 3
	defaultInterval = time.Millisecond * 100
	defaultTimeout  = time.Second * 2
)

func PingIP(ip string) PingResult {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return PingResult{
			IP:    ip,
			Error: err,
		}
	}
	pinger.SetPrivileged(true)
	pinger.Count = defaultCount
	pinger.Interval = defaultInterval
	pinger.Timeout = defaultTimeout

	err = pinger.Run()
	if err != nil {
		return PingResult{
			IP:    ip,
			Error: err,
		}
	}

	stats := pinger.Statistics()
	return PingResult{
		IP:          ip,
		RTT:         stats.AvgRtt,
		PacketLoss:  stats.PacketLoss,
		PacketsSent: stats.PacketsSent,
	}
}

func PingDNSResult(result dns.DNSResult) DNSPingResult {
	if result.ResolutionError != nil {
		return DNSPingResult{
			Domain:    result.Domain,
			DNSServer: result.Server,
			Error:     result.ResolutionError,
		}
	}

	pingResults := make([]PingResult, 0, len(result.Answers))
	var totalRTT time.Duration
	successfulPings := 0

	for _, ip := range result.Answers {
		pingResult := PingIP(ip)
		pingResults = append(pingResults, pingResult)

		// 只计算成功的 ping 结果
		if pingResult.Error == nil {
			totalRTT += pingResult.RTT
			successfulPings++
		}
	}

	// 计算平均 RTT
	var avgRTT time.Duration
	if successfulPings > 0 {
		avgRTT = totalRTT / time.Duration(successfulPings)
	}

	return DNSPingResult{
		Domain:      result.Domain,
		DNSServer:   result.Server,
		PingResults: pingResults,
		AvgRTT:      avgRTT,
	}
}
