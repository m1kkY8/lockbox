package config

import (
	"testing"
)

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected string
	}{
		{
			name: "basic config",
			config: Config{
				Host: "localhost:1337",
			},
			expected: "ws://localhost:1337/chat",
		},
		{
			name: "secure config",
			config: Config{
				Host:   "example.com:443",
				Secure: "true",
			},
			expected: "wss://example.com:443/chat",
		},
		{
			name: "host without port",
			config: Config{
				Host: "localhost",
			},
			expected: "ws://localhost:1337/chat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := GetUrl(tt.config)
			if url.String() != tt.expected {
				t.Errorf("GetUrl() = %v, want %v", url.String(), tt.expected)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				Username: "testuser",
				Host:     "localhost:1337",
				Color:    "93",
			},
			expectError: false,
		},
		{
			name:        "empty config",
			config:      Config{},
			expectError: true,
		},
		{
			name: "missing host",
			config: Config{
				Username: "testuser",
				Color:    "93",
			},
			expectError: true,
		},
		{
			name: "help color command",
			config: Config{
				Username: "testuser",
				Host:     "localhost:1337",
				Color:    "help",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateConfig() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestApplyDefaults(t *testing.T) {
	config := &Config{
		Username: "testuser",
		Host:     "localhost",
		// Color is empty, should get default
	}

	ApplyDefaults(config)

	if config.Color == "" {
		t.Error("ApplyDefaults() should set a default color")
	}

	// Verify other fields are unchanged
	if config.Username != "testuser" {
		t.Errorf("ApplyDefaults() should not change username, got %v", config.Username)
	}
	if config.Host != "localhost" {
		t.Errorf("ApplyDefaults() should not change host, got %v", config.Host)
	}
}

// Benchmark test for URL construction
func BenchmarkGetUrl(b *testing.B) {
	config := Config{
		Host:   "localhost:1337",
		Secure: "",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetUrl(config)
	}
}
