package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/helper"
	"github.com/KeyzarRasya/comdel-server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
)


func MainHandlers(c *fiber.Ctx) error {
	log.Info(string(c.Request().Header.Cookie("jwt")));
	return c.SendString("Hello");
}

func OAuthLoginHandler(c *fiber.Ctx) error {
	googleOAuth := config.OAuthConfig();
	var state string = helper.GenerateRandState();
	authUrl := googleOAuth.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce);

	var cookie fiber.Cookie = fiber.Cookie{
		Name: "state",
		Value: state,
		Expires: time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
		Path: "/",
		SameSite: "Lax",
	}

	c.Cookie(&cookie);

	return c.Redirect(authUrl);
}

func RedirectHandler(c *fiber.Ctx) error {
	state := c.Query("state");
	code := c.Query("code");
	expectedState := c.Cookies("state");
	googleOAuth := config.OAuthConfig();

	if expectedState == "" || expectedState != state {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid OAuth state")
	}	

	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Code");
	}

	token, err := googleOAuth.Exchange(context.Background(), code);

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to exchange the code");
	}

	client := googleOAuth.Client(context.Background(), token);

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo");

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("failed to get response of user info");
	}

	defer func () {
		resp.Body.Close();
	}()

	c.Cookie(&fiber.Cookie{
		Name:     "state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Waktu lampau
		HTTPOnly: true,
	})
		
	var userInfo dto.GoogleProfile;

	err = json.NewDecoder(resp.Body).Decode(&userInfo);

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to decode JSON");
	}
	
	var res dto.Response = services.SaveUser(userInfo, token);
	
	log.Info(res)

	if res.Status == fiber.StatusOK {
		profile, ok := res.Data.(dto.GoogleProfile);


		if !ok {
			panic("Data is not profileToken type");
		}

		cookies := fiber.Cookie{
			Name: "jwt",
			Value: profile.Token,
			HTTPOnly: true,
			Expires: time.Now().Add(30 * time.Hour),
			SameSite: "Lax",
		}

		log.Info(cookies)

		c.Cookie(&cookies);
	}


	return c.Redirect("http://localhost:5173/dashboard");
	
}

func GetUserInfo(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));

	var response dto.Response = services.GetUser(jwtCookies);
	return c.JSON(response.JSON());
}

func AddVideo(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));
	link := c.Query("vid");
	scheduler := c.Query("sc");
	strategy := c.Query("st");

	response := services.UploadVideo(jwtCookies, dto.UploadVideos{
		Link: link,
		Scheduler: scheduler,
		Strategy: strategy,
	})

	return c.JSON(response.JSON());
}

func RequestCommentTheadsHandler(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));
	var link string = c.Query("vid");
	var response dto.Response = services.GetComments(jwtCookies, link)

	return c.JSON(response.JSON());
}