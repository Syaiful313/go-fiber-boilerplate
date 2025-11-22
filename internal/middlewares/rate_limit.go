package middlewares

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type rateLimiter struct {
	mu          sync.Mutex
	requests    map[string]*rateRecord
	maxRequests int
	window      time.Duration
}

type rateRecord struct {
	count     int
	expiresAt time.Time
}

func newRateLimiter(max int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		requests:    make(map[string]*rateRecord),
		maxRequests: max,
		window:      window,
	}
}

func (r *rateLimiter) allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	rec, exists := r.requests[key]
	if !exists || now.After(rec.expiresAt) {
		if exists {
			delete(r.requests, key)
		}
		r.requests[key] = &rateRecord{
			count:     1,
			expiresAt: now.Add(r.window),
		}
		return true
	}

	if rec.count >= r.maxRequests {
		return false
	}

	rec.count++
	return true
}

func RateLimitMiddleware(max int, window time.Duration, keyFunc func(*fiber.Ctx) string) fiber.Handler {
	rl := newRateLimiter(max, window)

	return func(c *fiber.Ctx) error {
		key := keyFunc(c)
		if key == "" {
			key = c.IP()
		}

		if !rl.allow(key) {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}

		return c.Next()
	}
}
