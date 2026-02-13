package middleware

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-auth/internal/core/ports"
	"github.com/redis/go-redis/v9"
)

func IPRateLimiter(logger ports.Logger, rdb *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// Get the IP (simplified)
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr // fallback
			}
			key := "rl:" + ip

			logger.Info(ip)
			// 1. Atomically increment the counter
			// If the key doesn't exist, Redis creates it with value 1
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				// If Redis is down, we log it and allow the request (Fail Open)
				next.ServeHTTP(w, r)
				return
			}

			// 2. If it's a brand new key (count == 1), set the expiration
			if count == 1 {
				rdb.Expire(ctx, key, window)
			}

			// 3. Check the limit
			if int(count) > limit {
				// Return 429 Too Many Requests
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("Slow down! You've hit the limit."))
				return
			}

			// Allow the request to proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}
