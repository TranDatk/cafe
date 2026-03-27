package middleware

import (
	"cafe/domain"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

var clients sync.Map // Map of ip (string) -> *client

func init() {
	// Periodic cleanup of inactive clients every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			clients.Range(func(key, value interface{}) bool {
				ip := key.(string)
				c := value.(*client)
				
				c.mu.Lock()
				isInactive := time.Since(c.lastSeen) > 3*time.Minute
				c.mu.Unlock()

				if isInactive {
					clients.Delete(ip)
				}
				return true
			})
		}
	}()
}

// RateLimitMiddleware limits the number of requests per IP
func RateLimitMiddleware(requestsPerSecond float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Load or store a new client for the IP
		val, _ := clients.LoadOrStore(ip, &client{
			limiter:  rate.NewLimiter(rate.Limit(requestsPerSecond), burst),
			lastSeen: time.Now(),
		})
		
		client := val.(*client)

		// Update last seen time under client lock
		client.mu.Lock()
		client.lastSeen = time.Now()
		limiter := client.limiter
		client.mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, domain.ErrorResponse{
				Message: "Rate limit exceeded. Please try again later.",
			})
			return
		}

		c.Next()
	}
}
