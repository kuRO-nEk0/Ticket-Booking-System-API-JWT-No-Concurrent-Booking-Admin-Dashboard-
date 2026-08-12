package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"time"

	"ticket-booking-backend/core/repository"
	"ticket-booking-backend/core/routes"
)

func main() {
	// Initialize Database Connection
	repository.ConnectDB()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Ticket Booking API",
	})

	// Middleware
	app.Use(logger.New())  // Log HTTP requests
	app.Use(recover.New()) // Recover from panics to prevent server crash
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // In production, this should be the specific Vercel frontend URL
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Rate Limiting: 100 requests per 1 minute
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	// Setup API Routes
	routes.SetupRoutes(app)

	// Determine Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start Server
	log.Printf("Server listening on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
