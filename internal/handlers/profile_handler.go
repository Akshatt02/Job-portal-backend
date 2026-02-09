package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Akshatt02/job-portal-backend/internal/services"
)

type updateProfileRequest struct {
	Name          *string  `json:"name,omitempty"`
	Bio           *string  `json:"bio,omitempty"`
	LinkedinURL   *string  `json:"linkedin_url,omitempty"`
	Skills        []string `json:"skills,omitempty"`
	WalletAddress *string  `json:"wallet_address,omitempty"`
}

// GET /profile/:id
func GetProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}

	return c.JSON(user)
}

// GET /me
func Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	idStr := userID.(string)

	user, err := services.GetUserByID(idStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return c.JSON(user)
}

// PUT /profile  (auth required)
func UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	idStr := userID.(string)

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.LinkedinURL != nil {
		updates["linkedin_url"] = *req.LinkedinURL
	}
	if req.WalletAddress != nil {
		updates["wallet_address"] = *req.WalletAddress
	}
	if req.Skills != nil {
		updates["skills"] = req.Skills
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no updates provided"})
	}

	if err := services.UpdateUser(idStr, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update"})
	}

	// Return new profile
	u, err := services.GetUserByID(idStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch updated user"})
	}
	return c.JSON(u)
}
