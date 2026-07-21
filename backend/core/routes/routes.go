package routes

import (
	"github.com/gofiber/fiber/v2"

	"ticket-booking-backend/core/controllers"
	"ticket-booking-backend/core/middleware"
)

// SetupRoutes configures all the application endpoints
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Public Routes
	api.Post("/auth/register", controllers.Register)
	api.Post("/auth/login", controllers.Login)
	api.Get("/events", controllers.GetEvents)
	api.Get("/events/:id/seats", controllers.GetEventSeats)

	// Protected Routes
	api.Use(middleware.Protected()) // Require JWT from here on

	// Admin Event Creation & Deletion
	api.Post("/events", controllers.CreateEvent)
	api.Delete("/events/:id", controllers.DeleteEvent)

	// Booking Routes
	bookings := api.Group("/bookings")
	bookings.Post("/", controllers.BookSeat)
	bookings.Get("/", controllers.GetMyBookings)
	bookings.Delete("/:id", controllers.CancelBooking)
}
