package middleware

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	rateLimitMap = make(map[string]*bucketInfo)
	rateMutex    sync.RWMutex
	maxRequests  = 60
	windowTime   = time.Minute
)

type bucketInfo struct {
	requests  int
	resetTime time.Time
	mutex     sync.Mutex
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		clientIP := getClientIP(r)
		allowed, remaining := isRequestAllowed(clientIP)

		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(maxRequests))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			log.Printf("🚦 Rate limit exceeded for IP %s", clientIP)
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error": "Rate limit exceeded", "limit": 60, "window": "1 minute"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isRequestAllowed(clientIP string) (bool, int) {
	rateMutex.Lock()
	defer rateMutex.Unlock()

	now := time.Now()

	bucket, exists := rateLimitMap[clientIP]
	if !exists {
		bucket = &bucketInfo{
			requests:  0,
			resetTime: now.Add(windowTime),
		}
		rateLimitMap[clientIP] = bucket
	}

	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()

	if now.After(bucket.resetTime) {
		bucket.requests = 0
		bucket.resetTime = now.Add(windowTime)
	}

	if bucket.requests >= maxRequests {
		return false, 0
	}

	bucket.requests++
	remaining := maxRequests - bucket.requests
	return true, remaining
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
