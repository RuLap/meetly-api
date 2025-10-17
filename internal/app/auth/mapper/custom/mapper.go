package mapper

import (
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/dto"
	"github.com/RuLap/meetly-api/meetly/internal/app/auth/models"
)

func LoginRequestToUser(dto *dto.LoginRequest, hashedPassword string) *models.User {
	return &models.User{
		Email:    dto.Email,
		Password: &hashedPassword,
	}
}

func RegisterRequestToUser(dto *dto.RegisterRequest, hashedPassword string) *models.User {
	return &models.User{
		Email:          dto.Email,
		Password:       &hashedPassword,
		Provider:       models.LocalProvider,
		EmailConfirmed: false,
	}
}
