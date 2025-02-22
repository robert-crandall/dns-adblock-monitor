package handlers

import (
	"context"
	"net"
	"reflect"
	"testing"
)

func normalizeHostStatus(h HostStatus) HostStatus {
	// Helper to normalize slice
	normalizeSlice := func(s []string) []string {
		if s == nil || len(s) == 0 {
			return []string{}
		}
		return s
	}

	// Create a copy with normalized slices
	return HostStatus{
		Host:          h.Host,
		IPv4:          normalizeSlice(h.IPv4),
		IPv6:          normalizeSlice(h.IPv6),
		UnblockedIPv4: normalizeSlice(h.UnblockedIPv4),
		UnblockedIPv6: normalizeSlice(h.UnblockedIPv6),
		IsBlocked:     h.IsBlocked,
		Error:         h.Error,
	}
}

func TestCheckHost(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		blockingIPv4 []string
		blockingIPv6 []string
		mockIPs      []net.IPAddr
		mockErr      error
		want         HostStatus
	}{
		{
			name:         "all IPs blocked",
			host:         "blocked.example.com",
			blockingIPv4: []string{"127.0.0.0/8"},
			blockingIPv6: []string{"::1/128"},
			mockIPs: []net.IPAddr{
				{IP: net.ParseIP("127.0.0.1")},
				{IP: net.ParseIP("::1")},
			},
			want: HostStatus{
				Host:          "blocked.example.com",
				IPv4:          []string{"127.0.0.1"},
				IPv6:          []string{"::1"},
				UnblockedIPv4: []string{},
				UnblockedIPv6: []string{},
				IsBlocked:     true,
				Error:         "",
			},
		},
		{
			name:         "mixed blocking",
			host:         "mixed.example.com",
			blockingIPv4: []string{"127.0.0.0/8"},
			blockingIPv6: []string{"::1/128"},
			mockIPs: []net.IPAddr{
				{IP: net.ParseIP("127.0.0.1")},
				{IP: net.ParseIP("8.8.8.8")},
				{IP: net.ParseIP("2001:4860:4860::8888")},
			},
			want: HostStatus{
				Host:          "mixed.example.com",
				IPv4:          []string{"127.0.0.1", "8.8.8.8"},
				IPv6:          []string{"2001:4860:4860::8888"},
				UnblockedIPv4: []string{"8.8.8.8"},
				UnblockedIPv6: []string{"2001:4860:4860::8888"},
				IsBlocked:     false,
				Error:         "",
			},
		},
		{
			name:         "lookup error",
			host:         "error.example.com",
			blockingIPv4: []string{},
			blockingIPv6: []string{},
			mockErr:      &net.DNSError{Err: "no such host", Name: "error.example.com"},
			want: HostStatus{
				Host:          "error.example.com",
				IPv4:          []string{},
				IPv6:          []string{},
				UnblockedIPv4: []string{},
				UnblockedIPv6: []string{},
				IsBlocked:     true,
				Error:         "lookup error.example.com: no such host",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock resolver with test instance
			oldResolver := resolver

			mock := &mockResolver{
				t:   t,
				ips: tt.mockIPs,
				err: tt.mockErr,
			}
			resolver = mock

			defer func() {
				resolver = oldResolver
			}()

			// Initialize config with test blocking ranges
			Initialize([]string{tt.host}, tt.blockingIPv4, tt.blockingIPv6, "")

			// Run the test with debug logging
			got := checkHost(tt.host)

			// Compare results
			if !reflect.DeepEqual(normalizeHostStatus(got), normalizeHostStatus(tt.want)) {
				t.Errorf("checkHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Mock resolver for testing
type mockResolver struct {
	t   *testing.T // Add testing.T to access test logging
	ips []net.IPAddr
	err error
}

func (r *mockResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.ips, nil
}
