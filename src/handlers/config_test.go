package handlers

import (
	"net"
	"testing"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name         string
		hosts        []string
		blockingIPv4 []string
		blockingIPv6 []string
		resolver     string
		wantHosts    int
		wantIPv4     int
		wantIPv6     int
	}{
		{
			name:         "empty config",
			hosts:        []string{},
			blockingIPv4: []string{},
			blockingIPv6: []string{},
			resolver:     "",
			wantHosts:    0,
			wantIPv4:     0,
			wantIPv6:     0,
		},
		{
			name:         "basic config",
			hosts:        []string{"ads.example.com"},
			blockingIPv4: []string{"0.0.0.0/8"},
			blockingIPv6: []string{"::/128"},
			resolver:     "1.1.1.1:53",
			wantHosts:    1,
			wantIPv4:     1,
			wantIPv6:     1,
		},
		{
			name:         "multiple blocks",
			hosts:        []string{"ads1.example.com", "ads2.example.com"},
			blockingIPv4: []string{"0.0.0.0/8", "127.0.0.0/8"},
			blockingIPv6: []string{"::/128", "fc00::/7"},
			resolver:     "8.8.8.8:53",
			wantHosts:    2,
			wantIPv4:     2,
			wantIPv6:     2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Initialize(tt.hosts, tt.blockingIPv4, tt.blockingIPv6, tt.resolver)

			// Check hosts
			if len(config.Hosts) != tt.wantHosts {
				t.Errorf("Initialize() hosts = %v, want %v", len(config.Hosts), tt.wantHosts)
			}

			// Check IPv4 blocks
			if len(config.BlockingIPv4) != tt.wantIPv4 {
				t.Errorf("Initialize() IPv4 blocks = %v, want %v", len(config.BlockingIPv4), tt.wantIPv4)
			}

			// Check IPv6 blocks
			if len(config.BlockingIPv6) != tt.wantIPv6 {
				t.Errorf("Initialize() IPv6 blocks = %v, want %v", len(config.BlockingIPv6), tt.wantIPv6)
			}
		})
	}
}

func TestParseIPv4Blocks(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    int
		wantErr bool
	}{
		{
			name:    "valid CIDR",
			input:   []string{"192.168.0.0/24"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "valid single IP",
			input:   []string{"127.0.0.1"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "mixed valid inputs",
			input:   []string{"10.0.0.0/8", "192.168.1.1"},
			want:    2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.BlockingIPv4 = make([]net.IPNet, 0)
			parseIPv4Blocks(tt.input)

			if len(config.BlockingIPv4) != tt.want {
				t.Errorf("parseIPv4Blocks() got %v networks, want %v", len(config.BlockingIPv4), tt.want)
			}
		})
	}
}

func TestParseIPv6Blocks(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    int
		wantErr bool
	}{
		{
			name:    "valid CIDR",
			input:   []string{"2001:db8::/32"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "valid single IP",
			input:   []string{"::1"},
			want:    1,
			wantErr: false,
		},
		{
			name:    "mixed valid inputs",
			input:   []string{"2001:db8::/32", "::1"},
			want:    2,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.BlockingIPv6 = make([]net.IPNet, 0)
			parseIPv6Blocks(tt.input)

			if len(config.BlockingIPv6) != tt.want {
				t.Errorf("parseIPv6Blocks() got %v networks, want %v", len(config.BlockingIPv6), tt.want)
			}
		})
	}
}
