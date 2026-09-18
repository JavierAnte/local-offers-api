package httpx

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, body string
		wantErr    bool
	}{
		{"valid", `{"name":"local"}`, false},
		{"unknown field", `{"name":"local","extra":true}`, true},
		{"trailing value", `{"name":"local"} {}`, true},
		{"malformed", `{"name":`, true},
		{"empty", ``, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest("POST", "/", strings.NewReader(tt.body))
			var dst struct {
				Name string `json:"name"`
			}
			err := DecodeJSON(httptest.NewRecorder(), req, &dst)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DecodeJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriteError(t *testing.T) {
	t.Parallel()
	recorder := httptest.NewRecorder()
	WriteError(recorder, 400, "invalid_input", "Invalid input.")
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("content type = %q", recorder.Header().Get("Content-Type"))
	}
	if recorder.Body.String() != "{\"error\":{\"code\":\"invalid_input\",\"message\":\"Invalid input.\"}}\n" {
		t.Fatalf("body = %s", recorder.Body.String())
	}
}
