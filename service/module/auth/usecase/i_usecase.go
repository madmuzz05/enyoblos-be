package usecase

import (
	"github.com/gofiber/fiber/v3"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
	"github.com/madmuzz05/be-enyoblos/service/module/auth/dto"
	"github.com/madmuzz05/be-enyoblos/service/module/user/usecase"
)

type AuthUsecase struct {
	redisDb     *redisdb.RedisClient
	userUsecase usecase.IUserUsecase
}

func InitAuthUsecase(redisDb *redisdb.RedisClient, userUsecase usecase.IUserUsecase) IAuthUsecase {
	return &AuthUsecase{
		redisDb:     redisDb,
		userUsecase: userUsecase,
	}
}

type IAuthUsecase interface {
	Login(ctx fiber.Ctx, req dto.LoginRequest, deviceID string) (res dto.AuthResponse, sysError syserror.SysError)
	Register(ctx fiber.Ctx, req dto.RegisterRequest, deviceID string) (res dto.AuthResponse, sysError syserror.SysError)
	Logout(tokenStr string) (sysError syserror.SysError)
	RefreshToken(tokenStr string, oldAccessToken string) (res middleware.GenerateTokenRes, sysError syserror.SysError)
	RevokeAllTokens(userID int) (sysError syserror.SysError)
	RevokeDeviceTokens(userID int, deviceID string) (sysError syserror.SysError)
}
