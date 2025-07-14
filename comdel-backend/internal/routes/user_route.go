package routes

import (
	"comdel-backend/internal/handlers"
	"comdel-backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type Router interface {
	UserRoute(app *fiber.App) error;
}

type Route struct {
	UserHandlers handlers.UserHandlers
	VideoHandlers handlers.VideoHandlers
	SubscriptionMiddleware middleware.SubscriptionMiddleware
	PaymentHandlers handlers.PaymentHandlers
}

func NewRoute(
	userHandlers handlers.UserHandlers,
	videoHandlers handlers.VideoHandlers,
	paymentHanders handlers.PaymentHandlers,
	subscriptionMiddleware middleware.SubscriptionMiddleware,
) Route {
	return Route{
		UserHandlers: userHandlers, 
		VideoHandlers: videoHandlers,
		PaymentHandlers: paymentHanders,
		SubscriptionMiddleware: subscriptionMiddleware,
	}
}

func (r *Route) UserRoute(app *fiber.App) {
	
	app.Route("/user", func (route fiber.Router) {
		route.Get("/", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, r.UserHandlers.Main);
		route.Get("/info", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, r.SubscriptionMiddleware.LazyUnsubscribe, r.UserHandlers.UserInfo)
	})

	app.Route("/videos", func(route fiber.Router) {
		route.Get("/ownership", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, r.VideoHandlers.CheckOwnership);
		route.Get("/information", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, r.VideoHandlers.VideoInfo)
		route.Post("/upload", middleware.AuthenticationMiddleware, middleware.RefreshTokenMiddleware, r.VideoHandlers.AddVideo)
	})

	app.Route("/auth", func (route fiber.Router) {
		route.Get("/google", r.UserHandlers.OAuthLogin);
		route.Get("/google/redirect", r.UserHandlers.OAuthRedirect)
	})

	app.Route("/payment", func(route fiber.Router) {
		route.Get("/pay", r.PaymentHandlers.CreatePayment)
		route.Get("/finish", r.PaymentHandlers.FinishPayment)
	})

}