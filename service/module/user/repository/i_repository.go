package repository

import (
	"github.com/gofiber/fiber/v2"
	database "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/service/module/user/entity"
)

type UserRepository struct {
	mainDB *database.MainDB
}

func InitUserRepository(mainDB *database.MainDB) IUserRepository {
	return &UserRepository{
		mainDB: mainDB,
	}
}

func (r *UserRepository) GetMainDB(ctx *fiber.Ctx) (tx interface{}) {
	return r.mainDB.DB
}

type IUserRepository interface {
	GetMainDB(ctx *fiber.Ctx) (tx interface{})

	CreateUser(ctx *fiber.Ctx, user entity.User) (res entity.User, sysError syserror.SysError)
	GetUserByID(ctx *fiber.Ctx, Id string) (res entity.User, sysError syserror.SysError)
	GetUsers(ctx *fiber.Ctx) (res []entity.User, totalRecords int64, sysError syserror.SysError)
	DeleteUser(ctx *fiber.Ctx, Id string) syserror.SysError
	UpdateUser(ctx *fiber.Ctx, Id string, user entity.User) (res entity.User, sysError syserror.SysError)
	UpdateProfileUser(ctx *fiber.Ctx, Id string, user entity.User) (res entity.User, sysError syserror.SysError)
	GetUsersByOrganizationID(ctx *fiber.Ctx, organizationID string) (res []entity.User, totalRecords int64, sysError syserror.SysError)
	GetUserByEmailAndId(ctx *fiber.Ctx, email string, Id string) (res entity.User, sysError syserror.SysError)
	GetUserByEmail(ctx *fiber.Ctx, email string) (res entity.User, sysError syserror.SysError)
	GetPasswordById(ctx *fiber.Ctx, Id int) (password string, sysError syserror.SysError)
}
