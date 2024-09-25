package middlewares

import (
	"github.com/ahmaddzidnii/go-fiber-rest-api/config"
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/gofiber/fiber/v2"
)

func AdminOnly(c *fiber.Ctx) error {
	accsess_token := c.Cookies("access_token");

	if(accsess_token == ""){
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil);
	}

	claims, err := helpers.ClaimJWT(accsess_token);

	if(err != nil){
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil);
	}

	var user models.User;

	if err := config.DB.Where("id = ?", claims.Id).First(&user).Error; err != nil {
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil);
	}

	if(user.Role != "admin"){
		return helpers.Response(c, fiber.StatusForbidden, "Forbidden", nil);
	}
	return c.Next()
}