package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Akshatt02/job-portal-backend/internal/config"
	"github.com/Akshatt02/job-portal-backend/pkg/utils"
)

// AuthRequired verifies the Authorization header token and sets user_id in locals.
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header"})
		}

		tokenStr := parts[1]
		cfg := config.LoadConfig()
		userID, err := utils.ParseToken(tokenStr, cfg.JWTSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// put user id into locals for handlers
		c.Locals("user_id", userID)
		return c.Next()
	}
}
