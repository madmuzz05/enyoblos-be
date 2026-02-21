package middleware

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/madmuzz05/be-enyoblos/config"
	"github.com/madmuzz05/be-enyoblos/package/helper"
	"github.com/madmuzz05/be-enyoblos/package/redisdb"
)

type GenerateTokenRes struct {
	AccessToken string            `json:"access_token"`
	ExpiresIn   helper.CustomTime `json:"expires_in"`
}

// GenerateDeviceID - Generate unique device ID dari user agent + IP address
// Untuk identify device tertentu dan handle per-device logout
func GenerateDeviceID(c fiber.Ctx) string {
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// Combine user agent + IP untuk device fingerprint
	// Hash menggunakan simple string formatting
	deviceFingerprint := fmt.Sprintf("%s:%s:%d", userAgent, ipAddress, time.Now().UnixNano())

	// Create simple hash (cukup untuk fingerprint)
	hash := 0
	for _, char := range deviceFingerprint {
		hash = ((hash << 5) - hash) + int(char)
	}

	deviceID := fmt.Sprintf("%x", hash)[1:13] // Ambil 12 char
	return deviceID
}

// GenerateTokenHS256 creates a signed HS256 token with user ID, role, dan device ID
// userID = ID user yang login
// role = user role untuk authorization
// deviceID = unique device identifier (dari user agent + IP)
func GenerateTokenHS256(userID int, role string, deviceID string) (GenerateTokenRes, error) {
	ttl := time.Duration(config.AppConfig.JwtExpiresIn) * time.Second
	secret := config.AppConfig.JwtSecret
	key := config.AppConfig.JwtKey
	claims := jwt.MapClaims{
		"key":       key,
		"user_id":   userID,
		"device_id": deviceID, // 🆕 Tambah device ID untuk per-device logout
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(ttl).Unix(),
		"role":      role,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return GenerateTokenRes{}, err
	}
	expired, err := helper.ParseStringToCustomTime(time.Now().Add(ttl).Format("2006-01-02 15:04:05"))
	if err != nil {
		return GenerateTokenRes{}, err
	}
	return GenerateTokenRes{
		AccessToken: token,
		ExpiresIn:   expired,
	}, nil
}

// JWTHS256Middleware verifies token HS256 and sets claims to ctx.Locals("user_claims")
// Menerima RedisClient untuk check token blacklist
func JWTHS256Middleware(redisClient *redisdb.RedisClient, handler fiber.Handler, roles ...string) fiber.Handler {
	secret := config.AppConfig.JwtSecret
	return func(c fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "Authorization header required", nil)
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "invalid Authorization header", nil)
		}
		tokenStr := parts[1]

		// 🚫 Check token di Redis blacklist
		if redisClient != nil {
			if val, err := redisClient.Client.Get(redisClient.Ctx, "blacklist:"+tokenStr).Result(); err == nil && val == "true" {
				return helper.SendResponse(c, fiber.StatusUnauthorized, "Token revoked (logged out)", nil)
			}
		}

		// ✅ Parse token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", nil)
		}

		// attach claims
		var claims jwt.MapClaims
		if claimsData, ok := token.Claims.(jwt.MapClaims); ok {
			claims = claimsData
			c.Locals("user_claims", claims)
		} else {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "Invalid token claims", nil)
		}

		// 🚫 Check apakah user sudah di-revoke semua token-nya
		if redisClient != nil {
			if userID, ok := claims["user_id"].(float64); ok {
				revokeKey := fmt.Sprintf("revoke:user:%d", int(userID))
				if val, err := redisClient.Client.Get(redisClient.Ctx, revokeKey).Result(); err == nil && val == "true" {
					return helper.SendResponse(c, fiber.StatusUnauthorized, "User tokens have been revoked", nil)
				}

				// 🆕 Check device-specific revoke
				if deviceID, ok := claims["device_id"].(string); ok {
					deviceRevokeKey := fmt.Sprintf("revoke:user:%d:device:%s", int(userID), deviceID)
					if val, err := redisClient.Client.Get(redisClient.Ctx, deviceRevokeKey).Result(); err == nil && val == "true" {
						return helper.SendResponse(c, fiber.StatusUnauthorized, "This device has been logged out", nil)
					}
				}
			}
		}

		// role check
		if len(roles) > 0 {
			roleClaim, ok := claims["role"].(string)
			if !ok || !slices.Contains(roles, roleClaim) {
				return helper.SendResponse(c, fiber.StatusForbidden, "Forbidden: insufficient role", nil)
			}
		}

		return handler(c)
	}
}

func GenerateRefreshToken(userID int, role string, deviceID string) (GenerateTokenRes, error) {
	ttl := 7 * 24 * time.Hour                         // 7 hari
	secret := config.AppConfig.JwtSecret + "_refresh" // beda secret untuk refresh
	key := config.AppConfig.JwtKey
	claims := jwt.MapClaims{
		"key":       key,
		"user_id":   userID,
		"device_id": deviceID, // 🆕 Include device ID di refresh token juga
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(ttl).Unix(),
		"role":      role,
		"type":      "refresh",
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return GenerateTokenRes{}, err
	}
	expired, _ := helper.ParseStringToCustomTime(time.Now().Add(ttl).Format("2006-01-02 15:04:05"))
	return GenerateTokenRes{
		AccessToken: token,
		ExpiresIn:   expired,
	}, nil
}
