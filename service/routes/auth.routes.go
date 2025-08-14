package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
)

type authRoutes struct {
	Router fiber.Router
}

func InitAuthRoutes(router fiber.Router) *authRoutes {
	return &authRoutes{
		Router: router,
	}
}

func (r *authRoutes) Routes() {

	router := r.Router

	authGroup := router.Group("/auth")
	authGroup.Post("/login", func(c *fiber.Ctx) error {
		token, err := middleware.GenerateTokenHS256("user")
		if err != nil {
			return helper.SendResponse(c, fiber.StatusInternalServerError, "Failed to generate token", nil)
		}
		return helper.SendResponse(c, fiber.StatusOK, "Login successful", token)
	})
	authGroup.Post("/register", func(c *fiber.Ctx) error {
		return helper.SendResponse(c, fiber.StatusCreated, "Register successful", nil)
	})

	profile := router.Group("/user")
	profile.Get("/profile", middleware.JWTHS256Middleware(func(c *fiber.Ctx) error {
		return helper.SendResponse(c, fiber.StatusOK, "Profile retrieved successfully", c.Locals("user_claims"))
	}, "admin", "superadmin"))
}
