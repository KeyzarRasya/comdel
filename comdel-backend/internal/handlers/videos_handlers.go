package handlers

import (
	"comdel-backend/internal/dto"
	"comdel-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type VideoHandlers interface {
	AddVideo(c *fiber.Ctx) 			error;
	CheckOwnership(c *fiber.Ctx)	error;
	VideoInfo(c *fiber.Ctx)			error;
}

type VideoHandlersImpl struct {
	VideoService services.VideoService
}

func NewVideoHandlers(videoService services.VideoService) VideoHandlers {
	return &VideoHandlersImpl{VideoService: videoService}
}

func (vh *VideoHandlersImpl) AddVideo(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));
	link := c.Query("vid");
	scheduler := c.Query("sc");
	strategy := c.Query("st");

	uploadVideos := dto.UploadVideos{
		Link: link,
		Scheduler: scheduler,
		Strategy: strategy,
	}

	response := vh.VideoService.UploadVideo(jwtCookies, uploadVideos)

	return c.JSON(response.JSON());
}

func (vh *VideoHandlersImpl) CheckOwnership(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));
	link := c.Query("vid");

	var response dto.Response = vh.VideoService.IsCanUpload(link, jwtCookies);
	return c.Status(response.Status).JSON(response.JSON());
}

func (vh *VideoHandlersImpl) VideoInfo(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));

	videoId := c.Query("id");
	
	var response dto.Response = vh.VideoService.Info(videoId, jwtCookies)
	log.Info("Response : ", response);

	return c.Status(response.Status).JSON(response.JSON());
}