package handler

import "github.com/madmuzz05/be-enyoblos/service/module/auth/usecase"

type AuthHandler struct {
	AuthUsecase usecase.IAuthUsecase
}

func InitAuthHandler(authUsecase usecase.IAuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: authUsecase,
	}
}
