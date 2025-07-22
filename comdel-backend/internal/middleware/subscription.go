package middleware

import (
	"comdel-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionMiddleware interface {
	LazyUnsubscribe(c *fiber.Ctx) error;
}

type SubscriptionMiddlewareImpl struct {
	PaymentService services.PaymentService
}

func NewSubscriptionMiddleware(
	paymentService services.PaymentService,
) SubscriptionMiddleware {
	return &SubscriptionMiddlewareImpl{
		PaymentService:  paymentService,
	}
}

func (sm *SubscriptionMiddlewareImpl) LazyUnsubscribe(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"))

	sm.PaymentService.Unsubscribe(jwtCookies)
	return c.Next()
}