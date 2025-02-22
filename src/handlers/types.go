package handlers

import "net"

type Config struct {
	Hosts        []string
	BlockingIPv4 []net.IPNet
	BlockingIPv6 []net.IPNet
	Resolver     string
}

type HostStatus struct {
	Host          string   `json:"host"`
	IPv4          []string `json:"ipv4,omitempty"`
	IPv6          []string `json:"ipv6,omitempty"`
	UnblockedIPv4 []string `json:"unblocked_ipv4,omitempty"`
	UnblockedIPv6 []string `json:"unblocked_ipv6,omitempty"`
	Error         string   `json:"error,omitempty"`
	IsBlocked     bool     `json:"is_blocked"`
}

type CheckResponse struct {
	Status         string `json:"status"`
	AllBlocked     bool   `json:"all_blocked"`
	BlockingRanges struct {
		IPv4 []string `json:"ipv4"`
		IPv6 []string `json:"ipv6"`
	} `json:"blocking_ranges"`
	Hosts []HostStatus `json:"hosts"`
}
