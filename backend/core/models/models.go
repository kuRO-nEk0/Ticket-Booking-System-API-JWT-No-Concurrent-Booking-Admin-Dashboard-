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
	Category    string    `gorm:"type:varchar(40);index"`
	City        string    `gorm:"type:varchar(80);index"`
	Venue       string    `gorm:"type:varchar(160)"`
	Tags        string    `gorm:"type:text"` // JSON array of tag codes
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
	Version    int       `gorm:"default:1"`                            // Optimistic Concurrency Control
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

type ChatSession struct {
	ID        uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	UserID    uuid.UUID `gorm:"type:varchar(36);index;not null"`
	Slots     string    `gorm:"type:text"` // JSON, merged slot state
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatTurn struct {
	ID            uuid.UUID `gorm:"primaryKey;type:varchar(36)"`
	SessionID     uuid.UUID `gorm:"type:varchar(36);index"`
	Role          string    `gorm:"type:varchar(12)"` // "user" | "assistant"
	Content       string    `gorm:"type:text"`
	ExtractedJSON string    `gorm:"type:text"` // raw router output, for debugging
	Branch        string    `gorm:"type:varchar(24)"`
	ShownEventIDs string    `gorm:"type:text"` // JSON array; resolves "the jazz one"
	LatencyMS     int
	CreatedAt     time.Time
}

type BookingDraft struct {
	ID          uuid.UUID `gorm:"primaryKey;type:varchar(36)"` // idempotency key
	UserID      uuid.UUID `gorm:"type:varchar(36);index;not null"`
	EventID     uuid.UUID `gorm:"type:varchar(36);not null"`
	SeatIDs     string    `gorm:"type:text;not null"` // JSON array of seat UUIDs
	Quantity    int       `gorm:"not null"`
	ExpiresAt   time.Time `gorm:"index"` // 10 minutes
	ConfirmedAt *time.Time
	CreatedAt   time.Time
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

func (cs *ChatSession) BeforeCreate(tx *gorm.DB) (err error) {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return
}

func (ct *ChatTurn) BeforeCreate(tx *gorm.DB) (err error) {
	if ct.ID == uuid.Nil {
		ct.ID = uuid.New()
	}
	return
}

func (bd *BookingDraft) BeforeCreate(tx *gorm.DB) (err error) {
	if bd.ID == uuid.Nil {
		bd.ID = uuid.New()
	}
	return
}
