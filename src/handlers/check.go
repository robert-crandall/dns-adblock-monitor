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
func Initialize(hosts []string, blockingIPs []string, dnsResolver string) {
	mu.Lock()
	defer mu.Unlock()

	config.Hosts = hosts

	if len(blockingIPs) > 0 {
		config.BlockingIPs = blockingIPs
	} else {
		config.BlockingIPs = []string{"0.0.0.0", "127.0.0.1"}
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

	ips, err := resolver.LookupHost(context.Background(), host)
	if err != nil {
		status.Error = err.Error()
		status.IsBlocked = true
		return status
	}

	status.IPs = ips
	status.IsBlocked = isHostBlocked(ips)

	return status
}

func isHostBlocked(ips []string) bool {
	for _, ip := range ips {
		isBlocking := false
		for _, expectedIP := range config.BlockingIPs {
			if ip == expectedIP {
				isBlocking = true
				break
			}
		}
		if !isBlocking {
			return false
		}
	}
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
