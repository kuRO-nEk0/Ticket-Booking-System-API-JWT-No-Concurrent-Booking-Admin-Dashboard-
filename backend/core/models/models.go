package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	Email        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// Event represents a bookable event
type Event struct {
	ID          uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	Title       string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Date        time.Time `gorm:"not null"`
	Seats       []Seat    `gorm:"foreignKey:EventID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Seat represents a single seat for an event
type Seat struct {
	ID         uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	EventID    uuid.UUID `gorm:"type:varchar(36);not null;index;uniqueIndex:idx_event_seat"`
	SeatNumber string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_event_seat"`
	Status     string    `gorm:"type:varchar(20);default:'available'"` // 'available' or 'booked'
	Version    int       `gorm:"default:1"` // Optimistic Concurrency Control
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Booking represents a user's reservation of a specific seat
type Booking struct {
	ID        uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null"`
	EventID   uuid.UUID `gorm:"type:varchar(36);not null"`
	SeatID    uuid.UUID `gorm:"type:varchar(36);not null;unique"` // A seat can only be booked once
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations for easy fetching
	User  User  `gorm:"foreignKey:UserID"`
	Event Event `gorm:"foreignKey:EventID"`
	Seat  Seat  `gorm:"foreignKey:SeatID"`
}

// BeforeCreate hooks to generate UUIDs automatically
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return
}

func (s *Seat) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

func (b *Booking) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}
