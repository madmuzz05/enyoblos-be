package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/madmuzz05/be-enyoblos/config"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/package/logger"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	"github.com/madmuzz05/be-enyoblos/service/routes"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load config dari .env
	err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Init logger sesuai environment
	logger.InitLogger("development")

	db, err := database.Connect(
		config.AppConfig.DatabaseHost,
		config.AppConfig.DatabaseUsername,
		config.AppConfig.DatabasePassword,
		config.AppConfig.DatabaseName,
		config.AppConfig.DatabasePort,
		config.AppConfig.DatabaseSSL,
	)
	if err != nil {
		log.Fatal().Msg("Failed to connect to database")
	}
	log.Info().Msg("Connected to database successfully")

	database.RunMigrationsPostgres(db.DB)

	redistStringConn := fmt.Sprintf("%s:%s", config.AppConfig.RedisHost, config.AppConfig.RedisPort)
	redisDb, errRedis := redisdb.InitRedis(redistStringConn, config.AppConfig.RedisPassword, 0)
	if errRedis != nil {
		log.Fatal().Err(errRedis).Msg("Failed to connect to Redis")
	}

	// Fiber app
	app := fiber.New(fiber.Config{AppName: "enyoblos"})

	// Pasang middleware logger
	app.Use(logger.NewLogger())

	// Load routes
	app = routes.InitRoutes(app, db, redisDb)

	app.Use(func(c *fiber.Ctx) error {
		for _, routes := range app.Stack() {
			for _, r := range routes {
				if r.Path == c.Path() && r.Method != c.Method() {
					return helper.SendResponse(c, fiber.StatusMethodNotAllowed, "Method not allowed", nil)
				}
			}
		}
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		return helper.SendResponse(c, fiber.StatusNotFound, "Endpoint not found", nil)
	})

	// Ambil port dari env
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	if err := app.Listen(":" + port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
