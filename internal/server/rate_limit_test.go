package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRateLimitByIP(t *testing.T) {
	t.Parallel()

	handler := rateLimitByIP(2, time.Hour)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for attempt := 1; attempt <= 3; attempt++ {
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.RemoteAddr = "203.0.113.10:1234"
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)

		if attempt <= 2 && recorder.Code != http.StatusNoContent {
			t.Fatalf("attempt %d status = %d, want %d", attempt, recorder.Code, http.StatusNoContent)
		}
		if attempt == 3 {
			if recorder.Code != http.StatusTooManyRequests {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusTooManyRequests)
			}
			if !strings.Contains(recorder.Body.String(), `"code":"rate_limit_exceeded"`) {
				t.Fatalf("body = %s", recorder.Body.String())
			}
			if recorder.Header().Get("Retry-After") == "" {
				t.Fatal("Retry-After header is missing")
			}
		}
	}
}

func TestRateLimitByIPSeparatesClients(t *testing.T) {
	t.Parallel()

	handler := rateLimitByIP(1, time.Hour)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, remoteAddr := range []string{"203.0.113.10:1234", "203.0.113.11:1234"} {
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.RemoteAddr = remoteAddr
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusNoContent {
			t.Fatalf("client %s status = %d", remoteAddr, recorder.Code)
		}
	}
}
