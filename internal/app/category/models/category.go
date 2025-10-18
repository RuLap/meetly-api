package models

import "github.com/google/uuid"

type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}
