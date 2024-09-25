package main

import (
	"log"
	"time"

	"github.com/goccy/go-json"

	"github.com/ahmaddzidnii/go-fiber-rest-api/config"
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/ahmaddzidnii/go-fiber-rest-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup(app *fiber.App) {
	// Init Database
	config.ConnectDatabase();

	// Setup Cors
	app.Use(cors.New())

	// Logging
	app.Use(logger.New())

	// Rate Limiter
	app.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return helpers.Response(c, fiber.StatusTooManyRequests, "Too many requests", nil)
		},
		SkipFailedRequests: false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware: limiter.FixedWindow{},
	}))

	// Setup Router
	router.SetupRouter(app);
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	});

	// Setup
	Setup(app);

	log.Fatal(app.Listen("127.0.0.1:2000"))
}