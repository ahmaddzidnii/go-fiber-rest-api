package bookcontroller

import (
	"errors"

	"github.com/ahmaddzidnii/go-fiber-rest-api/config"
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func Index(c *fiber.Ctx) error {
	var books []models.Book;
	result := config.DB.Find(&books);

	if(result.Error != nil){
		return helpers.Response(c, fiber.StatusInternalServerError, result.Error.Error(), nil);
	}

	return helpers.Response(c, fiber.StatusOK,"Success retrive data", books);
}

func Show(c *fiber.Ctx) error {
	id:= c.Params("id");
	var book models.Book;
	
	result := config.DB.First(&book, id);

	if(result.Error != nil){
		if(errors.Is(result.Error,gorm.ErrRecordNotFound)){
			return helpers.Response(c, fiber.StatusNotFound, result.Error.Error(), nil);
		}

		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	return helpers.Response(c, fiber.StatusOK, "Succses retrive data", book);
}

func Create(c *fiber.Ctx) error {
	var book models.Book;

	if err := c.BodyParser(&book); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error(), nil);
	}
	
	if err := config.DB.Create(&book).Error; err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error(), nil);
	}

	return helpers.Response(c, fiber.StatusCreated, "Succses create data", book);
}

func Update(c *fiber.Ctx) error {
	id := c.Params("id");

	existingBook := config.DB.First(&models.Book{}, id);

	if errors.Is(existingBook.Error, gorm.ErrRecordNotFound) {
		return helpers.Response(c, fiber.StatusNotFound, "Book not found", nil);
	}


	var book models.Book;
	if err := c.BodyParser(&book); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error(), nil);
	}

	if config.DB.Where("id = ?", id).Updates(&book).RowsAffected == 0 {
		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	return helpers.Response(c, fiber.StatusOK, "Data updated successfully", book);
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id");

	existingBook := config.DB.First(&models.Book{}, id);

	if errors.Is(existingBook.Error, gorm.ErrRecordNotFound) {
		return helpers.Response(c, fiber.StatusNotFound, "Book not found", nil);
	}

	if config.DB.Delete(&models.Book{}, id).RowsAffected == 0 {
		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	return helpers.Response(c, fiber.StatusOK, "Succsess delete data", nil);
}