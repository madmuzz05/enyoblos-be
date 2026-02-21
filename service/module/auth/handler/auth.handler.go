package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/package/middleware"
	"github.com/madmuzz05/be-enyoblos/service/module/auth/dto"
)

// Login - User login endpoint
// @POST /auth/login
// @param LoginRequest (email, password, optional: device_id)
// @return AuthResponse
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if validationErrors, err := helper.ValidateRequest(c, &req); err != nil {
		return helper.SendResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	// 🆕 Gunakan device_id dari client jika ada, otherwise generate dari server
	deviceID := req.DeviceID
	if deviceID == "" {
		// Fallback: generate dari User-Agent + IP untuk web/browser clients
		deviceID = middleware.GenerateDeviceID(c)
	}

	res, sysErr := h.AuthUsecase.Login(c, req, deviceID)
	if sysErr != nil {
		return helper.SendErrorResponse(c, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(c, fiber.StatusOK, "Login successful", res)
}

// Register - User registration endpoint
// @POST /auth/register
// @param RegisterRequest (name, short_name, email, age, password, organization_id, optional: device_id)
// @return AuthResponse
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	if validationErrors, err := helper.ValidateRequest(c, &req); err != nil {
		return helper.SendResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	// 🆕 Gunakan device_id dari client jika ada, otherwise generate dari server
	deviceID := req.DeviceID
	if deviceID == "" {
		// Fallback: generate dari User-Agent + IP untuk web/browser clients
		deviceID = middleware.GenerateDeviceID(c)
	}

	res, sysErr := h.AuthUsecase.Register(c, req, deviceID)
	if sysErr != nil {
		return helper.SendErrorResponse(c, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(c, fiber.StatusCreated, "Registration successful", res)
}

// Logout - User logout endpoint
// @POST /auth/logout
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Missing Authorization header", nil)
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")

	if err := h.AuthUsecase.Logout(tokenStr); err != nil {
		return helper.SendErrorResponse(c, err.GetStatusCode(), err.GetMessage(), err.GetError())
	}

	return helper.SendResponse(c, fiber.StatusOK, "Logout successful", nil)
}

// RefreshToken - Refresh access token endpoint
// @POST /auth/refresh-token
// Body: {refresh_token: string, old_access_token?: string}
func (h *AuthHandler) RefreshToken(c fiber.Ctx) error {
	type RefreshTokenRequest struct {
		RefreshToken   string `json:"refresh_token" binding:"required"`
		OldAccessToken string `json:"old_access_token"` // Optional: old token to blacklist
	}

	var req RefreshTokenRequest
	if err := c.Bind().Body(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	res, sysErr := h.AuthUsecase.RefreshToken(req.RefreshToken, req.OldAccessToken)
	if sysErr != nil {
		return helper.SendErrorResponse(c, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(c, fiber.StatusOK, "Token refreshed", res)
}

// RevokeAllTokens - Revoke semua token untuk user
// Berguna untuk disable akses user atau logout dari semua device
// @POST /auth/revoke-all-tokens/:user_id
// Require: JWT Authorization
func (h *AuthHandler) RevokeAllTokens(c fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	sysErr := h.AuthUsecase.RevokeAllTokens(userID)
	if sysErr != nil {
		return helper.SendErrorResponse(c, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(c, fiber.StatusOK, "All tokens revoked successfully", nil)
}

// 🆕 RevokeDeviceTokens - Logout hanya di device tertentu
// Berguna untuk logout dari 1 device tanpa affect device lain
// @POST /auth/revoke-device-tokens/:user_id
// Body: {device_id: string}
// Require: JWT Authorization
func (h *AuthHandler) RevokeDeviceTokens(c fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID", err)
	}

	type RevokeDeviceRequest struct {
		DeviceID string `json:"device_id" binding:"required"`
	}

	var req RevokeDeviceRequest
	if err := c.Bind().Body(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	sysErr := h.AuthUsecase.RevokeDeviceTokens(userID, req.DeviceID)
	if sysErr != nil {
		return helper.SendErrorResponse(c, sysErr.GetStatusCode(), sysErr.GetMessage(), sysErr.GetError())
	}

	return helper.SendResponse(c, fiber.StatusOK, "Device tokens revoked successfully", nil)
}
