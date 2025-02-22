package handlers

import (
	"net"
	"net/http"
	"sync"
)

type Config struct {
	Hosts       []string
	BlockingIPs []string
}

var (
	config Config
	mu     sync.RWMutex
)

// Initialize sets up the handlers package with configuration
func Initialize(hosts []string, blockingIPs ...string) {
	mu.Lock()
	defer mu.Unlock()

	config.Hosts = hosts
	if len(blockingIPs) > 0 {
		config.BlockingIPs = blockingIPs
	} else {
		// Default blocking IPs if none provided
		config.BlockingIPs = []string{"0.0.0.0", "127.0.0.1"}
	}
}

// CheckHandler responds to HTTP requests by checking DNS resolution for configured hosts
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	for _, host := range config.Hosts {
		ips, err := net.LookupHost(host)
		if err != nil {
			// DNS resolution failed (good - it's blocked)
			continue
		}

		// If we got IPs back, check if they're all blocking IPs
		for _, ip := range ips {
			isBlocking := false
			for _, blockingIP := range config.BlockingIPs {
				if ip == blockingIP {
					isBlocking = true
					break
				}
			}
			if !isBlocking {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	// If we get here, all hosts were either unresolvable or returned blocking IPs
	w.WriteHeader(http.StatusOK)
}
