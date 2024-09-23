package middleware

import (
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Mengambil token dari cookie
	token := c.Cookies("access_token");

	// Jika token kosong atau tidak valid
	if token == "" {
		return helpers.Response(c, fiber.StatusUnauthorized, "Jwt must be provided", nil);
	}

	// Mengecek token
	user, err := helpers.ClaimJWT(token);

	if err != nil {
		return helpers.Response(c, fiber.StatusUnauthorized, err.Error(), nil);
	}

	// Membuat konteks user baru
	c.Locals("currentUser", user)

	// Continue to the next middleware
	return c.Next()
}