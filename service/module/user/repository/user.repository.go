package repository

import (
	"github.com/gofiber/fiber/v2"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/service/module/user/entity"
)

func (r *UserRepository) CreateUser(ctx *fiber.Ctx, user entity.User) (res entity.User, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	query := `INSERT INTO public.users (name, short_name, email, age, password, organization_id) 
             VALUES (?, ?, ?, ?, ?, ?) 
             RETURNING id, name, short_name, email, age, password, organization_id`

	model := db.Raw(query, user.Name, user.ShortName, user.Email, user.Age, user.Password, user.OrganizationID).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal membuat user")
		return
	}
	return
}

func (r *UserRepository) GetUserByID(ctx *fiber.Ctx, Id string) (res entity.User, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())
	query := `SELECT id, name, short_name, email, age, password, organization_id FROM public.users WHERE id = ?`

	model := db.Raw(query, Id).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil user")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetUsers(ctx *fiber.Ctx) (res []entity.User, totalRecords int64, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	// Parse pagination
	pagination := helper.ParsePaginationFromQuery(ctx)
	offset := helper.GetOffset(pagination.Page, pagination.PageSize)

	// Get total records
	countQuery := `SELECT COUNT(*) FROM public.users`
	db.Raw(countQuery).Scan(&totalRecords)

	// Get paginated data
	query := `SELECT id, name, short_name, email, age, password, organization_id 
	          FROM public.users 
	          LIMIT ? OFFSET ?`

	model := db.Raw(query, pagination.PageSize, offset).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil users")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Tidak ada users")
	}
	return
}

func (r *UserRepository) DeleteUser(ctx *fiber.Ctx, Id string) syserror.SysError {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	query := `DELETE FROM public.users WHERE id = ?`

	model := db.Exec(query, Id)
	if model.Error != nil {
		return syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal menghapus user")
	} else if model.RowsAffected == 0 {
		return syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return nil
}

func (r *UserRepository) UpdateUser(ctx *fiber.Ctx, Id string, user entity.User) (res entity.User, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	query := `UPDATE public.users 
			 SET name = ?, short_name = ?, email = ?, age = ?, password = ?, organization_id = ? 
			 WHERE id = ? 
			 RETURNING id, name, short_name, email, age, password, organization_id`
	model := db.Raw(query, user.Name, user.ShortName, user.Email, user.Age, user.Password, user.OrganizationID, Id).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal memperbarui user")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) UpdateProfileUser(ctx *fiber.Ctx, Id string, user entity.User) (res entity.User, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	query := `UPDATE public.users 
			 SET name = ?, short_name = ?, email = ?, age = ? 
			 WHERE id = ? 
			 RETURNING id, name, short_name, email, age, password, organization_id`
	model := db.Raw(query, user.Name, user.ShortName, user.Email, user.Age, Id).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal memperbarui profil user")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetUsersByOrganizationID(ctx *fiber.Ctx, organizationID string) (res []entity.User, totalRecords int64, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())

	// Parse pagination
	pagination := helper.ParsePaginationFromQuery(ctx)
	offset := helper.GetOffset(pagination.Page, pagination.PageSize)

	// Get total records
	countQuery := `SELECT COUNT(*) FROM public.users WHERE organization_id = ?`
	db.Raw(countQuery, organizationID).Scan(&totalRecords)

	// Get paginated data
	query := `SELECT id, name, short_name, email, age, password, organization_id 
	          FROM public.users 
	          WHERE organization_id = ? 
	          LIMIT ? OFFSET ?`

	model := db.Raw(query, organizationID, pagination.PageSize, offset).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil users berdasarkan organization_id")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "Tidak ada users untuk organization_id tersebut")
	}
	return
}

func (r *UserRepository) GetUserByEmailAndId(ctx *fiber.Ctx, email string, Id string) (res entity.User, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())
	query := `SELECT id, name, short_name, email, age, password, organization_id FROM public.users WHERE email = ? AND id != ?`

	model := db.Raw(query, email, Id).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil user")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetUserByEmail(ctx *fiber.Ctx, email string) (res entity.User, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())
	query := `SELECT id, name, short_name, email, age, password, organization_id FROM public.users WHERE email = ?`

	model := db.Raw(query, email).Scan(&res)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil user")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}

func (r *UserRepository) GetPasswordById(ctx *fiber.Ctx, Id int) (password string, sysError syserror.SysError) {
	db := r.mainDB.DB.WithContext(ctx.UserContext())
	query := `SELECT password FROM public.users WHERE id = ?`

	model := db.Raw(query, Id).Scan(&password)
	if model.Error != nil {
		sysError = syserror.CreateError(model.Error, fiber.StatusInternalServerError, "Gagal mengambil password user")
		return
	} else if model.RowsAffected == 0 {
		sysError = syserror.CreateError(fiber.ErrNotFound, fiber.StatusNotFound, "User tidak ditemukan")
	}
	return
}
