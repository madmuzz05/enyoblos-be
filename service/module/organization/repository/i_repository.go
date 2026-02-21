package repository

import (
	"github.com/gofiber/fiber/v3"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/entity"
)

type OrganizationRepository struct {
	mainDB *database.MainDB
}

func InitOrganizationRepository(mainDB *database.MainDB) IOrganizationRepository {
	return &OrganizationRepository{
		mainDB: mainDB,
	}
}

func (r *OrganizationRepository) GetMainDB(ctx fiber.Ctx) (tx interface{}) {
	return r.mainDB.DB
}

type IOrganizationRepository interface {
	GetMainDB(ctx fiber.Ctx) (tx interface{})

	GetOrganizations(ctx fiber.Ctx) (res []entity.Organization, totalRecords int64, sysError syserror.SysError)
	GetOrganizationByID(ctx fiber.Ctx, Id int) (res entity.Organization, sysError syserror.SysError)
	CreateOrganization(ctx fiber.Ctx, req dto.CreateOrganizationRequest) (res entity.Organization, sysError syserror.SysError)
	UpdateOrganization(ctx fiber.Ctx, id int, req dto.UpdateOrganizationRequest) (res entity.Organization, sysError syserror.SysError)
	DeleteOrganization(ctx fiber.Ctx, id int) (sysError syserror.SysError)
}
