package dns

import (
	"time"
)

type Protocol string

const (
	ProtocolUDP Protocol = "UDP"
	ProtocolTCP Protocol = "TCP"
	ProtocolTLS Protocol = "TLS"
)

type DNSResult struct {
	Server          string
	Domain          string
	Protocol        Protocol
	ResponseTime    time.Duration
	ResolutionError error
	RetryCount      int
	Answers         []string    // 新增：保存DNS应答记录
}

type Resolver interface {
	Resolve(domain string, timeout time.Duration) DNSResult
}

type UDPResolver struct {
	server string
}

type TCPResolver struct {
	server string
}

type TLSResolver struct {
	server string
}
