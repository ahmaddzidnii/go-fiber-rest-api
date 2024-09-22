package main

import (
	"log"

	"github.com/ahmaddzidnii/go-fiber-rest-api/controllers/bookcontroller"
	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/gofiber/fiber/v2"
)

func main() {
	models.ConnectDatabase();

	app := fiber.New();

	api := app.Group("/api");	
	book := api.Group("/books");

	book.Get("/", bookcontroller.Index);
	book.Get("/:id", bookcontroller.Show);
	book.Post("", bookcontroller.Create);
	book.Put("/:id", bookcontroller.Update);
	book.Delete("/:id", bookcontroller.Delete);

	log.Fatal(app.Listen(":3000"))
}