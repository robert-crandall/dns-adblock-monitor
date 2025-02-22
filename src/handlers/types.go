package handlers

type Config struct {
	Hosts        []string
	BlockingIPv4 []string
	BlockingIPv6 []string
	Resolver     string
}

type HostStatus struct {
	Host      string   `json:"host"`
	IPv4      []string `json:"ipv4,omitempty"`
	IPv6      []string `json:"ipv6,omitempty"`
	Error     string   `json:"error,omitempty"`
	IsBlocked bool     `json:"is_blocked"`
}

type CheckResponse struct {
	Status     string       `json:"status"`
	AllBlocked bool         `json:"all_blocked"`
	Hosts      []HostStatus `json:"hosts"`
}
