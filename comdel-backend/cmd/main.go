package main

import (
	"comdel-backend/internal/config"
	"comdel-backend/internal/handlers"
	"comdel-backend/internal/middleware"
	"comdel-backend/internal/repository"
	"comdel-backend/internal/routes"
	"comdel-backend/internal/services"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

	dbLoader := config.DBLoaderImpl{}
	conn, err := dbLoader.Load()

	if err != nil {
		log.Info("Failed to load database at start")
		return;
	}

	/* Google OAuth Config */
	var scopes []string = []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/youtube.force-ssl",
	}
	oauthConfig := oauth2.Config{
		RedirectURL: 	"http://localhost:8080/auth/google/redirect",
		ClientID: 		os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: 	os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: 		scopes,
		Endpoint: 		google.Endpoint,
		
	}

	/* Google Oauth Provider */
	googleOauth := config.NewGoogleOAuth(&oauthConfig)

	/*
		Repository Dependency
	*/
	transactionRepository := repository.NewTransactionRepository(conn)
	userRepository := repository.NewUserRepository(conn)
	videoRepository := repository.NewVideoRepository(conn)
	tokenRepository := repository.NewTokenRepository(conn)
	subscriptionRepository := repository.NewSubscriptionRepository(conn)
	commentRepository := repository.NewCommentRepository(conn)

	/* Service Dependency */
	youtubeService := services.NewYoutubeService(googleOauth)

	commentService := services.NewCommentService(
		userRepository,
		tokenRepository,
		&youtubeService,
		commentRepository,
		videoRepository,
		googleOauth,
		&dbLoader,
	)

	/* Dependency Config*/
	auth := services.Authentication{}
	
	/*
	===START===
	Service Dependency
	*/
	// 1. User Service Dependency Injection
	ytService := services.YoutubeServiceImpl{OAuthProvider: googleOauth}
	userService := services.NewUserService(
		userRepository,
		tokenRepository,
		videoRepository,
		&auth,
		&dbLoader,
		googleOauth, 
		&ytService,

	)

	//2. Video Service Dependency Injection
	videoService := services.NewVideoService(
		userRepository,
		videoRepository,
		tokenRepository,
		&commentService,
		commentRepository,
		&ytService,
		&dbLoader,
		&auth,
	)

	// 3. Payment Service Dependency Injection
	paymentService := services.NewPaymentService(
		userRepository,
		transactionRepository,
		subscriptionRepository,
		&dbLoader,
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