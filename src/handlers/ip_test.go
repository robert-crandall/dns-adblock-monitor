package handlers

import (
	"net"
	"testing"
)

func TestIPv4CIDRMatching(t *testing.T) {
	tests := []struct {
		name         string
		ip           string
		blockingIPv4 []string
		want         bool
	}{
		{
			name:         "exact match with network",
			ip:           "127.0.0.1",
			blockingIPv4: []string{"127.0.0.0/8"},
			want:         true,
		},
		{
			name:         "ip within larger subnet",
			ip:           "10.0.0.5",
			blockingIPv4: []string{"10.0.0.0/24"},
			want:         true,
		},
		{
			name:         "ip outside subnet",
			ip:           "192.168.1.1",
			blockingIPv4: []string{"10.0.0.0/8"},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config = Config{
				BlockingIPv4: make([]net.IPNet, 0),
				BlockingIPv6: make([]net.IPNet, 0),
			}

			for _, cidr := range tt.blockingIPv4 {
				_, network, err := net.ParseCIDR(cidr)
				if err != nil {
					t.Fatalf("Failed to parse CIDR %s: %v", cidr, err)
				}
				config.BlockingIPv4 = append(config.BlockingIPv4, *network)
			}

			got := isIPv4Blocked(tt.ip)
			if got != tt.want {
				t.Errorf("%s: isIPv4Blocked(%s) = %v, want %v", tt.name, tt.ip, got, tt.want)
			}
		})
	}
}

func TestIPv6CIDRMatching(t *testing.T) {
	tests := []struct {
		name         string
		ip           string
		blockingIPv6 []string
		want         bool
	}{
		{
			name:         "exact match with network",
			ip:           "2a07:a8c0::",
			blockingIPv6: []string{"2a07:a8c0::/31"},
			want:         true,
		},
		{
			name:         "ip within larger subnet",
			ip:           "2a07:a8c0:4::",
			blockingIPv6: []string{"2a07:a8c0::/31"},
			want:         true,
		},
		{
			name:         "ip in second half of subnet",
			ip:           "2a07:a8c1::",
			blockingIPv6: []string{"2a07:a8c0::/31"},
			want:         true,
		},
		{
			name:         "ip outside subnet",
			ip:           "2a07:a8c2::",
			blockingIPv6: []string{"2a07:a8c0::/31"},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize config
			config = Config{
				BlockingIPv4: make([]net.IPNet, 0),
				BlockingIPv6: make([]net.IPNet, 0),
			}

			// Parse blocking networks
			for _, cidr := range tt.blockingIPv6 {
				_, network, err := net.ParseCIDR(cidr)
				if err != nil {
					t.Fatalf("Failed to parse CIDR %s: %v", cidr, err)
				}
				config.BlockingIPv6 = append(config.BlockingIPv6, *network)
			}

			// Test IP blocking
			got := isIPv6Blocked(tt.ip)
			if got != tt.want {
				t.Errorf("%s: isIPv6Blocked(%s) = %v, want %v", tt.name, tt.ip, got, tt.want)
			}
		})
	}
}
