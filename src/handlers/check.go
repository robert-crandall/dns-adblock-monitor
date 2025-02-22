package handlers

import (
	"net/http"

	"github.com/robert-crandall/dns-adblock-monitor/src/dns"
)

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	hostname := r.URL.Query().Get("host")
	if hostname == "" {
		http.Error(w, "Host parameter is required", http.StatusBadRequest)
		return
	}

	err := dns.ResolveDNS(hostname)
	if err != nil {
		http.Error(w, "DNS resolution failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
