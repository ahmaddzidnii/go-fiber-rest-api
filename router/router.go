package router

import (
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/auth"
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/bookcontroller"
	usercontroller "github.com/ahmaddzidnii/go-fiber-rest-api/controllers/user"
	"github.com/ahmaddzidnii/go-fiber-rest-api/middleware"
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
	authRoute.Post("/register", auth.Register);
	authRoute.Post("/login", auth.Login);
	authRoute.Post("/renew", auth.Renew);

	// User
	user := api.Group("/users")

	// Protected Route
	user.Use(middleware.AuthMiddleware)
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