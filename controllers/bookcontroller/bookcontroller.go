package bookcontroller

import (
	"errors"

	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func Index(c *fiber.Ctx) error {
	var books []models.Book;
	result := models.DB.Find(&books);

	if(result.Error != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(books);
}
func Show(c *fiber.Ctx) error {
	id:= c.Params("id");
	var book models.Book;
	
	result := models.DB.First(&book, id);

	if(result.Error != nil){
		if(errors.Is(result.Error,gorm.ErrRecordNotFound)){
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error.Error(),
		})
	}

	return c.JSON(book);
}
func Create(c *fiber.Ctx) error {
	var book models.Book;

	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	
	if err := models.DB.Create(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}
func Update(c *fiber.Ctx) error {
	id := c.Params("id");

	existingBook := models.DB.First(&models.Book{}, id);

	if errors.Is(existingBook.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}


	var book models.Book;
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if models.DB.Where("id = ?", id).Updates(&book).RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to update data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data updated successfully",
	});
}
func Delete(c *fiber.Ctx) error {
	id := c.Params("id");

	existingBook := models.DB.First(&models.Book{}, id);

	if errors.Is(existingBook.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	if models.DB.Delete(&models.Book{}, id).RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to delete data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data deleted successfully",
	});
}