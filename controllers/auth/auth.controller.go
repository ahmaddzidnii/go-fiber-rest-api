package auth

import (
	"time"

	"github.com/ahmaddzidnii/go-fiber-rest-api/config"
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Register(c *fiber.Ctx) error {
	// Mendefinisikan struct Register
	var register models.Register;

	// Memasukkan data dari request ke struct Register
	if err := c.BodyParser(&register); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error(), nil);	
	}

	// Cek password dan konfirmasi password
	if(register.Password != register.ConfirmPassword){
		return helpers.Response(c, fiber.StatusBadRequest, "Password and confirm password not match", nil);
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(register.Password);

	if(err != nil){
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error(), nil);
	}
	
	// Membuat struct User
	user := models.User{
		FullName: register.FullName,
		Username: register.Username,
		Email: register.Email,
		Password: hashedPassword,
	}

	// Simpan data user ke database
	if err := config.DB.Create(&user).Error; err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	response := models.ResponseRegister{
		Id: user.Id,
		FullName: user.FullName,
		Username: user.Username,
		Email: user.Email,
	}
	// Mengembalikan data register
	return helpers.Response(c, fiber.StatusCreated, "Success register user", response);
}

func Login(c *fiber.Ctx) error {
	// Mendefinisikan struct Login
	var login models.Login;

	// Memasukkan data dari request ke struct Login
	if err := c.BodyParser(&login); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, err.Error(), nil);
	}

	// Mencari user berdasarkan email apakah ditemukan?
	var user models.User;
	if err := config.DB.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return helpers.Response(c, fiber.StatusNotFound, "Wrong email or password", nil);
	}

	// Cek apakah password benar
	if err := helpers.ComparePassword(user.Password, login.Password); err != nil {
		return helpers.Response(c, fiber.StatusBadRequest, "Wrong email or password", nil);
	}

	// generate refresh token
	refresh_token,errGenerateRefreshToken := helpers.GenerateJWT(&user, jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)));

	if errGenerateRefreshToken != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, errGenerateRefreshToken.Error(), nil);
	}

	// Simpan refresh token ke database
	if err := config.DB.Model(&user).Update("refresh_token", refresh_token).Error; err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refresh_token
	cookie.HTTPOnly = true
	
	// Set cookie
	c.Cookie(cookie)

	// Generate accses token
	access_token, err := helpers.GenerateJWT(&user,jwt.NewNumericDate(time.Now().Add(time.Minute * 2)));

	if(err != nil){
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error(), nil);
	}
	// Create cookie
	cookie.Name = "access_token"
	cookie.Value = access_token
	cookie.HTTPOnly = true
	
	// Set cookie
	c.Cookie(cookie)


	// Mengembalikan data login
	return helpers.Response(c, fiber.StatusOK, "Success login", fiber.Map{
		"access_token": access_token,
		"refresh_token": refresh_token,
	})
}

func Renew(c *fiber.Ctx) error {
	// Get refresh token from cookie
	refresh_token := c.Cookies("refresh_token")

	// Find user by refresh token
	var user models.User
	if err := config.DB.Where("refresh_token = ?", refresh_token).First(&user).Error; err != nil {
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	// Generate new access token
	access_token, err := helpers.GenerateJWT(&user, jwt.NewNumericDate(time.Now().Add(time.Minute * 2)))
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// Create cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "access_token"
	cookie.Value = access_token
	cookie.HTTPOnly = true

	// Set cookie
	c.Cookie(cookie)

	// Return new access token
	return helpers.Response(c, fiber.StatusOK, "Success renew refresh token", fiber.Map{
		"access_token": access_token,
	})
}