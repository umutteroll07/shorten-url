package main

import (
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload"
	"github.com/umutteroll07/ShortenURL/app/route"
	_ "github.com/umutteroll07/ShortenURL/internal/database"
	_ "github.com/umutteroll07/ShortenURL/internal/migration"
)

func main() {
	app := fiber.New()
	route.SetupRoute(app)
	app.Listen(":3000")
}
