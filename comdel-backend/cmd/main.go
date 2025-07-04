package main

import (
	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env");
	// c := cron.New()

	// err = c.AddFunc("@every 10m", func() {
	// 	if err := services.CronFetchDelete(); err != nil {
	// 		log.Info("Error running CronFetchDelete:")
	// 		log.Info(err.Error())
	// 	} else {
	// 		log.Info("Success doing cron job")
	// 	}
	// })

	// if err != nil {

	// 	log.Fatal("Failed to do cron jobs")

	// }

	// c.Start()

	if err != nil {
		log.Fatal("Failed to load .env files");
	}

	serverConfig := config.ServerConfig();
	app := fiber.New(serverConfig);

	corsConfig := cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}
	


	app.Use(cors.New(corsConfig));
	routes.UserRoute(app);
	app.Listen(":8080");
}