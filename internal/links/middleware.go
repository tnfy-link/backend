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

func NewIDValidator(validator func(string) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := validator(id); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.Next()
	}
}
