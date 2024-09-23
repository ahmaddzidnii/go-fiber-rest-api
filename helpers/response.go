package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type ResponseWithData struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

type ResponseWithoutData struct {
	Status int `json:"status"`
	Message string `json:"message"`
}

func Response(c *fiber.Ctx, statusCode int, message string, payload interface{}) error {
	c.Set("Content-Type", "application/json");
	c.Set("X-Powered-By", "Go Fiber");

	var response interface{};

	if payload != nil {
		response = &ResponseWithData{
			Status: statusCode,
			Message: message,
			Data: payload,
		}
	} else {
		response = &ResponseWithoutData{
			Status: statusCode,
			Message: message,
		}
	}

	return c.Status(statusCode).JSON(response);
}