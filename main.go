package main

import (
	"log"

	"github.com/goccy/go-json"

	"github.com/ahmaddzidnii/go-fiber-rest-api/config"
	"github.com/ahmaddzidnii/go-fiber-rest-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {
	// Setup Cors
	app.Use(cors.New())

	// Init Database
	config.ConnectDatabase();

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