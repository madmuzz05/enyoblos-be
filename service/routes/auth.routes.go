package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	authHandler "github.com/madmuzz05/be-enyoblos/service/module/auth/handler"
)

type authRoutes struct {
	Router      fiber.Router
	AuthHandler *authHandler.AuthHandler
	RedisClient *redisdb.RedisClient
}

func InitAuthRoutes(router fiber.Router, handler *authHandler.AuthHandler, redis *redisdb.RedisClient) *authRoutes {
	return &authRoutes{
		Router:      router,
		AuthHandler: handler,
		RedisClient: redis,
	}
}

func (r *authRoutes) Routes() {
	authGroup := r.Router.Group("/auth")

	// Public routes
	authGroup.Post("/login", r.AuthHandler.Login)
	authGroup.Post("/register", r.AuthHandler.Register)
	authGroup.Post("/refresh-token", r.AuthHandler.RefreshToken)

	// Protected routes
	authGroup.Post("/logout", middleware.JWTHS256Middleware(r.RedisClient, r.AuthHandler.Logout))
	authGroup.Post("/revoke-all-tokens/:user_id", middleware.JWTHS256Middleware(r.RedisClient, r.AuthHandler.RevokeAllTokens))
	authGroup.Post("/revoke-device-tokens/:user_id", middleware.JWTHS256Middleware(r.RedisClient, r.AuthHandler.RevokeDeviceTokens)) // 🆕
}
