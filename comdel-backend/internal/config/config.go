package config

import (
	"context"
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func ServerConfig() fiber.Config {
	config := fiber.Config{
		AppName: "ComdelServer",

	}
	return config;
}

func OAuthConfig() oauth2.Config {
	var scopes []string = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/youtube.force-ssl",
	}
	googleOAuth := oauth2.Config{
		RedirectURL: 	"http://localhost:8080/auth/google/redirect",
		ClientID: 		os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: 	os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: 		scopes,
		Endpoint: 		google.Endpoint,
		
	}

	return googleOAuth;
}

func LoadDatabase() *pgx.Conn {

	var databaseUri string;

	if os.Getenv("DEV_ENV") == "dev" {
		databaseUri = os.Getenv("LOCAL_DATABASE_URI")
	} else {
		databaseUri = os.Getenv("DATABASE_URI")
	}

	log.Info(databaseUri)

	conn, err := pgx.Connect(context.Background(), databaseUri);

	if err != nil {
		log.Info(os.Getenv("DATABASE_URI"))
		
		log.Error("Failed to Load Database");
		log.Error(err.Error())
		return nil;
	}

	log.Info("Creating Database Connection");
	return conn;
}

func PaymentConfig() (*snap.Client, error) {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY");

	if serverKey == "" {
		return nil, errors.New("Missing Server Key");
	}
	var s snap.Client;

	s.New(serverKey, midtrans.Sandbox);

	return &s, nil;
}
