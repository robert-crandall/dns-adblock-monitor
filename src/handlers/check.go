package handlers

import (
	"context"
	"encoding/json"
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

	if len(blockingIPv4) > 0 {
		config.BlockingIPv4 = blockingIPv4
	} else {
		config.BlockingIPv4 = []string{"0.0.0.0", "127.0.0.1"}
	}

	if len(blockingIPv6) > 0 {
		config.BlockingIPv6 = blockingIPv6
	} else {
		config.BlockingIPv6 = []string{"::", "::1"}
	}

	initResolver(dnsResolver)
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
		isBlocking := false
		for _, expectedIP := range config.BlockingIPv4 {
			if ip == expectedIP {
				isBlocking = true
				break
			}
		}
		if !isBlocking {
			return false
		}
	}

	// Check IPv6 addresses
	for _, ip := range ipv6 {
		isBlocking := false
		for _, expectedIP := range config.BlockingIPv6 {
			if ip == expectedIP {
				isBlocking = true
				break
			}
		}
		if !isBlocking {
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
