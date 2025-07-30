package auth

import (
	"sync"
	"time"
)

// RateLimiter implements a simple rate limiting mechanism
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// IsAllowed checks if a request from the given IP is allowed
func (r *RateLimiter) IsAllowed(ip string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()

	// Get existing requests for this IP
	requests, exists := r.requests[ip]
	if !exists {
		r.requests[ip] = []time.Time{now}
		return true
	}

	// Remove old requests outside the window
	var validRequests []time.Time
	for _, requestTime := range requests {
		if now.Sub(requestTime) <= r.window {
			validRequests = append(validRequests, requestTime)
		}
	}

	// Check if we're within the limit
	if len(validRequests) >= r.limit {
		r.requests[ip] = validRequests
		return false
	}

	// Add current request
	validRequests = append(validRequests, now)
	r.requests[ip] = validRequests
	return true
}

// Clean removes old entries to prevent memory leaks
func (r *RateLimiter) Clean() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	for ip, requests := range r.requests {
		var validRequests []time.Time
		for _, requestTime := range requests {
			if now.Sub(requestTime) <= r.window {
				validRequests = append(validRequests, requestTime)
			}
		}

		if len(validRequests) == 0 {
			delete(r.requests, ip)
		} else {
			r.requests[ip] = validRequests
		}
	}
}
