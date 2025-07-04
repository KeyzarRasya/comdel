package handlers

import (
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CheckOwnership(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));
	link := c.Query("vid");

	var response dto.Response = services.CheckVideoOwnership(link, jwtCookies);

	return c.Status(response.Status).JSON(response.JSON());
}

func FetchVideoInfo(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));

	videoId := c.Query("id");
	
	var response dto.Response = services.GetVideoInformation(videoId, jwtCookies);
	log.Info("Response : ", response);

	return c.Status(response.Status).JSON(response.JSON());
}