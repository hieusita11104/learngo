package common

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const (
	SecretKey = "your-secret-key" // Replace with your actual secret key
)

func AuthMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return fiber.NewError(fiber.StatusUnauthorized, "Authorization header missing or invalid format")
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))
	c.Locals("user_id", userID)

	return c.Next()
}
