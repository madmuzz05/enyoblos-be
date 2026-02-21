package helper

import (
	"github.com/gofiber/fiber/v3"
)

func SendResponse(c fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"data":    data,
		"message": message,
		"success": status < 400,
	})
}

func SendErrorResponse(c fiber.Ctx, status int, message string, err error) error {
	var errMsg interface{}
	if err != nil {
		errMsg = fiber.Map{"error": err.Error()}
	}
	return c.Status(status).JSON(fiber.Map{
		"data":    errMsg,
		"message": message,
		"success": false,
	})
}
