package handlers

import (
	"context"
	"encoding/json"
	"time"

	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/helper"
	"comdel-backend/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
)

type UserHandlers interface {
	Main(c *fiber.Ctx)			error
	OAuthLogin(c *fiber.Ctx)	error
	OAuthRedirect(c *fiber.Ctx)	error
	UserInfo(c *fiber.Ctx)		error
}

type UserHandlersImpl struct {
	UserService services.UserService
}

func NewUserHandlers(UserService services.UserService) UserHandlers {
	return &UserHandlersImpl{UserService: UserService}
}


func (uh *UserHandlersImpl) Main(c *fiber.Ctx) error {
	log.Info(string(c.Request().Header.Cookie("jwt")));
	return c.SendString("Hello");
}

func (uh *UserHandlersImpl) OAuthLogin(c *fiber.Ctx) error {
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

func (uh *UserHandlersImpl) OAuthRedirect(c *fiber.Ctx) error {
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

	defer resp.Body.Close()

	c.Cookie(&fiber.Cookie{
		Name:     "state",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Waktu lampau
		HTTPOnly: true,
	})
		
	var userInfo dto.GoogleProfile;

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to decode JSON");
	}

	
	var res dto.Response = uh.UserService.SaveUser(userInfo, token);
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

func (uh *UserHandlersImpl) UserInfo(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"));

	var response dto.Response = uh.UserService.GetUser(jwtCookies)
	return c.JSON(response.JSON());
}