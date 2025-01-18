package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
}

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
    }
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr

        rl.mu.Lock()
        now := time.Now()
        
        // Clean old requests
        if times, exists := rl.requests[ip]; exists {
            var valid []time.Time
            for _, t := range times {
                if now.Sub(t) < time.Minute {
                    valid = append(valid, t)
                }
            }
            rl.requests[ip] = valid
        }

        // Check rate limit (5 requests per minute)
        if len(rl.requests[ip]) >= 5 {
            rl.mu.Unlock()
            w.Write([]byte(`<div class="error">Too many attempts. Please try again later.</div>`))
            return
        }

        // Add new request
        rl.requests[ip] = append(rl.requests[ip], now)
        rl.mu.Unlock()

        next(w, r)
    }
} 