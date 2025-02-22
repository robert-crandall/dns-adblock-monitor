package handlers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
)

var (
	config Config
	mu     sync.RWMutex
)

// Initialize sets up the handlers package with configuration
func Initialize(hosts []string, blockingIPv4 []string, blockingIPv6 []string, dnsResolver string) {
	mu.Lock()
	defer mu.Unlock()

	config.Hosts = hosts

	// Parse IPv4 CIDR blocks
	config.BlockingIPv4 = make([]net.IPNet, 0)
	for _, ipv4 := range blockingIPv4 {
		// Try parsing as CIDR first
		_, network, err := net.ParseCIDR(ipv4)
		if err == nil && network != nil {
			config.BlockingIPv4 = append(config.BlockingIPv4, *network)
			continue
		}

		// If not CIDR, treat as single IP
		ip := net.ParseIP(ipv4)
		if ip != nil {
			// Convert single IP to /32 network
			_, network, _ := net.ParseCIDR(ipv4 + "/32")
			if network != nil {
				config.BlockingIPv4 = append(config.BlockingIPv4, *network)
			}
		}
	}

	// Parse IPv6 CIDR blocks
	config.BlockingIPv6 = make([]net.IPNet, 0)
	for _, ipv6 := range blockingIPv6 {
		// Try parsing as CIDR first
		_, network, err := net.ParseCIDR(ipv6)
		if err == nil && network != nil {
			config.BlockingIPv6 = append(config.BlockingIPv6, *network)
			continue
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

	initResolver(dnsResolver)
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

	for _, network := range config.BlockingIPv6 {
		if network.Contains(parsedIP) {
			return true
		}
	}
	return false
}

// CheckHandler responds to HTTP requests by checking DNS resolution
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	response := CheckResponse{
		Hosts: make([]HostStatus, 0, len(config.Hosts)),
	}

	allBlocked := true
	for _, host := range config.Hosts {
		status := checkHost(host)
		if !status.IsBlocked {
			allBlocked = false
		}
		response.Hosts = append(response.Hosts, status)
	}

	response.AllBlocked = allBlocked
	response.Status = "ok"

	writeResponse(w, response, allBlocked)
}

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

	// Separate IPv4 and IPv6 addresses
	for _, ip := range ips {
		if ip4 := ip.IP.To4(); ip4 != nil {
			status.IPv4 = append(status.IPv4, ip4.String())
		} else {
			status.IPv6 = append(status.IPv6, ip.IP.String())
		}
	}

	status.IsBlocked = isHostBlocked(status.IPv4, status.IPv6)
	return status
}

func isHostBlocked(ipv4 []string, ipv6 []string) bool {
	// Check IPv4 addresses
	for _, ip := range ipv4 {
		if !isIPv4Blocked(ip) {
			return false
		}
	}

	// Check IPv6 addresses
	for _, ip := range ipv6 {
		if !isIPv6Blocked(ip) {
			return false
		}
	}

	// If we get here, all IPs were blocking IPs
	return true
}

func writeResponse(w http.ResponseWriter, response CheckResponse, allBlocked bool) {
	w.Header().Set("Content-Type", "application/json")
	if !allBlocked {
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = "error"
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}
