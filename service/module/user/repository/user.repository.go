package repository

import (
	"github.com/gofiber/fiber/v3"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/service/module/user/entity"
)

func (r *UserRepository) CreateUser(ctx fiber.Ctx, user entity.User) (res entity.User, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	query := `INSERT INTO public.users (name, short_name, email, age, password, organization_id) 
             VALUES ($1, $2, $3, $4, $5, $6) 
             RETURNING id, name, short_name, email, age, password, organization_id`

	model := db.Get(&res, query, user.Name, user.ShortName, user.Email, user.Age, user.Password, user.OrganizationID)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal membuat user")
		return
	}
	return
}

func (r *UserRepository) GetUserByID(ctx fiber.Ctx, Id string) (res entity.User, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}
	query := `SELECT id, name, short_name, email, age, password, organization_id FROM public.users WHERE id = $1`

	model := db.Get(&res, query, Id)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil user")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetUsers(ctx fiber.Ctx) (res []entity.User, totalRecords int64, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	// Parse pagination
	pagination := helper.ParsePaginationFromQuery(ctx)
	offset := helper.GetOffset(pagination.Page, pagination.PageSize)

	// Get total records
	countQuery := `SELECT COUNT(*) FROM public.users`
	model := db.Get(&totalRecords, countQuery)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil total records")
		return
	}

	// Get paginated data
	query := `SELECT id, name, short_name, email, age, password, organization_id 
	          FROM public.users 
	          LIMIT $1 OFFSET $2`

	model = db.Select(&res, query, pagination.PageSize, offset)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil users")
		return
	} else if len(res) == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Tidak ada users")
	}
	return
}

func (r *UserRepository) DeleteUser(ctx fiber.Ctx, Id string) syserror.SysError {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	query := `DELETE FROM public.users WHERE id = $1`

	result, model := db.Exec(query, Id)
	rowsAffected, _ := result.RowsAffected()
	if model != nil {
		return syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal menghapus user")
	} else if rowsAffected == 0 {
		return syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx fiber.Ctx, Id string, user entity.User) (res entity.User, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	query := `UPDATE public.users 
			 SET name = $1, short_name = $2, email = $3, age = $4, password = $5, organization_id = $6 
			 WHERE id = $7 
			 RETURNING id, name, short_name, email, age, password, organization_id`
	model := db.Get(&res, query, user.Name, user.ShortName, user.Email, user.Age, user.Password, user.OrganizationID, Id)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal memperbarui user")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) UpdateProfileUser(ctx fiber.Ctx, Id string, user entity.User) (res entity.User, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	query := `UPDATE public.users 
			 SET name = $1, short_name = $2, email = $3, age = $4 
			 WHERE id = $5 
			 RETURNING id, name, short_name, email, age, password, organization_id`
	model := db.Get(&res, query, user.Name, user.ShortName, user.Email, user.Age, Id)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal memperbarui profil user")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetUsersByOrganizationID(ctx fiber.Ctx, organizationID string) (res []entity.User, totalRecords int64, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}

	// Parse pagination
	pagination := helper.ParsePaginationFromQuery(ctx)
	offset := helper.GetOffset(pagination.Page, pagination.PageSize)

	// Get total records
	countQuery := `SELECT COUNT(*) FROM public.users WHERE organization_id = $1`
	model := db.Get(&totalRecords, countQuery, organizationID)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil total records")
		return
	}

	// Get paginated data
	query := `SELECT id, name, short_name, email, age, password, organization_id 
	          FROM public.users 
	          WHERE organization_id = $1 
	          LIMIT $2 OFFSET $3`

	model = db.Select(&res, query, organizationID, pagination.PageSize, offset)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil users berdasarkan organization_id")
		return
	} else if len(res) == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Tidak ada users untuk organization_id tersebut")
	}
	return
}

func (r *UserRepository) GetUserByEmailAndId(ctx fiber.Ctx, email string, Id string) (res entity.User, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}
	query := `SELECT id, name, short_name, email, age, password, organization_id FROM public.users WHERE email = $1 AND id != $2`

	model := db.Get(&res, query, email, Id)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil user")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetUserByEmail(ctx fiber.Ctx, email string) (res entity.User, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}
	query := `SELECT id, name, short_name, email, age, password, organization_id FROM public.users WHERE email = $1`

	model := db.Get(&res, query, email)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil user")
		return
	} else if res.ID == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetPasswordById(ctx fiber.Ctx, Id int) (password string, sysError syserror.SysError) {
	db := database.DBWithCtx{
		DB:  r.mainDB.DB,
		Ctx: ctx.Context(),
	}
	query := `SELECT password FROM public.users WHERE id = $1`

	model := db.Get(&password, query, Id)
	if model != nil {
		sysError = syserror.CreateError(model, fiber.StatusInternalServerError, "Gagal mengambil password user")
		return
	} else if password == "" {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}
