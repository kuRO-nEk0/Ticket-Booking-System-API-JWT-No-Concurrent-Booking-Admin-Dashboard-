package controllers

import (
	"github.com/gofiber/fiber/v2"

	"ticket-booking-backend/internal/models"
	"ticket-booking-backend/internal/repository"
)

// GetEvents returns a list of all events
func GetEvents(c *fiber.Ctx) error {
	var events []models.Event
	if err := repository.DB.Find(&events).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch events"})
	}

	return c.JSON(events)
}

// GetEventSeats returns all seats and their availability for a specific event
func GetEventSeats(c *fiber.Ctx) error {
	eventID := c.Params("id")

	var seats []models.Seat
	if err := repository.DB.Where("event_id = ?", eventID).Find(&seats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch seats"})
	}

	return c.JSON(seats)
}
