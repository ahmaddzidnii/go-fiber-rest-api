package router

import (
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/authcontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/bookcontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/usercontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message" : "Api is running",
		})
	})

	api := app.Group("/api");
	// Auth
	authRoute := api.Group("/auth");
	authRoute.Post("/register", authcontroller.Register);
	authRoute.Post("/login", authcontroller.Login);
	authRoute.Get("/renew", authcontroller.Renew);
	authRoute.Get("/logout", authcontroller.Logout);
	authRoute.Get("/session", authcontroller.Session);

	// User
	user := api.Group("/users")
	user.Use(middlewares.AuthMiddleware)
	user.Get("/me", usercontroller.GetMe)

	// Book
	
	book := api.Group("/books");
	book.Get("/",middlewares.AdminOnly,bookcontroller.Index);
	book.Get("/:id", bookcontroller.Show);
	book.Post("/", middlewares.AdminOnly, bookcontroller.Create);
	book.Put("/:id", bookcontroller.Update);
	book.Delete("/:id", bookcontroller.Delete);
}