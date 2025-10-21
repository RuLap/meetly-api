package event

import "time"

type GetEventResponse struct {
	ID           string                   `json:"id"`
	Title        string                   `json:"title"`
	Description  string                   `json:"description"`
	Latitude     float64                  `json:"latitude"`
	Longitude    float64                  `json:"longitude"`
	Address      *string                  `json:"address"`
	StartsAt     *time.Time               `json:"starts_at"`
	EndsAt       *time.Time               `json:"ends_at"`
	IsPublic     bool                     `json:"is_public"`
	CreatedAt    time.Time                `json:"created_at"`
	Creator      GetParticipantResponse   `json:"creator"`
	Category     GetCategoryResponse      `json:"category"`
	Participants []GetParticipantResponse `json:"participants"`
}

type GetShortEventResponse struct {
	ID                string                 `json:"id"`
	Title             string                 `json:"title"`
	Latitude          float64                `json:"latitude"`
	Longitude         float64                `json:"longitude"`
	StartsAt          *time.Time             `json:"starts_at"`
	EndsAt            *time.Time             `json:"ends_at"`
	IsPublic          bool                   `json:"is_public"`
	ParticipantsCount int                    `json:"participants_count"`
	Category          GetCategoryResponse    `json:"category"`
	Creator           GetParticipantResponse `json:"creator"`
}

type CreateEventRequest struct {
	CategoryID  string     `json:"category_id" validate:"required,uuid"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Latitude    float64    `json:"latitude" validate:"required"`
	Longitude   float64    `json:"longitude" validate:"required"`
	Address     *string    `json:"address"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	IsPublic    bool       `json:"is_public"`
}

type GetCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetParticipantResponse struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarUrl string `json:"avatar_url"`
}
