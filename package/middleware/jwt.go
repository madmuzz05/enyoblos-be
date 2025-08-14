package middleware

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/madmuzz05/be-enyoblos/config"
	"github.com/madmuzz05/be-enyoblos/package/helper"
)

type GenerateTokenRes struct {
	AccessToken string            `json:"access_token"`
	ExpiresIn   helper.CustomTime `json:"expires_in"`
}

// GenerateTokenHS256 creates a signed HS256 token with "sub" claim as subject.
// Returns token string and error.
func GenerateTokenHS256(role string) (GenerateTokenRes, error) {
	ttl := time.Duration(config.AppConfig.JwtExpiresIn) * time.Hour
	secret := config.AppConfig.JwtSecret
	key := config.AppConfig.JwtKey
	claims := jwt.MapClaims{
		"key":  key,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(ttl).Unix(),
		"role": role,
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
func JWTHS256Middleware(handler fiber.Handler, roles ...string) fiber.Handler {
	secret := config.AppConfig.JwtSecret
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "Authorization header required", nil)
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "invalid Authorization header", nil)
		}
		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// ensure HMAC and HS256
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			if alg, ok := t.Header["alg"].(string); !ok || alg != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing alg: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "Invalid or expired token", nil)
		}
		// attach claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_claims", claims)
		}
		claims := c.Locals("user_claims")
		if claims == nil {
			return helper.SendResponse(c, fiber.StatusUnauthorized, "Unauthorized: no claims found", nil)
		}
		userClaims := claims.(jwt.MapClaims)
		roleClaim, ok := userClaims["role"].(string)
		if !ok {
			return helper.SendResponse(c, fiber.StatusForbidden, "Forbidden: role claim missing or invalid", nil)
		}
		allowed := slices.Contains(roles, roleClaim)
		if !allowed {
			return helper.SendResponse(c, fiber.StatusForbidden, "Forbidden: insufficient role", nil)
		}
		return handler(c)
	}
}
