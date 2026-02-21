package usecase

import (
	"github.com/gofiber/fiber/v3"
	dbpostgres "github.com/madmuzz05/be-enyoblos/package/database/postgres"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	organizationUsecase "github.com/madmuzz05/be-enyoblos/service/module/organization/usecase"
	"github.com/madmuzz05/be-enyoblos/service/module/user/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/user/repository"
)

type UserUsecase struct {
	userRepo        repository.IUserRepository
	organizationUse organizationUsecase.IOrganizationUsecase
	redisDb         *redisdb.RedisClient
	mainDB          *dbpostgres.MainDB
}

func InitUserUsecase(userRepo repository.IUserRepository, organizationUse organizationUsecase.IOrganizationUsecase, redisDb *redisdb.RedisClient, mainDB *dbpostgres.MainDB) IUserUsecase {
	return &UserUsecase{
		userRepo:        userRepo,
		organizationUse: organizationUse,
		redisDb:         redisDb,
		mainDB:          mainDB,
	}
}

type IUserUsecase interface {
	CreateUser(ctx fiber.Ctx, req dto.CreateUserRequest) (res dto.GetUserResponse, sysError syserror.SysError)
	GetUserByID(ctx fiber.Ctx, Id string) (res dto.GetUserResponse, sysError syserror.SysError)
	GetUsers(ctx fiber.Ctx) (res []dto.GetUserResponse, totalRecords int64, sysError syserror.SysError)
	UpdateProfileUser(ctx fiber.Ctx, Id string, req dto.CreateUserRequest) (res dto.GetUserResponse, sysError syserror.SysError)
	GetUserByEmailAndId(ctx fiber.Ctx, email string, Id string) (res dto.GetUserResponse, sysError syserror.SysError)
	DeleteUser(ctx fiber.Ctx, Id string) (sysError syserror.SysError)
	UpdateUser(ctx fiber.Ctx, Id string, req dto.CreateUserRequest) (res dto.GetUserResponse, sysError syserror.SysError)
	GetUsersByOrganizationID(ctx fiber.Ctx, organizationID string) (res []dto.GetUserResponse, totalRecords int64, sysError syserror.SysError)
	GetUserByEmail(ctx fiber.Ctx, email string) (res dto.GetUserResponse, sysError syserror.SysError)
	GetPasswordById(ctx fiber.Ctx, Id int) (password string, sysError syserror.SysError)
}
