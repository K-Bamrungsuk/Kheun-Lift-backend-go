package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit() gin.HandlerFunc {
	const (
		limit  = 5
		window = time.Minute
		maxIPs = 10_000
	)

	type entry struct {
		count   int
		resetAt time.Time
	}

	// shortcut: limits are local to this process; use shared storage when scaling to multiple instances.
	entries := make(map[string]entry)
	var mu sync.Mutex
	var nextCleanup time.Time

	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		ip := c.ClientIP()

		mu.Lock()
		now := time.Now()
		if !now.Before(nextCleanup) {
			for key, value := range entries {
				if !now.Before(value.resetAt) {
					delete(entries, key)
				}
			}
			nextCleanup = now.Add(window)
		}

		current := entries[ip]
		if !now.Before(current.resetAt) {
			delete(entries, ip)
			current = entry{resetAt: now.Add(window)}
		}

		retryAfter := 0
		switch {
		case current.count >= limit:
			retryAfter = int((current.resetAt.Sub(now) + time.Second - 1) / time.Second)
		case current.count == 0 && len(entries) >= maxIPs:
			retryAfter = int(window / time.Second)
		default:
			current.count++
			entries[ip] = current
		}
		mu.Unlock()

		if retryAfter > 0 {
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
