package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ticket-booking-backend/core/controllers"
	"ticket-booking-backend/core/models"
	"ticket-booking-backend/core/repository"
)

func main() {
	log.Println("Starting seed script...")
	
	repository.ConnectDB()
	db := repository.DB

	// Guard: check if seeded
	var count int64
	db.Model(&models.Event{}).Where("title = ?", "Seeding Marker Event").Count(&count)
	if count > 0 {
		log.Println("Database already seeded. Exiting.")
		return
	}

	// 1. UPDATE the existing "Tech Conference 2026" event
	res := db.Exec(`UPDATE events SET category = ?, city = ?, venue = ?, tags = ? WHERE title = ?`,
		"conference", "San Francisco", "Moscone Center", `["conference", "seated"]`, "Tech Conference 2026")
	log.Printf("Updated existing Tech Conference: %d rows affected", res.RowsAffected)

	// 2. Create the marker
	db.Create(&models.Event{
		ID:          uuid.New(),
		Title:       "Seeding Marker Event",
		Description: "Do not delete. Used to prevent duplicate seeding.",
		Date:        time.Now().AddDate(1, 0, 0),
		Category:    "conference",
		City:        "Test City",
		Tags:        `[]`,
	})

	// 3. Generate events
	cities := []string{"Guwahati", "Guwahati", "Guwahati", "Mumbai", "Delhi", "Bengaluru", "Kolkata"}
	now := time.Now()

	type template struct {
		Title string
		Desc  string
		Cat   string
		Tags  []string
	}

	templates := []template{
		// Music
		{"Symphony in C", "A grand orchestral performance.", "music", []string{"classical", "seated"}},
		{"Indie Rock Night", "Local indie bands performing live.", "music", []string{"indie", "standing", "18_plus"}},
		{"Electronic Dance Fest", "The biggest EDM party of the year.", "music", []string{"electronic", "standing", "late_night"}},
		{"Acoustic Evening", "Relaxing acoustic covers.", "music", []string{"acoustic", "seated", "all_ages"}},
		// Comedy
		{"Standup Special", "A night of hilarious standup comedy.", "comedy", []string{"seated", "18_plus", "late_night"}},
		{"Improv Show", "Audience-driven improv comedy.", "comedy", []string{"seated", "all_ages"}},
		{"Comedy Festival", "Featuring top comedians from around the country.", "comedy", []string{"seated", "outdoor"}},
		{"Late Night Laughs", "Uncensored comedy night.", "comedy", []string{"seated", "18_plus", "late_night"}},
		// Theatre
		{"Shakespeare in the Park", "A classic tragedy performed outdoors.", "theatre", []string{"outdoor", "seated", "family_friendly"}},
		{"Broadway Musical", "A spectacular musical production.", "theatre", []string{"seated", "all_ages"}},
		{"Modern Play", "A thought-provoking contemporary drama.", "theatre", []string{"seated", "18_plus"}},
		{"Puppet Show", "A magical theatre experience for kids.", "theatre", []string{"seated", "family_friendly"}},
		// Sports
		{"Local Football Derby", "High-stakes local football match.", "sports", []string{"outdoor", "standing", "family_friendly"}},
		{"Basketball Championship", "The final showdown.", "sports", []string{"seated", "family_friendly"}},
		{"Marathon Event", "Annual city marathon.", "sports", []string{"outdoor", "all_ages"}},
		{"Tennis Finals", "Thrilling tennis match.", "sports", []string{"outdoor", "seated"}},
		// Workshop
		{"Photography Masterclass", "Learn from the pros.", "workshop", []string{"seated", "all_ages"}},
		{"Pottery Workshop", "Hands-on pottery making class.", "workshop", []string{"seated", "family_friendly"}},
		{"Coding Bootcamp", "Intensive web development weekend.", "workshop", []string{"seated", "18_plus"}},
		{"Creative Writing", "Unlock your storytelling potential.", "workshop", []string{"seated", "all_ages"}},
		// Conference
		{"Tech Summit", "Future of AI and Web3.", "conference", []string{"seated", "18_plus"}},
		{"Business Expo", "Networking and industry trends.", "conference", []string{"seated", "18_plus"}},
		{"Medical Symposium", "Advances in modern medicine.", "conference", []string{"seated"}},
		{"Design Conference", "UX/UI and graphic design talks.", "conference", []string{"seated"}},
		// Film
		{"Indie Film Festival", "Showcasing independent cinema.", "film", []string{"seated", "18_plus"}},
		{"Classic Movie Night", "Screening a timeless classic.", "film", []string{"seated", "family_friendly"}},
		{"Documentary Premiere", "Eye-opening new documentary.", "film", []string{"seated", "18_plus"}},
		{"Outdoor Cinema", "Watch a blockbuster under the stars.", "film", []string{"outdoor", "seated", "family_friendly"}},
		// Festival
		{"Food Truck Festival", "Taste flavors from around the world.", "festival", []string{"outdoor", "family_friendly", "standing"}},
		{"Spring Carnival", "Rides, games, and fun.", "festival", []string{"outdoor", "family_friendly", "all_ages"}},
		{"Literary Festival", "Meet your favorite authors.", "festival", []string{"seated", "family_friendly"}},
		{"Cultural Fest", "Celebrating local art and traditions.", "festival", []string{"outdoor", "family_friendly", "standing"}},
	}

	// Dates: 4 in the past, rest spread over 8 weeks, many on weekends
	var dates []time.Time
	// 4 in the past
	for i := 0; i < 4; i++ {
		dates = append(dates, now.AddDate(0, 0, -(i+1)*5)) // -5, -10, -15, -20 days
	}
	// The rest in the future (8 weeks = 56 days)
	for i := 0; i < 34; i++ {
		// heavily weight weekends
		offset := rand.Intn(56) + 1
		t := now.AddDate(0, 0, offset)
		if rand.Float32() < 0.6 {
			// Try to shift to weekend
			daysToWeekend := 6 - int(t.Weekday()) // Shift towards Saturday
			if daysToWeekend < 0 {
				daysToWeekend = 0
			}
			t = t.AddDate(0, 0, daysToWeekend)
		}
		dates = append(dates, t)
	}

	// Capacities (20-300)
	capacities := []int{20, 50, 100, 150, 200, 300}

	eventsCreated := 0

	// Ambiguity fixtures
	tagsJSON1, _ := json.Marshal([]string{"jazz", "seated", "acoustic"})
	createEvent(db, "Midnight Jazz Quartet", "A smooth evening of jazz music.", "music", "Guwahati", "The Blue Room", string(tagsJSON1), dates[4], 50, false, false)
	
	tagsJSON2, _ := json.Marshal([]string{"jazz", "outdoor"})
	createEvent(db, "Sunday Jazz & Blues", "Jazz fusion and classic blues.", "music", "Guwahati", "Open Air Theatre", string(tagsJSON2), dates[5], 100, false, false)
	eventsCreated += 2

	// Sold out trackers
	soldOutCount := 0
	nearSoldOutCount := 0

	dateIdx := 0

	for _, tmpl := range templates {
		city := cities[rand.Intn(len(cities))]
		capacity := capacities[rand.Intn(len(capacities))]
		tagsJSON, _ := json.Marshal(tmpl.Tags)

		isPast := false
		if dateIdx < len(dates) && dates[dateIdx].Before(now) {
			isPast = true
		}

		soldOut := false
		nearSoldOut := false

		// Make 2 fully sold out (but they must be in the future to matter for availability filtering)
		if !isPast && soldOutCount < 2 {
			soldOut = true
			soldOutCount++
		} else if !isPast && nearSoldOutCount < 3 {
			nearSoldOut = true
			nearSoldOutCount++
		}

		createEvent(db, tmpl.Title, tmpl.Desc, tmpl.Cat, city, fmt.Sprintf("%s Venue", city), string(tagsJSON), dates[dateIdx], capacity, soldOut, nearSoldOut)
		dateIdx++
		eventsCreated++
	}

	log.Printf("Seeding complete! Generated %d events.", eventsCreated)
	reportDistribution(db)
}

func createEvent(db *gorm.DB, title, desc, cat, city, venue, tags string, date time.Time, capacity int, soldOut bool, nearSoldOut bool) {
	eventID := uuid.New()
	event := models.Event{
		ID:          eventID,
		Title:       title,
		Description: desc,
		Date:        date,
		Category:    cat,
		City:        city,
		Venue:       venue,
		Tags:        tags,
	}

	tx := db.Begin()
	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to create event %s: %v", title, err)
		return
	}

	// Use existing automated seat-generation path
	if err := controllers.GenerateSeats(tx, eventID, capacity); err != nil {
		tx.Rollback()
		log.Printf("Failed to generate seats for %s: %v", title, err)
		return
	}

	// Handle sold out / near sold out
	if soldOut || nearSoldOut {
		var limit int
		if soldOut {
			limit = capacity
		} else {
			limit = capacity - (rand.Intn(3) + 2) // leaves 2-4 seats
		}

		if limit > 0 {
			// Find seats and book them
			var seats []models.Seat
			tx.Where("event_id = ?", eventID).Limit(limit).Find(&seats)

			// We need a dummy user for the bookings
			var dummyUser models.User
			tx.FirstOrCreate(&dummyUser, models.User{
				ID:           uuid.New(),
				Email:        fmt.Sprintf("dummy_%s@test.com", uuid.New().String()[:8]),
				PasswordHash: "dummy",
			})

			for _, s := range seats {
				tx.Model(&s).Updates(map[string]interface{}{"status": "booked", "version": s.Version + 1})
				tx.Create(&models.Booking{
					ID:      uuid.New(),
					UserID:  dummyUser.ID,
					EventID: eventID,
					SeatID:  s.ID,
				})
			}
		}
	}

	tx.Commit()
}

func reportDistribution(db *gorm.DB) {
	type Result struct {
		Key   string
		Count int
	}

	var results []Result
	db.Model(&models.Event{}).Select("category as key, count(*) as count").Where("title != 'Seeding Marker Event'").Group("category").Scan(&results)
	fmt.Println("\n--- Events by Category ---")
	for _, r := range results {
		fmt.Printf("%-15s : %d\n", r.Key, r.Count)
	}

	db.Model(&models.Event{}).Select("city as key, count(*) as count").Where("title != 'Seeding Marker Event'").Group("city").Scan(&results)
	fmt.Println("\n--- Events by City ---")
	for _, r := range results {
		fmt.Printf("%-15s : %d\n", r.Key, r.Count)
	}

	// Past vs Future
	var pastCount, futureCount int64
	db.Model(&models.Event{}).Where("date < ? AND title != 'Seeding Marker Event'", time.Now()).Count(&pastCount)
	db.Model(&models.Event{}).Where("date >= ? AND title != 'Seeding Marker Event'", time.Now()).Count(&futureCount)
	fmt.Printf("\nPast Events: %d | Future Events: %d\n", pastCount, futureCount)

	// Sold out counts
	var events []models.Event
	db.Preload("Seats").Where("title != 'Seeding Marker Event'").Find(&events)

	fullySoldOut := 0
	nearSoldOut := 0

	for _, e := range events {
		avail := 0
		for _, s := range e.Seats {
			if s.Status == "available" {
				avail++
			}
		}
		if avail == 0 && len(e.Seats) > 0 {
			fullySoldOut++
		} else if avail >= 2 && avail <= 4 {
			nearSoldOut++
		}
	}

	fmt.Printf("\nFully Sold Out Events: %d\n", fullySoldOut)
	fmt.Printf("Near Sold Out Events (2-4 seats left): %d\n\n", nearSoldOut)
}
