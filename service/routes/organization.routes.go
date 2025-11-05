package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/handler"
)

type organizationRoutes struct {
	Handler     *handler.OrganizationHandler
	Router      fiber.Router
	RedisClient *redisdb.RedisClient
}

func InitOrganizationRoutes(router fiber.Router, orgHandler *handler.OrganizationHandler, redis *redisdb.RedisClient) *organizationRoutes {
	return &organizationRoutes{
		Handler:     orgHandler,
		Router:      router,
		RedisClient: redis,
	}
}

func (r *organizationRoutes) Routes() {
	router := r.Router
	org := router.Group("/organization")

	// GET /organization - Get all organizations (public)
	org.Get("/", r.Handler.GetOrganizations)

	// GET /organization/:id - Get organization by ID (public)
	org.Get("/:id", r.Handler.GetOrganizationByID)

	// ============ Protected Routes (requires JWT) ============

	// POST /organization - Create new organization
	org.Post("/", middleware.JWTHS256Middleware(r.RedisClient, r.Handler.CreateOrganization))

	// PUT /organization/:id - Update organization
	org.Put("/:id", middleware.JWTHS256Middleware(r.RedisClient, r.Handler.UpdateOrganization))

	// DELETE /organization/:id - Delete organization
	org.Delete("/:id", middleware.JWTHS256Middleware(r.RedisClient, r.Handler.DeleteOrganization))
}
