package middleware

import (
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func LazyUnsubscribe(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"))

	var response dto.Response = services.Unsubscribe(jwtCookies);

	log.Info(response.Message)
	log.Info(response.Data)

	return c.Next()
}