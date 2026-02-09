package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/Akshatt02/job-portal-backend/internal/services"
)

type extractSkillsRequest struct {
	Bio string `json:"bio"`
}

// POST /ai/extract-skills  (auth required)
func ExtractSkills(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	uidStr := userID.(string)

	var req extractSkillsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Bio == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bio is required"})
	}

	skills, err := services.ExtractSkillsFromText(context.Background(), req.Bio)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ai extraction failed", "details": err.Error()})
	}

	// update user skills in DB
	updates := map[string]interface{}{"skills": skills}
	if err := services.UpdateUser(uidStr, updates); err != nil {
		// still return the skills even if DB update fails
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user skills", "details": err.Error()})
	}

	return c.JSON(fiber.Map{"skills": skills})
}
