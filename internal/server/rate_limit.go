package server

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/JavierAnte/local-offers-api/internal/httpx"
)

type rateLimitWindow struct {
	count   int
	resetAt time.Time
}

func rateLimitByIP(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]rateLimitWindow)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			client := requestIP(r)

			mu.Lock()
			for key, current := range clients {
				if !now.Before(current.resetAt) {
					delete(clients, key)
				}
			}

			current, ok := clients[client]
			if !ok {
				current = rateLimitWindow{resetAt: now.Add(window)}
			}
			if current.count >= maxRequests {
				retryAfter := max(1, int(time.Until(current.resetAt).Seconds()))
				mu.Unlock()
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				httpx.WriteError(w, http.StatusTooManyRequests, "rate_limit_exceeded", "Too many requests. Please try again later.")
				return
			}

			current.count++
			clients[client] = current
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

func requestIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
