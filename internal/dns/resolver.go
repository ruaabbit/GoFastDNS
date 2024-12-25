package dns

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func NewResolver(address string) (Resolver, error) {
	if strings.HasPrefix(address, "[/") {
		parts := strings.SplitN(address, "/]", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid domain-specific format")
		}
		address = parts[1]
	}

	// 解析协议
	if strings.Contains(address, "://") {
		u, err := url.Parse(address)
		if err != nil {
			return nil, err
		}

		switch u.Scheme {
		case "udp":
			return &UDPResolver{server: u.Host}, nil
		case "tcp":
			return &TCPResolver{server: u.Host}, nil
		case "tls":
			return &TLSResolver{server: u.Host}, nil
		default:
			return nil, fmt.Errorf("unsupported protocol: %s", u.Scheme)
		}
	}

	// 默认 UDP
	return &UDPResolver{server: address}, nil
}

func (r *UDPResolver) Resolve(domain string, timeout time.Duration) DNSResult {
	c := dns.Client{
		Timeout: timeout,
	}

	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	_, duration, err := c.Exchange(&m, net.JoinHostPort(r.server, "53"))

	return DNSResult{
		Server:          r.server,
		Domain:          domain,
		Protocol:        ProtocolUDP,
		ResponseTime:    duration,
		ResolutionError: err,
	}
}

func (r *TCPResolver) Resolve(domain string, timeout time.Duration) DNSResult {
	c := dns.Client{
		Net:     "tcp",
		Timeout: timeout,
	}

	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	_, duration, err := c.Exchange(&m, net.JoinHostPort(r.server, "53"))

	return DNSResult{
		Server:          r.server,
		Domain:          domain,
		Protocol:        ProtocolTCP,
		ResponseTime:    duration,
		ResolutionError: err,
	}
}

func (r *TLSResolver) Resolve(domain string, timeout time.Duration) DNSResult {
	c := dns.Client{
		Net:     "tcp-tls",
		Timeout: timeout,
	}

	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)

	_, duration, err := c.Exchange(&m, net.JoinHostPort(r.server, "853"))

	return DNSResult{
		Server:          r.server,
		Domain:          domain,
		Protocol:        ProtocolTLS,
		ResponseTime:    duration,
		ResolutionError: err,
	}
}
