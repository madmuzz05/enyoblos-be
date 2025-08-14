package helper

import (
	"github.com/gofiber/fiber/v2"
)

func SendResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"data":    data,
		"message": message,
		"success": status < 400,
	})
}
