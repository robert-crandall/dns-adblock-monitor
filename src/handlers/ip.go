package handlers

import "net"

func addIPv4Block(ipv4 string) {
	// Try parsing as CIDR first
	_, network, err := net.ParseCIDR(ipv4)
	if err == nil && network != nil {
		config.BlockingIPv4 = append(config.BlockingIPv4, *network)
		return
	}

	// If not CIDR, treat as single IP
	ip := net.ParseIP(ipv4)
	if ip != nil {
		_, network, _ := net.ParseCIDR(ipv4 + "/32")
		if network != nil {
			config.BlockingIPv4 = append(config.BlockingIPv4, *network)
		}
	}
}

func addIPv6Block(ipv6 string) {
	// Try parsing as CIDR first
	_, network, err := net.ParseCIDR(ipv6)
	if err == nil && network != nil {
		config.BlockingIPv6 = append(config.BlockingIPv6, *network)
		return
	}

	// If not CIDR, treat as single IP
	ip := net.ParseIP(ipv6)
	if ip != nil {
		_, network, _ := net.ParseCIDR(ipv6 + "/128")
		if network != nil {
			config.BlockingIPv6 = append(config.BlockingIPv6, *network)
		}
	}
}

func isIPv4Blocked(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, network := range config.BlockingIPv4 {
		if network.Contains(parsedIP) {
			return true
		}
	}
	return false
}

func isIPv6Blocked(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Ensure we're using the full IPv6 representation
	parsedIP = parsedIP.To16()
	if parsedIP == nil {
		return false
	}

	for _, network := range config.BlockingIPv6 {
		if network.Contains(parsedIP) {
			return true
		}
	}
	return false
}
