package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
)

func SetupRoutes(app *fiber.App) *fiber.App {

	// Global middlewares
	app.Use(middleware.NewRateLimiter())

	// set up middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		AllowMethods:  "*",
		AllowHeaders:  "*",
		ExposeHeaders: "*",
	}))

	return app
}

func InitRoutes(app *fiber.App) *fiber.App {

	router := SetupRoutes(app)
	api := router.Group("/api/v1")

	InitAuthRoutes(api).Routes()
	// define your routes here

	return router
}
