package handlers

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

// Config holds the handler configuration
type Config struct {
	Hosts       []string
	BlockingIPs []string
	Resolver    string
}
