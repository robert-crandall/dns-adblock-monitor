package handlers

import (
	"context"
)

func checkHost(host string) HostStatus {
	status := HostStatus{
		Host: host,
	}

	ips, err := resolver.LookupIPAddr(context.Background(), host)
	if err != nil {
		status.Error = err.Error()
		status.IsBlocked = true
		return status
	}

	status.IPv4 = make([]string, 0)
	status.IPv6 = make([]string, 0)
	status.UnblockedIPv4 = make([]string, 0)
	status.UnblockedIPv6 = make([]string, 0)

	// Separate IPv4 and IPv6 addresses
	for _, ip := range ips {
		if ip4 := ip.IP.To4(); ip4 != nil {
			ipStr := ip4.String()
			status.IPv4 = append(status.IPv4, ipStr)
			if !isIPv4Blocked(ipStr) {
				status.UnblockedIPv4 = append(status.UnblockedIPv4, ipStr)
			}
		} else {
			ipStr := ip.IP.String()
			status.IPv6 = append(status.IPv6, ipStr)
			if !isIPv6Blocked(ipStr) {
				status.UnblockedIPv6 = append(status.UnblockedIPv6, ipStr)
			}
		}
	}

	status.IsBlocked = len(status.UnblockedIPv4) == 0 && len(status.UnblockedIPv6) == 0
	return status
}
