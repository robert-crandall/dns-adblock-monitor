package handlers

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

// HostStatus represents the resolution status of a single host
type HostStatus struct {
	Host      string   `json:"host"`
	IPs       []string `json:"ips,omitempty"`
	Error     string   `json:"error,omitempty"`
	IsBlocked bool     `json:"is_blocked"`
}

// CheckResponse represents the complete check response
type CheckResponse struct {
	Status     string       `json:"status"`
	AllBlocked bool         `json:"all_blocked"`
	Hosts      []HostStatus `json:"hosts"`
}

type Config struct {
	Hosts       []string
	BlockingIPs []string
	Resolver    string
}

var (
	config   Config
	resolver *net.Resolver
	mu       sync.RWMutex
)

// Initialize sets up the handlers package with configuration
func Initialize(hosts []string, blockingIPs []string, dnsResolver string) {
	mu.Lock()
	defer mu.Unlock()

	config.Hosts = hosts

	if len(blockingIPs) > 0 {
		config.BlockingIPs = blockingIPs
	} else {
		// Default blocking IPs if none provided
		config.BlockingIPs = []string{"0.0.0.0", "127.0.0.1"}
	}

	if dnsResolver != "" {
		resolver = &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: time.Second * 10,
				}
				return d.DialContext(ctx, "udp", dnsResolver)
			},
		}
	} else {
		resolver = net.DefaultResolver
	}
}

// CheckHandler responds to HTTP requests by checking DNS resolution for configured hosts
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	response := CheckResponse{
		Hosts: make([]HostStatus, 0, len(config.Hosts)),
	}

	allBlocked := true
	for _, host := range config.Hosts {
		status := HostStatus{
			Host: host,
		}

		ips, err := resolver.LookupHost(context.Background(), host)
		if err != nil {
			status.Error = err.Error()
			status.IsBlocked = true
			response.Hosts = append(response.Hosts, status)
			continue
		}

		status.IPs = ips
		status.IsBlocked = true

		// Check if all returned IPs are blocking IPs
		for _, ip := range ips {
			isBlocking := false
			for _, expectedIP := range config.BlockingIPs {
				if ip == expectedIP {
					isBlocking = true
					break
				}
			}
			if !isBlocking {
				status.IsBlocked = false
				allBlocked = false
			}
		}

		response.Hosts = append(response.Hosts, status)
	}

	response.AllBlocked = allBlocked
	response.Status = "ok"

	w.Header().Set("Content-Type", "application/json")
	if !allBlocked {
		w.WriteHeader(http.StatusInternalServerError)
		response.Status = "error"
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
