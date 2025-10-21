package user

import (
	"context"
	"github.com/google/uuid"

	"github.com/RuLap/meetly-api/meetly/internal/pkg/providers"
)

type userProvider struct {
	service Service
}

func NewUserProvider(service Service) providers.UserProvider {
	return &userProvider{service: service}
}

func (p *userProvider) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) (map[string]providers.UserInfo, error) {
	users, err := p.service.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[string]providers.UserInfo)
	for _, user := range users {
		result[user.ID] = providers.UserInfo{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			BirthDate: user.BirthDate,
			Gender:    user.Gender,
			AvatarUrl: user.AvatarUrl,
		}
	}

	return result, nil
}

func (p *userProvider) GetUserByID(ctx context.Context, userID uuid.UUID) (*providers.UserInfo, error) {
	user, err := p.service.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := providers.UserInfo{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		BirthDate: user.BirthDate,
		Gender:    user.Gender,
		AvatarUrl: user.AvatarUrl,
	}

	return &result, err
}
