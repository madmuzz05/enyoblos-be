package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/madmuzz05/be-enyoblos/config"
	"github.com/madmuzz05/be-enyoblos/package/helper"
)

func NewRateLimiter() fiber.Handler {
	// use env or passed config; simple example:
	max := config.AppConfig.RateLimitMax
	expiration := time.Duration(config.AppConfig.RateLimitWindow) * time.Second
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		LimitReached: func(c *fiber.Ctx) error {
			return helper.SendResponse(c, fiber.StatusTooManyRequests, "Rate limit exceeded", nil)
		},
	})
}
