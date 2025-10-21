package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*GetUserResponse, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) (map[string]GetUserResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *SaveUserRequest) (*GetUserResponse, error)
}

type service struct {
	log  *slog.Logger
	repo Repository
}

func NewService(log *slog.Logger, repo Repository) Service {
	return &service{log: log, repo: repo}
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*GetUserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get user by id", "id", id, "error", err)
		return nil, err
	}

	result := UserToGetResponse(user)

	return result, nil
}

func (s *service) GetByIDs(ctx context.Context, ids []uuid.UUID) (map[string]GetUserResponse, error) {
	return nil, nil
}

func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req *SaveUserRequest) (*GetUserResponse, error) {
	user, err := SaveRequestToUser(req, id)
	if err != nil {
		s.log.Error("failed to map save request to user with id", "id", id, "error", err)
		return nil, fmt.Errorf("произошла ошибка")
	}

	entity, err := s.repo.Update(ctx, user)
	if err != nil {
		s.log.Error("failed to save user with id", "id", id, "error", err)
		return nil, err
	}

	result := UserToGetResponse(entity)

	return result, nil
}
