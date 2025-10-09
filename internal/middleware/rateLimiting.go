package middleware

import (
	"net/http"
	"sync"

	"github.com/alex-cos/scvtemplate/observability"
	"github.com/gin-gonic/gin"
	ratelimit "golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors map[string]*ratelimit.Limiter
	mu       sync.Mutex
	rate     float64
	burst    int
}

func NewRateLimiter(rate float64, burst int) *RateLimiter {
	r := &RateLimiter{
		visitors: make(map[string]*ratelimit.Limiter),
		mu:       sync.Mutex{},
		rate:     rate,
		burst:    burst,
	}

	return r
}

func (r *RateLimiter) getVisitor(ip string) *ratelimit.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	limiter, exists := r.visitors[ip]
	if !exists {
		limiter = ratelimit.NewLimiter(ratelimit.Limit(r.rate), r.burst)
		r.visitors[ip] = limiter
	}
	return limiter
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := r.getVisitor(ip)
		if !limiter.Allow() {
			path := c.FullPath()
			method := c.Request.Method
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			observability.RateLimitBlockedTotal.WithLabelValues(method, path).Inc()
			return
		}
		c.Next()
	}
}
