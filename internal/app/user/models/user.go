package models

import (
	"github.com/google/uuid"
	"time"
)

type Gender string

const (
	GenderMale   Gender = "Мужской"
	GenderFemale Gender = "Женский"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	BirthDate time.Time `db:"birth_date"`
	Gender    Gender    `db:"Gender"`
	AvatarUrl string    `db:"avatar_url"`
}
