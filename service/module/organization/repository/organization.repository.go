package repository

import (
	"github.com/gofiber/fiber/v3"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/entity"
)

func (r *OrganizationRepository) GetOrganizations(ctx fiber.Ctx) (res []entity.Organization, totalRecords int64, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	// Parse pagination
	pagination := helper.ParsePaginationFromQuery(ctx)
	offset := helper.GetOffset(pagination.Page, pagination.PageSize)

	// Get total records
	countQuery := `SELECT COUNT(*) FROM public.organizations`
	db.Get(&totalRecords, countQuery)

	// Get paginated data
	query := `SELECT * FROM public.organizations ORDER BY id ` + pagination.Sort + ` LIMIT $1 OFFSET $2`
	err := db.Select(&res, query, pagination.PageSize, offset)

	if err != nil {
		sysError = syserror.CreateError(err, fiber.StatusInternalServerError, "Terjadi Kesalahan Internal")
		return
	} else if len(res) == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.ErrNotFound.Code, fiber.ErrNotFound.Message)
	}
	return
}

func (r *OrganizationRepository) GetOrganizationByID(ctx fiber.Ctx, Id int) (res entity.Organization, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}
	query := `SELECT * FROM public.organizations WHERE id = $1`

	model := db.Get(&res, query, Id)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil organization")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Organization tidak ditemukan")
	}
	return
}

// CreateOrganization - Create new organization
func (r *OrganizationRepository) CreateOrganization(ctx fiber.Ctx, req dto.CreateOrganizationRequest) (res entity.Organization, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	res = entity.Organization{
		Name:      req.Name,
		ShortName: req.ShortName,
		Address:   req.Address,
	}

	query := `INSERT INTO public.organizations (name, short_name, address) VALUES ($1, $2, $3) RETURNING *`
	model := db.Get(&res, query, res.Name, res.ShortName, res.Address)

	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal membuat organization")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Organization tidak ditemukan")
		return
	}

	return
}

// UpdateOrganization - Update organization data
func (r *OrganizationRepository) UpdateOrganization(ctx fiber.Ctx, id int, req dto.UpdateOrganizationRequest) (res entity.Organization, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	// Update hanya field yang diberikan
	if req.Name != "" {
		res.Name = req.Name
	}
	if req.ShortName != "" {
		res.ShortName = req.ShortName
	}
	if req.Address != "" {
		res.Address = req.Address
	}

	query := `UPDATE public.organizations SET name = $1, short_name = $2, address = $3 WHERE id = $4 RETURNING *`
	model := db.Get(&res, query, res.Name, res.ShortName, res.Address, id)

	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengupdate organization")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Organization tidak ditemukan")
		return
	}

	return
}

// DeleteOrganization - Delete organization
func (r *OrganizationRepository) DeleteOrganization(ctx fiber.Ctx, id int) (sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	query := `DELETE FROM public.organizations WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		sysError = syserror.CreateError(err, fiber.StatusInternalServerError, "Gagal menghapus organization")
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		sysError = syserror.CreateError(nil, fiber.StatusNotFound, "Organization tidak ditemukan")
		return
	}

	return
}
