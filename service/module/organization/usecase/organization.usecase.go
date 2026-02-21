package usecase

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/entity"
)

func (u *OrganizationUsecase) GetOrganizations(ctx fiber.Ctx) (res []entity.Organization, totalRecords int64, sysError syserror.SysError) {
	res, totalRecords, sysError = u.organizationRepo.GetOrganizations(ctx)
	return
}

func (u *OrganizationUsecase) GetOrganizationByID(ctx fiber.Ctx, Id int) (res entity.Organization, sysError syserror.SysError) {
	res, sysError = u.organizationRepo.GetOrganizationByID(ctx, Id)
	return
}

// CreateOrganization - Create new organization
func (u *OrganizationUsecase) CreateOrganization(ctx fiber.Ctx, req dto.CreateOrganizationRequest) (res entity.Organization, sysError syserror.SysError) {
	tx, errTx := database.TxCreate(ctx, u.mainDB.DB)
	if errTx != nil {
		sysError = errTx
		return
	}
	defer func() {
		database.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	res, sysError = u.organizationRepo.CreateOrganization(ctx, req)
	return
}

// UpdateOrganization - Update organization
func (u *OrganizationUsecase) UpdateOrganization(ctx fiber.Ctx, id int, req dto.UpdateOrganizationRequest) (res entity.Organization, sysError syserror.SysError) {
	tx, errTx := database.TxCreate(ctx, u.mainDB.DB)
	if errTx != nil {
		sysError = errTx
		return
	}
	defer func() {
		database.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	// check if organization exists
	_, sysError = u.organizationRepo.GetOrganizationByID(ctx, id)
	if sysError != nil {
		return
	}

	// Update organization
	res, sysError = u.organizationRepo.UpdateOrganization(ctx, id, req)
	return
}

// DeleteOrganization - Delete organization
func (u *OrganizationUsecase) DeleteOrganization(ctx fiber.Ctx, id int) (sysError syserror.SysError) {
	tx, errTx := database.TxCreate(ctx, u.mainDB.DB)
	if errTx != nil {
		sysError = errTx
		return
	}
	defer func() {
		database.TxSubmitTerr(ctx, sysError)
	}()

	if tx == nil {
		sysError = syserror.CreateError(fmt.Errorf("failed to begin transaction"), fiber.StatusInternalServerError, "Gagal memulai transaksi")
		return
	}

	// check if organization exists
	res, sysError := u.GetOrganizationByID(ctx, id)
	if sysError != nil {
		return
	}
	fmt.Println("Deleting organization:", res)

	sysError = u.organizationRepo.DeleteOrganization(ctx, id)
	return
}
