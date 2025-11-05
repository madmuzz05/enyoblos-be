package repository

import (
	"github.com/gofiber/fiber/v2"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/organization/entity"
)

func (r *OrganizationRepository) GetOrganizations(ctx *fiber.Ctx) (res []entity.Organization, totalRecords int64, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	// Parse pagination
	pagination := helper.ParsePaginationFromQuery(ctx)
	offset := helper.GetOffset(pagination.Page, pagination.PageSize)

	// Get total records
	countQuery := `SELECT COUNT(*) FROM organizations`
	db.Raw(countQuery).Scan(&totalRecords)

	// Get paginated data
	query := `SELECT * FROM organizations ORDER BY id ` + pagination.Sort + ` LIMIT ? OFFSET ?`
	model := db.Raw(query, pagination.PageSize, offset).Scan(&res)

	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Terjadi Kesalahan Internal")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.ErrNotFound.Code, fiber.ErrNotFound.Message)
	}
	return
}

func (r *OrganizationRepository) GetOrganizationByID(ctx *fiber.Ctx, Id int) (res entity.Organization, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())
	query := `SELECT * FROM organizations WHERE id = ?`

	model := db.Raw(query, Id).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil organization")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Organization tidak ditemukan")
	}
	return
}

// CreateOrganization - Create new organization
func (r *OrganizationRepository) CreateOrganization(ctx *fiber.Ctx, req dto.CreateOrganizationRequest) (res entity.Organization, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	res = entity.Organization{
		Name:      req.Name,
		ShortName: req.ShortName,
		Address:   req.Address,
	}

	query := `INSERT INTO organizations (name, short_name, address) VALUES (?, ?, ?) RETURNING *`
	model := db.Raw(query, res.Name, res.ShortName, res.Address).Scan(&res)

	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal membuat organization")
		return
	}

	return
}

// UpdateOrganization - Update organization data
func (r *OrganizationRepository) UpdateOrganization(ctx *fiber.Ctx, id int, req dto.UpdateOrganizationRequest) (res entity.Organization, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

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

	query := `UPDATE organizations SET name = ?, short_name = ?, address = ? WHERE id = ? RETURNING *`
	model := db.Raw(query, res.Name, res.ShortName, res.Address, id).Scan(&res)

	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengupdate organization")
		return
	}

	return
}

// DeleteOrganization - Delete organization
func (r *OrganizationRepository) DeleteOrganization(ctx *fiber.Ctx, id int) (sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	query := `DELETE FROM organizations WHERE id = ?`
	model := db.Exec(query, id)

	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal menghapus organization")
		return
	}

	return
}
