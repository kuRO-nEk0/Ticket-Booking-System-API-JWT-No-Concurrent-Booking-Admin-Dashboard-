package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"ticket-booking-backend/core/models"
	"ticket-booking-backend/core/repository"
)

type BookingRequest struct {
	EventID string `json:"event_id"`
	SeatID  string `json:"seat_id"`
}

// BookSeat handles the concurrent booking logic using Optimistic Concurrency Control
func BookSeat(c *fiber.Ctx) error {
	// Extract user ID from JWT token (passed by middleware)
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Parse request body
	var req BookingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	eventID, err := uuid.Parse(req.EventID)
	seatID, err := uuid.Parse(req.SeatID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid IDs"})
	}

	// --- CONCURRENCY HANDLING ALGORITHM: OPTIMISTIC LOCKING ---
	// 1. Read the current state of the seat (No locks acquired)
	var seat models.Seat
	if err := repository.DB.Where("id = ? AND event_id = ?", seatID, eventID).First(&seat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Seat not found"})
	}

	// 2. Check Availability
	if seat.Status != "available" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Seat is already booked"})
	}

	// 3. Attempt the update with the exact version we just read
	// The database will only update the row if the version hasn't changed since we read it.
	result := repository.DB.Model(&seat).
		Where("version = ?", seat.Version).
		Updates(map[string]interface{}{
			"status":  "booked",
			"version": seat.Version + 1, // Increment version for next time
		})

	if result.Error != nil {
		log.Println("Database error during optimistic update:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "An error occurred while booking"})
	}

	// If RowsAffected is 0, it means someone else beat us to it and incremented the version!
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Seat was just booked by someone else. Please choose another seat."})
	}

	// 4. Create the Booking Record
	// Since we successfully updated the seat status, we own the seat. We can now create the record.
	booking := models.Booking{
		UserID:  userID,
		EventID: eventID,
		SeatID:  seatID,
	}

	if err := repository.DB.Create(&booking).Error; err != nil {
		log.Println("Booking creation error:", err)
		// Note: In a production system, if this fails, we should run a compensating transaction
		// to revert the seat status back to 'available'.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create booking record."})
	}

	// Preload the related Event and Seat so the frontend can immediately render the virtual ticket
	repository.DB.Preload("Event").Preload("Seat").First(&booking, booking.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Seat booked successfully!",
		"booking": booking,
	})
}

// GetMyBookings fetches all bookings for the logged in user
func GetMyBookings(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var bookings []models.Booking
	if err := repository.DB.Preload("Event").Preload("Seat").Where("user_id = ?", userID).Order("created_at desc").Find(&bookings).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch bookings"})
	}

	return c.JSON(bookings)
}

// CancelBooking handles the cancellation of a booked ticket
func CancelBooking(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userIDStr := claims["user_id"].(string)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	bookingIDStr := c.Params("id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid booking ID"})
	}

	var booking models.Booking
	if err := repository.DB.Where("id = ? AND user_id = ?", bookingID, userID).First(&booking).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Booking not found or unauthorized"})
	}

	var seat models.Seat
	if err := repository.DB.First(&seat, booking.SeatID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Seat not found"})
	}

	// Begin a transaction to safely delete booking and restore seat
	tx := repository.DB.Begin()

	// 1. Delete the booking
	if err := tx.Delete(&booking).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to cancel booking"})
	}

	// 2. Update the seat status to available and increment version for OCC
	if err := tx.Model(&seat).Updates(map[string]interface{}{
		"status":  "available",
		"version": seat.Version + 1,
	}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update seat status"})
	}

	tx.Commit()

	return c.JSON(fiber.Map{"message": "Ticket cancelled successfully"})
}
