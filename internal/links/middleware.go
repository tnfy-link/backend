package links

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func NewLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                1,
		SkipFailedRequests: true,
		Expiration:         time.Second,
	})
}
