package event

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID  `db:"id"`
	CreatorID   uuid.UUID  `db:"creator_id"`
	CategoryID  uuid.UUID  `db:"category_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Latitude    float64    `db:"latitude"`
	Longitude   float64    `db:"longitude"`
	Address     *string    `db:"address"`
	StartsAt    *time.Time `db:"starts_at"`
	EndsAt      *time.Time `db:"ends_at"`
	IsPublic    bool       `db:"is_public"`
	CreatedAt   time.Time  `db:"created_at"`
}

type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

type Participant struct {
	UserID   uuid.UUID `db:"user_id"`
	EventID  uuid.UUID `db:"event_id"`
	JoinedAt time.Time `db:"joined_at"`
}
