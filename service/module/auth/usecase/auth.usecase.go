package usecase

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/madmuzz05/be-enyoblos/config"
	syserror "github.com/madmuzz05/be-enyoblos/package/error"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	"github.com/madmuzz05/be-enyoblos/service/module/auth/dto"
	userDTO "github.com/madmuzz05/be-enyoblos/service/module/user/dto"
	"golang.org/x/crypto/bcrypt"
)

func (u *AuthUsecase) Logout(tokenStr string) (sysError syserror.SysError) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return syserror.CreateError(err, fiber.StatusInternalServerError, err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))
	duration := time.Until(time.Unix(exp, 0))
	userID, _ := claims["user_id"].(float64)
	deviceID, _ := claims["device_id"].(string) // 🆕 Get device ID dari token

	// 🧹 Blacklist access token dengan TTL = sisa waktu token
	err = u.redisDb.Client.Set(u.redisDb.Ctx, "blacklist:"+tokenStr, "true", duration).Err()
	if err != nil {
		return syserror.CreateError(err, fiber.StatusInternalServerError, err.Error())
	}

	// 🆕 Invalidate HANYA di device ini (bukan semua device)
	// Set device-specific logout key untuk user ini = waktu logout
	// Ini akan mencegah refresh token dari device ini saja untuk digunakan
	deviceLogoutKey := fmt.Sprintf("revoke:user:%d:device:%s", int(userID), deviceID)
	err = u.redisDb.Client.Set(u.redisDb.Ctx, deviceLogoutKey, "true", 7*24*time.Hour).Err()
	if err != nil {
		return syserror.CreateError(err, fiber.StatusInternalServerError, err.Error())
	}

	return
}

// Login - Authenticate user dengan email dan password
func (u *AuthUsecase) Login(ctx *fiber.Ctx, req dto.LoginRequest, deviceID string) (res dto.AuthResponse, sysError syserror.SysError) {
	// Get user by email
	userRes, err := u.userUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err.GetStatusCode() == fiber.StatusInternalServerError {
			sysError = err
			return
		}
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Email tidak ditemukan")
		return
	}

	password := req.Password // password placeholder untuk login verification
	userPassword, err := u.userUsecase.GetPasswordById(ctx, userRes.ID)
	if err != nil {
		sysError = err
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Password salah")
		return
	}

	// 🆕 Generate access token dengan userID dan deviceID
	accessToken, tokenErr := middleware.GenerateTokenHS256(userRes.ID, "user", deviceID)
	if tokenErr != nil {
		sysError = syserror.CreateError(tokenErr, fiber.StatusInternalServerError, "Gagal generate token")
		return
	}

	// 🆕 Generate refresh token dengan userID dan deviceID
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userRes.ID, "user", deviceID)
	if refreshErr != nil {
		sysError = syserror.CreateError(refreshErr, fiber.StatusInternalServerError, "Gagal generate refresh token")
		return
	}

	// 🧹 Clear device-specific logout marker ketika user login kembali
	// Ini memungkinkan user untuk refresh token setelah re-login di device yang sama
	deviceLogoutKey := fmt.Sprintf("revoke:user:%d:device:%s", userRes.ID, deviceID)
	u.redisDb.Client.Del(u.redisDb.Ctx, deviceLogoutKey)

	res = dto.AuthResponse{
		User:         &userRes,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}

	return
}

// Register - Create new user account
func (u *AuthUsecase) Register(ctx *fiber.Ctx, req dto.RegisterRequest, deviceID string) (res dto.AuthResponse, sysError syserror.SysError) {
	// Create user menggunakan UserUsecase
	userReq := userDTO.CreateUserRequest{
		Name:           req.Name,
		ShortName:      req.ShortName,
		Email:          req.Email,
		Age:            req.Age,
		Password:       req.Password,
		OrganizationID: req.OrganizationID,
	}

	userRes, err := u.userUsecase.CreateUser(ctx, userReq)
	if err != nil {
		sysError = err
		return
	}

	// 🆕 Generate access token dengan userID dan deviceID
	accessToken, tokenErr := middleware.GenerateTokenHS256(userRes.ID, "user", deviceID)
	if tokenErr != nil {
		sysError = syserror.CreateError(tokenErr, fiber.StatusInternalServerError, "Gagal generate token")
		return
	}

	// 🆕 Generate refresh token dengan userID dan deviceID
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userRes.ID, "user", deviceID)
	if refreshErr != nil {
		sysError = syserror.CreateError(refreshErr, fiber.StatusInternalServerError, "Gagal generate refresh token")
		return
	}

	// 🧹 Clear device-specific logout marker untuk new user (safe)
	deviceLogoutKey := fmt.Sprintf("revoke:user:%d:device:%s", userRes.ID, deviceID)
	u.redisDb.Client.Del(u.redisDb.Ctx, deviceLogoutKey)

	res = dto.AuthResponse{
		User:         &userRes,
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}

	return
}

// RefreshToken - Generate new access token dari refresh token
// oldAccessToken = old access token yang ingin di-blacklist (optional)
func (u *AuthUsecase) RefreshToken(tokenStr string, oldAccessToken string) (res middleware.GenerateTokenRes, sysError syserror.SysError) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JwtSecret + "_refresh"), nil
	})

	// ❌ Check explicit token expiry
	if err != nil {
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Invalid refresh token")
		return
	}

	if !token.Valid {
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Refresh token expired or invalid")
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	// ❌ Explicit check expiry time
	exp, ok := claims["exp"].(float64)
	if !ok {
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Invalid token expiry claim")
		return
	}

	if time.Now().Unix() > int64(exp) {
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Refresh token has expired")
		return
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "Token is not a refresh token")
		return
	}

	// 🆕 Get device ID dari token untuk validate device consistency
	deviceID, _ := claims["device_id"].(string)
	userID, _ := claims["user_id"].(float64)

	// 🆕 Check apakah device ini sudah di-revoke
	if deviceID != "" {
		deviceRevokeKey := fmt.Sprintf("revoke:user:%d:device:%s", int(userID), deviceID)
		if val, err := u.redisDb.Client.Get(u.redisDb.Ctx, deviceRevokeKey).Result(); err == nil && val == "true" {
			sysError = syserror.CreateError(fiber.ErrUnauthorized, fiber.StatusUnauthorized, "This device has been logged out")
			return
		}
	}

	// 🧹 Blacklist old access token jika diberikan
	// Ini mencegah token lama untuk digunakan setelah refresh
	if oldAccessToken != "" {
		oldToken, err := jwt.Parse(oldAccessToken, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JwtSecret), nil
		})

		if err == nil && oldToken.Valid {
			oldClaims := oldToken.Claims.(jwt.MapClaims)
			// Verify oldToken adalah milik user yang sama
			oldUserID, _ := oldClaims["user_id"].(float64)
			if int(oldUserID) == int(userID) {
				// Blacklist old token
				oldExp := int64(oldClaims["exp"].(float64))
				oldDuration := time.Until(time.Unix(oldExp, 0))
				if oldDuration > 0 {
					u.redisDb.Client.Set(u.redisDb.Ctx, "blacklist:"+oldAccessToken, "true", oldDuration)
				}
			}
		}
	}

	// Generate new access token dengan userID dan deviceID yang sama
	role, _ := claims["role"].(string)
	res, tokenErr := middleware.GenerateTokenHS256(int(userID), role, deviceID)
	if tokenErr != nil {
		sysError = syserror.CreateError(tokenErr, fiber.StatusInternalServerError, "Gagal generate token")
		return
	}

	return
}

// RevokeAllTokens - Revoke semua token untuk user tertentu
// Berguna untuk logout dari semua device atau disable akses user
// userID = ID user yang akan di-revoke tokennya
func (u *AuthUsecase) RevokeAllTokens(userID int) (sysError syserror.SysError) {
	// Set marker di Redis untuk user ini
	// Setiap kali user login, generate token baru dengan user ID
	// Middleware akan check apakah user sudah di-revoke atau tidak

	revokeKey := fmt.Sprintf("revoke:user:%d", userID)
	ttl := 7 * 24 * time.Hour // Keep revoke marker untuk 7 hari

	// Set flag di Redis
	err := u.redisDb.Client.Set(u.redisDb.Ctx, revokeKey, "true", ttl).Err()
	if err != nil {
		return syserror.CreateError(err, fiber.StatusInternalServerError, "Gagal revoke token user")
	}

	return nil
}

// 🆕 RevokeDeviceTokens - Logout hanya di device tertentu
// Berguna untuk logout dari 1 device tanpa affect device lain
// userID = ID user
// deviceID = device yang ingin di-logout
func (u *AuthUsecase) RevokeDeviceTokens(userID int, deviceID string) (sysError syserror.SysError) {
	// Set device-specific revoke flag di Redis
	deviceRevokeKey := fmt.Sprintf("revoke:user:%d:device:%s", userID, deviceID)
	ttl := 7 * 24 * time.Hour // Keep untuk 7 hari

	// Set flag di Redis
	err := u.redisDb.Client.Set(u.redisDb.Ctx, deviceRevokeKey, "true", ttl).Err()
	if err != nil {
		return syserror.CreateError(err, fiber.StatusInternalServerError, "Gagal revoke token device")
	}

	return nil
}
