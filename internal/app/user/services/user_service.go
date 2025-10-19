package services

import (
	"context"
	"fmt"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/dto"
	mapper "github.com/RuLap/meetly-api/meetly/internal/app/user/mapper/custom"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/repository"
	"github.com/google/uuid"
	"log/slog"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*dto.GetUserResponse, error)
	UpdateUser(ctx context.Context, id string, req *dto.SaveUserRequest) (*dto.GetUserResponse, error)
}

type userService struct {
	log      *slog.Logger
	userRepo repository.UserRepository
}

func NewUserService(log *slog.Logger, userRepo repository.UserRepository) UserService {
	return &userService{log: log, userRepo: userRepo}
}

func (s *userService) GetByID(ctx context.Context, id string) (*dto.GetUserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get user by id", "id", id, "error", err)
		return nil, err
	}

	result := mapper.UserToGetResponse(user)

	return result, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req *dto.SaveUserRequest) (*dto.GetUserResponse, error) {
	user, err := mapper.SaveRequestToUser(req)
	if err != nil {
		s.log.Error("failed to map save request to user with id", "id", id, "error", err)
		return nil, fmt.Errorf("произошла ошибка")
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		s.log.Error("failed to parse user with user id", "id", id, "error", err)
	}
	user.ID = uuid

	entity, err := s.userRepo.Update(ctx, user)
	if err != nil {
		s.log.Error("failed to save user with id", "id", id, "error", err)
		return nil, err
	}

	result := mapper.UserToGetResponse(entity)

	return result, nil
}
