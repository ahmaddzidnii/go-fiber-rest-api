package router

import (
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/authcontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/bookcontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/usercontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message" : "Api is running",
		})
	})

	// Auth
	authRoute := api.Group("/auth");
	authRoute.Post("/register", authcontroller.Register);
	authRoute.Post("/login", authcontroller.Login);
	authRoute.Get("/renew", authcontroller.Renew);
	authRoute.Get("/logout", authcontroller.Logout);

	// User
	user := api.Group("/users")

	// Protected Route
	user.Use(middlewares.AuthMiddleware)
	user.Get("/me", usercontroller.GetMe)
	// user.Post("/", helpers.CreateUser)
	// user.Patch("/:id", middleware.Protected(), helpers.UpdateUser)
	// user.Delete("/:id", middleware.Protected(), helpers.DeleteUser)

	// Book
	
	book := api.Group("/books");
	book.Get("/", bookcontroller.Index);
	book.Get("/:id", bookcontroller.Show);
	book.Post("/", bookcontroller.Create);
	book.Put("/:id", bookcontroller.Update);
	book.Delete("/:id", bookcontroller.Delete);
}