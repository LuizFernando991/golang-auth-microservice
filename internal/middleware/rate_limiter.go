package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/LuizFernando991/golang-auth-microservice/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisRateLimiter struct {
	Client *redis.Client
	Logger *config.Logger
	Limit  int
	Window time.Duration
}

func NewRedisRateLimiter(client *redis.Client, logger *config.Logger, limit int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{Client: client, Logger: logger, Limit: limit, Window: window}
}

func (r *RedisRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := "rl:" + ip + ":" + c.FullPath()
		res, err := r.Client.Incr(ctx, key).Result()
		if err != nil {
			r.Logger.Warn("redis incr failed", zap.Error(err))
			c.Next()
			return
		}
		if res == 1 {
			_ = r.Client.Expire(ctx, key, r.Window).Err()
		}
		if res > int64(r.Limit) {
			ttl, _ := r.Client.TTL(ctx, key).Result()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded", "retry_after_seconds": int(ttl.Seconds())})
			return
		}
		c.Next()
	}
}
