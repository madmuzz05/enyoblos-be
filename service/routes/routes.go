package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	authHandler "github.com/madmuzz05/be-enyoblos/service/module/auth/handler"
	authUsecase "github.com/madmuzz05/be-enyoblos/service/module/auth/usecase"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/handler"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/repository"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/usecase"
	userRepository "github.com/madmuzz05/be-enyoblos/service/module/user/repository"
	userUsecase "github.com/madmuzz05/be-enyoblos/service/module/user/usecase"
)

func SetupRoutes(app *fiber.App) *fiber.App {

	// Global middlewares
	app.Use(middleware.NewRateLimiter())

	// set up middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},
	}))

	return app
}

func InitRoutes(app *fiber.App, db *database.MainDB, redisDb *redisdb.RedisClient) *fiber.App {
	router := SetupRoutes(app)
	api := router.Group("/api/v1")

	// Initialize Organization
	orgRepository := repository.InitOrganizationRepository(db)
	orgUsecase := usecase.InitOrganizationUsecase(orgRepository, redisDb, db)
	orgHandler := handler.InitOrganizationHandler(orgUsecase)

	// Initialize User
	userRepo := userRepository.InitUserRepository(db)
	userUC := userUsecase.InitUserUsecase(userRepo, orgUsecase, redisDb, db)

	// Initialize Auth
	authUC := authUsecase.InitAuthUsecase(redisDb, userUC)
	authHdl := authHandler.InitAuthHandler(authUC)

	InitAuthRoutes(api, authHdl, redisDb).Routes()
	InitOrganizationRoutes(api, orgHandler, redisDb).Routes()
	// define your routes here

	return router
}
