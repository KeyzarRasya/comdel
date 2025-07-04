package routes

import (
	"github.com/KeyzarRasya/comdel-server/internal/handlers"
	"github.com/KeyzarRasya/comdel-server/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	
	app.Route("/user", func (route fiber.Router) {
		route.Get("/", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, handlers.MainHandlers);
		route.Get("/info", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, middleware.LazyUnsubscribe, handlers.GetUserInfo)
	})

	app.Route("/videos", func(route fiber.Router) {
		route.Get("/ownership", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, handlers.CheckOwnership);
		route.Get("/information", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, handlers.FetchVideoInfo)
		route.Post("/upload", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, handlers.AddVideo)
		route.Get("/comment", middleware.AuthenticationMiddleware, handlers.RequestCommentTheadsHandler)
	})

	app.Route("/auth", func (route fiber.Router) {
		route.Get("/google", handlers.OAuthLoginHandler);
		route.Get("/google/redirect", handlers.RedirectHandler)
	})

	app.Route("/payment", func(route fiber.Router) {
		route.Get("/pay", handlers.CreatePayment)
		route.Get("/finish", handlers.HandleFinishPayment)
	})

}