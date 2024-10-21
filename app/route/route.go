package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/umutteroll07/ShortenURL/app/handler"
	"github.com/umutteroll07/ShortenURL/app/middleware"
)

func SetupRoute(app *fiber.App) {
	app.Post("/shortUrl", middleware.RateLimiter, handler.ShortenUrl)
	app.Get("/shortenUrl/:shortUrl", handler.RedirectUrl)
	app.Get("/short_url", handler.GetAllShortenUrl)
	app.Delete("shortenUrl", handler.DeleteShortenUrl)
}
