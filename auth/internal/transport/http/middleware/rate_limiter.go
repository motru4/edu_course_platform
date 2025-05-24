package middleware

import (
	"net/http"
	"sync"
	"time"

	"auth-service/internal/config"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	tokens    map[string][]time.Time
	mu        sync.Mutex
	maxTokens int
	period    time.Duration
}

func NewRateLimiter(cfg config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		tokens:    make(map[string][]time.Time),
		maxTokens: cfg.Requests,
		period:    cfg.Period,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()

		// Очищаем устаревшие токены
		if timestamps, exists := rl.tokens[clientIP]; exists {
			var validTokens []time.Time
			for _, ts := range timestamps {
				if now.Sub(ts) <= rl.period {
					validTokens = append(validTokens, ts)
				}
			}
			rl.tokens[clientIP] = validTokens
		}

		// Проверяем количество запросов
		if len(rl.tokens[clientIP]) >= rl.maxTokens {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"details": "please try again later",
			})
			c.Abort()
			return
		}

		// Добавляем новый токен
		rl.tokens[clientIP] = append(rl.tokens[clientIP], now)

		c.Next()
	}
}
