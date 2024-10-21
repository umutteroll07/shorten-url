package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/umutteroll07/ShortenURL/app/helper"
)



var ctx = context.Background()

func RateLimiter(c *fiber.Ctx) error {
	
		ttl := time.Second * 60
		limit := 5

		userIP := c.IP()
		key := "rate:" + userIP
		rdb := helper.RedisClient()

		requestCount, err := rdb.Get(ctx, key).Int()
		if err == redis.Nil {
			requestCount = 0
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Redis error" + err.Error())
		}

		if requestCount >= limit {
			return c.Status(fiber.StatusTooManyRequests).SendString("Too many requests")
		}

		rdb.Incr(ctx, key)
		rdb.Expire(ctx, key, ttl)


		
	return c.Next()

}
