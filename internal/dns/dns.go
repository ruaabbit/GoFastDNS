package dns

import (
	"time"
)

func ResolveDNS(address, domain string, attempts int, timeout time.Duration) DNSResult {
	resolver, err := NewResolver(address)
	if err != nil {
		return DNSResult{
			Server:          address,
			Domain:          domain,
			ResolutionError: err,
		}
	}

	var result DNSResult
	for i := 0; i <= attempts; i++ {
		result = resolver.Resolve(domain, timeout)
		result.RetryCount = i
		if result.ResolutionError == nil {
			break
		}
	}

	return result
}
