package usecase

import (
	"github.com/gofiber/fiber/v2"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/entity"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/repository"
)

type OrganizationUsecase struct {
	organizationRepo repository.IOrganizationRepository
	redisDb          *redisdb.RedisClient
}

func InitOrganizationUsecase(organizationRepo repository.IOrganizationRepository, redisDb *redisdb.RedisClient) IOrganizationUsecase {
	return &OrganizationUsecase{
		organizationRepo: organizationRepo,
		redisDb:          redisDb,
	}
}

type IOrganizationUsecase interface {
	GetOrganizations(ctx *fiber.Ctx) (res []entity.Organization, totalRecords int64, sysError syserror.SysError)
	GetOrganizationByID(ctx *fiber.Ctx, Id int) (res entity.Organization, sysError syserror.SysError)
	CreateOrganization(ctx *fiber.Ctx, req dto.CreateOrganizationRequest) (res entity.Organization, sysError syserror.SysError)
	UpdateOrganization(ctx *fiber.Ctx, id int, req dto.UpdateOrganizationRequest) (res entity.Organization, sysError syserror.SysError)
	DeleteOrganization(ctx *fiber.Ctx, id int) (sysError syserror.SysError)
}
