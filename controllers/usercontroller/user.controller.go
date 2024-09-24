package usercontroller

import (
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetMe(c *fiber.Ctx) error {
	user := c.Locals("currentUser").(*helpers.MyJwtClaims)

	return helpers.Response(c, fiber.StatusOK, "Success retrive data", user)
}