package authcontroller

import (
	"errors"
	"time"

	"github.com/ahmaddzidnii/go-fiber-rest-api/config"
	"github.com/ahmaddzidnii/go-fiber-rest-api/helpers"
	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"github.com/ahmaddzidnii/go-fiber-rest-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
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
	// Mendapatkan refresh token dari cookie
	refresh_token := c.Cookies("refresh_token")

	// mencari user berdasarkan refresh token
	var user models.User
	if err := config.DB.Where("refresh_token = ?", refresh_token).First(&user).Error; err != nil {
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	_,err := helpers.ClaimJWT(refresh_token);

	if(err != nil){
		return helpers.Response(c, fiber.StatusUnauthorized, err.Error(), nil);
	}

	// membuat access token baru
	access_token, err := helpers.GenerateJWT(&user, jwt.NewNumericDate(time.Now().Add(time.Minute * 2)))
	if err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, err.Error(), nil)
	}

	// membuat cookie
	cookie := new(fiber.Cookie)
	cookie.Name = "access_token"
	cookie.Value = access_token
	cookie.HTTPOnly = true

	// mengirimkan cookie
	c.Cookie(cookie)

	// mengembalikan response access token yang baru
	return helpers.Response(c, fiber.StatusOK, "Success renew access token", fiber.Map{
		"access_token": access_token,
	})
}

func Logout(c *fiber.Ctx) error {
	// Mendapatkan access token dari cookie
	access_token := c.Cookies("access_token");

	if(access_token == ""){
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil);
	}

	cc, err := helpers. ClaimJWT(access_token);

	if(err != nil){
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", nil);
	}

	// Mencari user berdasarkan access token
	var user models.User;
	if err := config.DB.Where("id = ?", cc.Id).First(&user).Error; err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return helpers.Response(c, fiber.StatusNotFound, "User not found", nil);
		}

		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	// Menghapus refresh token dari database
	if err := config.DB.Model(&user).Update("refresh_token", nil).Error; err != nil {
		return helpers.Response(c, fiber.StatusInternalServerError, "Internal server error", nil);
	}

	// Menghapus cookie
	utils.ClearCookies((c), "access_token", "refresh_token");

	// Mengembalikan response logout
	return helpers.Response(c, fiber.StatusOK, "Success logout", nil)
}

func Session(c *fiber.Ctx) error {
	accsess_token := c.Cookies("access_token");

	if(accsess_token == ""){
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", fiber.Map{
			"session": nil,
		});
	}
	s, err := helpers. ClaimJWT(accsess_token);

	if(err != nil){
		return helpers.Response(c, fiber.StatusUnauthorized, "Unauthorized", fiber.Map{
			"session": nil,
		});
	}

	return helpers.Response(c, fiber.StatusOK, "Succsess retrieve data", fiber.Map{
		"session": s,
	});
}