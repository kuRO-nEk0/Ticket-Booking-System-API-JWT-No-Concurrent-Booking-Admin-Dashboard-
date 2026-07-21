package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"ticket-booking-backend/core/models"
)

var DB *gorm.DB

// ConnectDB establishes a connection to the database and runs auto-migrations
func ConnectDB() {
	var db *gorm.DB
	var err error

	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		// Use PostgreSQL in Production (Render automatically injects DATABASE_URL)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		log.Println("Connected Successfully to Production PostgreSQL Database")
	} else {
		// Fallback to SQLite for local development (and write to /tmp for serverless compatibility)
		db, err = gorm.Open(sqlite.Open("/tmp/ticket_booking.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		log.Println("Connected Successfully to Local SQLite Database (in /tmp)")
	}

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	DB = db

	// Migrate the schema automatically
	log.Println("Running Migrations...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Seat{},
		&models.Booking{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database. \n", err)
	}

	seedDatabase()
}

func seedDatabase() {
	var count int64
	DB.Model(&models.Event{}).Count(&count)
	if count == 0 {
		log.Println("Database is empty. Seeding test event and seats...")
		eventID := uuid.New()
		event := models.Event{
			ID:          eventID,
			Title:       "Tech Conference 2026",
			Description: "The biggest AI and Web3 summit in the valley.",
			Date:        time.Now().AddDate(0, 1, 0), // 1 month from now
		}
		DB.Create(&event)

		for i := 1; i <= 20; i++ {
			seat := models.Seat{
				ID:         uuid.New(),
				EventID:    eventID,
				SeatNumber: fmt.Sprintf("A%d", i),
				Status:     "available",
				Version:    1,
			}
			DB.Create(&seat)
		}
		log.Println("Seed complete.")
	}
}
