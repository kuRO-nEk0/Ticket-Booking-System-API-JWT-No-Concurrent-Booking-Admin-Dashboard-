package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"ticket-booking-backend/internal/models"
	"ticket-booking-backend/internal/repository"
)

type EventRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	SeatsCount  int    `json:"seats_count"`
}

// CreateEvent handles creating a new event and its corresponding seats
func CreateEvent(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	if email != "tmarked4l@gmail.com" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Super Admin access required"})
	}

	var req EventRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Basic Validation
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Event title is required"})
	}
	
	eventDate, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format. Use RFC3339."})
	}
	if eventDate.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Event date must be in the future"})
	}

	if req.SeatsCount < 1 || req.SeatsCount > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Seats count must be between 1 and 500"})
	}

	// Create Event
	eventID := uuid.New()
	event := models.Event{
		ID:          eventID,
		Title:       req.Title,
		Description: req.Description,
		Date:        eventDate,
	}

	// Use a transaction
	tx := repository.DB.Begin()

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create event"})
	}

	// Generate Seats
	var seats []models.Seat
	for i := 1; i <= req.SeatsCount; i++ {
		seats = append(seats, models.Seat{
			ID:         uuid.New(),
			EventID:    eventID,
			SeatNumber: fmt.Sprintf("A%d", i),
			Status:     "available",
			Version:    1,
		})
	}

	if err := tx.Create(&seats).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate seats"})
	}

	tx.Commit()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Event and seats created successfully",
		"event":   event,
	})
}

// DeleteEvent handles the deletion of an event and all associated seats/bookings
func DeleteEvent(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	if email != "tmarked4l@gmail.com" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Super Admin access required"})
	}

	eventIDStr := c.Params("id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	// Use a transaction
	tx := repository.DB.Begin()

	// 1. Delete associated bookings
	if err := tx.Where("event_id = ?", eventID).Delete(&models.Booking{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete associated bookings"})
	}

	// 2. Delete associated seats
	if err := tx.Where("event_id = ?", eventID).Delete(&models.Seat{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete associated seats"})
	}

	// 3. Delete the event itself
	var event models.Event
	if err := tx.Where("id = ?", eventID).Delete(&event).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete event"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{"message": "Event deleted successfully"})
}
