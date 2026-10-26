package config

import "testing"

func TestValidatePublicBaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    string
		wantErr bool
	}{
		{"https origin", "https://api.example.com", "https://api.example.com", false},
		{"local http origin", "http://localhost:8080/", "http://localhost:8080", false},
		{"missing", "", "", true},
		{"relative", "api.example.com", "", true},
		{"unsupported scheme", "ftp://api.example.com", "", true},
		{"path", "https://api.example.com/api/v1", "", true},
		{"query", "https://api.example.com?source=test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := validatePublicBaseURL(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validatePublicBaseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("validatePublicBaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
