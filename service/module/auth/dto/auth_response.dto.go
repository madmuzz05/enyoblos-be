package dto

import (
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	userDTO "github.com/madmuzz05/be-enyoblos/service/module/user/dto"
)

type AuthResponse struct {
	User         *userDTO.GetUserResponse     `json:"user"`
	AccessToken  *middleware.GenerateTokenRes `json:"access_token"`
	RefreshToken *middleware.GenerateTokenRes `json:"refresh_token,omitempty"`
}
