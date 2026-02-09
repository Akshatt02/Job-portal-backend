package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Akshatt02/job-portal-backend/internal/services"
	// "github.com/Akshatt02/job-portal-backend/internal/models"
)

// payload for creating job
type createJobRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Skills        []string `json:"skills,omitempty"`
	Salary        string   `json:"salary,omitempty"`
	Location      string   `json:"location,omitempty"`
	PaymentTxHash string   `json:"payment_tx_hash,omitempty"`
}

// POST /jobs (protected)
func CreateJob(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	uidStr := userID.(string)

	var req createJobRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Title == "" || req.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "title and description required"})
	}

	jobID, err := services.CreateJob(req.Title, req.Description, req.Skills, req.Salary, req.Location, uidStr, req.PaymentTxHash)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Respond with created job id
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": jobID})
}

// GET /jobs
func ListJobs(c *fiber.Ctx) error {
	// optional ?limit= query can be added later
	jobs, err := services.ListJobs(100)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list jobs"})
	}
	return c.JSON(jobs)
}

// GET /jobs/:id
func GetJob(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "id required"})
	}

	job, err := services.GetJobByID(id)
	if err != nil {
		if err == services.ErrJobNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "job not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch job"})
	}
	return c.JSON(job)
}
