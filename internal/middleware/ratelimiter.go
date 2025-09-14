package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/checkspeed/sc-backend/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ClientLimiter holds rate limiters for each client (IP address)
type ClientLimiter struct {
	ips   map[string]*rate.Limiter
	mutex sync.RWMutex
	limit rate.Limit
	burst int
}

// NewClientLimiter creates a new client limiter
// For 1 request per minute: rate.Every(time.Minute) with burst of 1
func NewClientLimiter() *ClientLimiter {
	return &ClientLimiter{
		ips:   make(map[string]*rate.Limiter),
		limit: rate.Every(time.Minute), // 1 request per minute
		burst: 1,
	}
}

// GetLimiter returns the rate limiter for a specific IP
func (cl *ClientLimiter) GetLimiter(ip string) *rate.Limiter {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	limiter, exists := cl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.limit, cl.burst)
		cl.ips[ip] = limiter
	}
	return limiter

}

// CleanupOldLimiters removes limiters that haven't been used recently
// Call this periodically to prevent memory leaks
func (cl *ClientLimiter) CleanupStaleIPs() {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	toDelete := make([]string, 0)

	for ip, limiter := range cl.ips {
		// If limiter has full tokens available, it hasn't been used recently
		// We can safely remove it to free up memory
		if limiter.Tokens() >= float64(cl.burst) {
			toDelete = append(toDelete, ip)
		}
	}
	// Delete the unused limiters
	for _, ip := range toDelete {
		delete(cl.ips, ip)
	}
}

// RateLimit middleware returns a Gin middleware that implements rate limiting
func RateLimit(clientLimiter *ClientLimiter) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Get client IP address
		ip := c.ClientIP()

		// Get the rate limiter for this IP
		limiter := clientLimiter.GetLimiter(ip)

		// check if reuqest is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.ApiResp{
				Status:  models.StatusError,
				Message: "Rate limit exceeded. Please wait 1 minute between requests",
				Code:    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		// continue to next handler
		c.Next()
	}
}
