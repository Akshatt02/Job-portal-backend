package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Akshatt02/job-portal-backend/internal/services"
	"github.com/Akshatt02/job-portal-backend/pkg/utils"
	"github.com/Akshatt02/job-portal-backend/internal/config"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name, email and password required"})
	}

	id, err := services.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	cfg := config.LoadConfig()
	token, err := utils.GenerateJWT(id, cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create token"})
	}

	return c.JSON(fiber.Map{"token": token})
}

func Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email and password required"})
	}

	id, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	cfg := config.LoadConfig()
	token, err := utils.GenerateJWT(id, cfg.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
