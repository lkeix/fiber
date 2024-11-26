package proxy

import (
	"github.com/gofiber/fiber/v3"
)

type RetryIf func(ctx fiber.Ctx, cb CircuitBreaker, err error) bool