package middleware

import (
	"context"
	"time"

	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
)

func AuthenticationMiddleware(c *fiber.Ctx) error {
	var jwtCookies string = string(c.Request().Header.Cookie("jwt"))

	if jwtCookies == "" {
		return c.Redirect("/auth/google");
	}

	return c.Next();
}

func RefreshTokenMiddleware(c *fiber.Ctx) error {
    conn := config.LoadDatabase()
    oauthConfig := config.OAuthConfig()
    jwtCookie := string(c.Request().Header.Cookie("jwt"))
    userId, err := helper.VerifyAndGet(jwtCookie)
    log.Info("Refresh: ", userId);
    if err != nil {
		var response dto.Response = dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to verify token",
			Data: nil,
		}
        return c.Status(fiber.StatusBadRequest).JSON(response.JSON())
    }

    var refreshToken, accessToken, tokenId string
    var expiry time.Time

    err = conn.QueryRow(
        context.Background(),
        "SELECT access_token, refresh_token, expiry, token_id FROM oauth_token WHERE owner_id=$1",
        userId,
    ).Scan(&accessToken, &refreshToken, &expiry, &tokenId)
    if err != nil {
		var response dto.Response = dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to fetch token",
			Data: nil,
		}
        return c.Status(fiber.StatusInternalServerError).JSON(response.JSON())
    }

    if time.Now().After(expiry) {
        token := &oauth2.Token{
			AccessToken: accessToken,
            RefreshToken: refreshToken,
			Expiry: expiry,
        }
        ts := oauthConfig.TokenSource(context.Background(), token)
        newToken, err := ts.Token()
        if err != nil {
			log.Info(err);
			var response dto.Response = dto.Response{
                Status: fiber.StatusBadRequest,
                Message: "Failed to refresh token",
                Data: nil,
			}
            return c.Status(fiber.StatusBadRequest).JSON(response.JSON())
        }

        _, err = conn.Exec(
            context.Background(),
            "UPDATE oauth_token SET access_token=$1, expiry=$2 WHERE token_id=$3",
            newToken.AccessToken, newToken.Expiry, tokenId,
        )
        if err != nil {
			var response dto.Response = dto.Response{
                Status: fiber.StatusBadRequest,
                Message: "Failed to update token",
                Data: nil,
			}
            return c.Status(fiber.StatusInternalServerError).JSON(response.JSON())
        }
    }

    return c.Next()
}
