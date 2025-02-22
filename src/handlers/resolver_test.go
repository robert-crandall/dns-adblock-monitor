package handlers

import (
	"context"
	"testing"
)

func TestResolverInitialization(t *testing.T) {
	tests := []struct {
		name         string
		resolverAddr string
		expectCustom bool
		expectError  bool
	}{
		{
			name:         "system resolver",
			resolverAddr: "",
			expectCustom: false,
			expectError:  false,
		},
		{
			name:         "custom resolver",
			resolverAddr: "1.1.1.1:53",
			expectCustom: true,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initResolver(tt.resolverAddr)

			_, err := resolver.LookupIPAddr(context.Background(), "example.com")
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}
