package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"time"
)

// RateLimiter stores request information
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter creates a new rate limiter
// limit - maximum number of requests
// window - time window
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine for old entries
	go rl.cleanup()

	return rl
}

// cleanup removes old entries from memory
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()

		for key, timestamps := range rl.requests {
			// Remove entries older than 2*window
			cutoff := now.Add(-2 * rl.window)
			valid := make([]time.Time, 0)

			for _, ts := range timestamps {
				if ts.After(cutoff) {
					valid = append(valid, ts)
				}
			}

			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}

		rl.mu.Unlock()
	}
}

// isAllowed checks if request is allowed
func (rl *RateLimiter) isAllowed(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Get or create timestamps list for this key
	timestamps := rl.requests[key]

	// Keep only requests from current window
	validTimestamps := make([]time.Time, 0)
	for _, ts := range timestamps {
		if ts.After(windowStart) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	// Check limit
	if len(validTimestamps) >= rl.limit {
		return false
	}

	// Add current request
	validTimestamps = append(validTimestamps, now)
	rl.requests[key] = validTimestamps

	return true
}

// getRateLimitKey returns rate limiting key
// Priority: User ID (from JWT) > IP address
func getRateLimitKey(c *gin.Context) string {
	// Check for userID in context (from JWT)
	if userID, exists := c.Get("userID"); exists {
		return fmt.Sprintf("user:%v", userID)
	}

	// If no JWT, use IP address
	return fmt.Sprintf("ip:%s", c.ClientIP())
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute, time.Minute)

	return func(c *gin.Context) {
		key := getRateLimitKey(c)

		if !limiter.isAllowed(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded: maximum %d requests per minute", requestsPerMinute),
			})
			c.Abort()
			return
		}

		// Add headers with limit information
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerMinute))
		// In production app you can add remaining requests count

		c.Next()
	}
}
