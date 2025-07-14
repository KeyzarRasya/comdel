package main

import (
	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/handlers"
	"github.com/KeyzarRasya/comdel-server/internal/middleware"
	"github.com/KeyzarRasya/comdel-server/internal/repository"
	"github.com/KeyzarRasya/comdel-server/internal/routes"
	"github.com/KeyzarRasya/comdel-server/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

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

func main() {
	err := godotenv.Load("../.env");
	conn := config.LoadDatabase()

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

	/*
		Repository Dependency
	*/
	userRepository := repository.NewUserRepository(conn)
	videoRepository := repository.NewVideoRepository(conn)
	transactionRepository := repository.NewTransactionRepository(conn)
	tokenRepository := repository.NewTokenRepository(conn)
	subscriptionRepository := repository.NewSubscriptionRepository(conn)
	commentRepository := repository.NewCommentRepository(conn)


	/*
		===START===
		Service Dependency
	*/
	// 1. User Service Dependency Injection
	userService := services.NewUserService(
		userRepository,
		tokenRepository,
		videoRepository,
	)

	//2. Video Service Dependency Injection
	videoService := services.NewVideoService(
		userRepository,
		videoRepository,
		tokenRepository,
		commentRepository,
	)

	// 3. Payment Service Dependency Injection
	paymentService := services.NewPaymentService(
		userRepository,
		transactionRepository,
		subscriptionRepository,
	)
	/*
		Service Dependency
		===END===
	*/



	/* Handler Dependency */
	userHandlers := handlers.NewUserHandlers(userService)
	videoHandlers := handlers.NewVideoHandlers(videoService)
	paymentHandlers := handlers.NewPaymentHandlers(paymentService);

	/* Middleware Injecting*/
	subsciptionMiddleware := middleware.NewSubscriptionMiddleware(paymentService)

	/*
		Route
	*/
	route := routes.NewRoute(
		userHandlers,
		videoHandlers,
		paymentHandlers,
		subsciptionMiddleware,
	)
	
	


	app.Use(cors.New(corsConfig));
	route.UserRoute(app);
	app.Listen(":8080");
}