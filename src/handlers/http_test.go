package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestResponseFormat(t *testing.T) {
	tests := []struct {
		name          string
		response      CheckResponse
		expectStatus  int
		expectBlocked bool
	}{
		{
			name: "all blocked response",
			response: CheckResponse{
				Status:     "ok",
				AllBlocked: true,
				Hosts: []HostStatus{
					{
						Host:      "ads.example.com",
						IsBlocked: true,
					},
				},
			},
			expectStatus:  200,
			expectBlocked: true,
		},
		{
			name: "partially blocked response",
			response: CheckResponse{
				Status:     "ok",
				AllBlocked: false,
				Hosts: []HostStatus{
					{
						Host:      "ads.example.com",
						IsBlocked: false,
					},
				},
			},
			expectStatus:  500,
			expectBlocked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeResponse(w, tt.response, tt.expectBlocked)

			if w.Code != tt.expectStatus {
				t.Errorf("expected status code %d, got %d", tt.expectStatus, w.Code)
			}

			var decoded CheckResponse
			if err := json.NewDecoder(w.Body).Decode(&decoded); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if decoded.AllBlocked != tt.expectBlocked {
				t.Errorf("expected AllBlocked to be %v, got %v", tt.expectBlocked, decoded.AllBlocked)
			}
		})
	}
}
