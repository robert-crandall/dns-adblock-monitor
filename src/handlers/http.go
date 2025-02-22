package handlers

import (
	"encoding/json"
	"net/http"
)

// CheckHandler responds to HTTP requests by checking DNS resolution
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	response := CheckResponse{
		Hosts: make([]HostStatus, 0, len(config.Hosts)),
	}

	// Add blocking ranges to response
	for _, network := range config.BlockingIPv4 {
		response.BlockingRanges.IPv4 = append(response.BlockingRanges.IPv4, network.String())
	}
	for _, network := range config.BlockingIPv6 {
		response.BlockingRanges.IPv6 = append(response.BlockingRanges.IPv6, network.String())
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
