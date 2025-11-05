package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	// 🆕 Optional: device_id dari client (untuk aplikasi mobile/desktop)
	// Jika tidak diberikan, server akan generate dari User-Agent + IP
	// Format: UUID atau simple string identifier (max 64 chars)
	DeviceID string `json:"device_id"`
}

type RegisterRequest struct {
	Name           string `json:"name" binding:"required"`
	ShortName      string `json:"short_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Age            int    `json:"age" binding:"required,gte=0"`
	Password       string `json:"password" binding:"required,min=8"`
	OrganizationID int    `json:"organization_id" binding:"required"`
	// 🆕 Optional: device_id dari client
	DeviceID string `json:"device_id"`
}
