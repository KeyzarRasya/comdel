package middleware

import (
	"comdel-backend/internal/dto"
	"comdel-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

	var response dto.Response = sm.PaymentService.Unsubscribe(jwtCookies)

	log.Info(response.Message)
	log.Info(response.Data)

	return c.Next()
}