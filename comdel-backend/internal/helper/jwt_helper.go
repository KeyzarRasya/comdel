package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId string) (string, error) {
	claim := jwt.MapClaims{
		"userId": userId,
		"exp": time.Now().Add((24 * time.Hour) * 30).Unix(),
	};
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim);

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")));
	
	if err != nil {
		log.Info(err);
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	// Extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Could not extract claims")
	}
}

func VerifyAndGet(cookie string) (string, error) {
	verifyCookie, err := VerifyToken(cookie);

	if err != nil {
		return "", fmt.Errorf("Failed to verify token");
	}

	userId, ok := verifyCookie["userId"].(string)

	if !ok || userId == "" {
		return "", fmt.Errorf("Invalid Token");
	}

	return userId, nil;
}
