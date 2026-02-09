package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/Akshatt02/job-portal-backend/internal/config"
	"github.com/Akshatt02/job-portal-backend/internal/db"
	"github.com/Akshatt02/job-portal-backend/internal/handlers"
	"github.com/Akshatt02/job-portal-backend/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()
	db.Connect(cfg.DatabaseURL)
	defer db.Close()

	app := fiber.New()
	app.Use(logger.New())

	// Auth routes
	app.Post("/auth/register", handlers.Register)
	app.Post("/auth/login", handlers.Login)

	// Public profile read
	app.Get("/profile/:id", handlers.GetProfile)

	app.Get("/jobs", handlers.ListJobs)
	app.Get("/jobs/:id", handlers.GetJob)

	// Protected
	protected := app.Group("", middleware.AuthRequired())
	protected.Get("/me", handlers.Me)
	protected.Put("/profile", handlers.UpdateProfile)
	protected.Post("/jobs", handlers.CreateJob)

	log.Println("Starting server on port", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
