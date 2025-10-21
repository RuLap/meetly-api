package providers

import (
	"context"
	"github.com/google/uuid"
)

type UserProvider interface {
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (map[string]UserInfo, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
}

type UserInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	AvatarUrl string `json:"avatar_url"`
}
