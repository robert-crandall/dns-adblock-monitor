package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckHandler(t *testing.T) {
	tests := []struct {
		name          string
		hosts         []string
		blockingIPv4  []string
		blockingIPv6  []string
		resolver      string
		expectedCode  int
		expectBlocked bool
	}{
		{
			name:          "all hosts blocked",
			hosts:         []string{"blocked.example.com"},
			blockingIPv4:  []string{"0.0.0.0"},
			blockingIPv6:  []string{"::"},
			expectedCode:  http.StatusOK,
			expectBlocked: true,
		},
		{
			name:          "host resolves to non-blocking IP",
			hosts:         []string{"google.com"},
			blockingIPv4:  []string{"0.0.0.0"},
			blockingIPv6:  []string{"::"},
			expectedCode:  http.StatusInternalServerError,
			expectBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize handler with test config
			Initialize(tt.hosts, tt.blockingIPv4, tt.blockingIPv6, tt.resolver)

			// Create test request
			req := httptest.NewRequest("GET", "/check", nil)
			w := httptest.NewRecorder()

			// Call handler
			CheckHandler(w, req)

			// Check status code
			if w.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			// Parse response
			var response CheckResponse
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			// Check blocking status
			if response.AllBlocked != tt.expectBlocked {
				t.Errorf("expected AllBlocked to be %v, got %v", tt.expectBlocked, response.AllBlocked)
			}
		})
	}
}

func TestIsHostBlocked(t *testing.T) {
	tests := []struct {
		name         string
		ipv4         []string
		ipv6         []string
		blockingIPv4 []string
		blockingIPv6 []string
		want         bool
	}{
		{
			name:         "all IPs blocked",
			ipv4:         []string{"0.0.0.0"},
			ipv6:         []string{"::"},
			blockingIPv4: []string{"0.0.0.0"},
			blockingIPv6: []string{"::"},
			want:         true,
		},
		{
			name:         "non-blocking IPv4",
			ipv4:         []string{"8.8.8.8"},
			ipv6:         []string{"::"},
			blockingIPv4: []string{"0.0.0.0"},
			blockingIPv6: []string{"::"},
			want:         false,
		},
		{
			name:         "non-blocking IPv6",
			ipv4:         []string{"0.0.0.0"},
			ipv6:         []string{"2001:4860:4860::8888"},
			blockingIPv4: []string{"0.0.0.0"},
			blockingIPv6: []string{"::"},
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up config for test
			config.BlockingIPv4 = tt.blockingIPv4
			config.BlockingIPv6 = tt.blockingIPv6

			if got := isHostBlocked(tt.ipv4, tt.ipv6); got != tt.want {
				t.Errorf("isHostBlocked() = %v, want %v", got, tt.want)
			}
		})
	}
}
