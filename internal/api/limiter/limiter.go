package limiter

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func New(limit int) fiber.Handler {
	if limit <= 0 {
		panic("limit must be greater than 0")
	}

	return limiter.New(limiter.Config{
		Max:                limit,
		SkipFailedRequests: true,
		Expiration:         time.Second,
	})
}
