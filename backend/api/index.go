package api

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"ticket-booking-backend/internal/repository"
	"ticket-booking-backend/internal/routes"
)

var app *fiber.App

// init runs automatically before the Vercel handler
func init() {
	// Initialize Database Connection
	repository.ConnectDB()

	// Initialize Fiber app
	app = fiber.New(fiber.Config{
		AppName: "Ticket Booking API (Vercel Serverless)",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup API Routes
	routes.SetupRoutes(app)
}

// Handler is the Vercel Serverless Function entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
	// Adapt standard net/http handler to Fiber
	adaptor.FiberApp(app)(w, r)
}
