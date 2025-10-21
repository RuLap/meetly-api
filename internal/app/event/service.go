package event

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/RuLap/meetly-api/meetly/internal/pkg/providers"
	"github.com/google/uuid"
)

type Service interface {
	GetShortEvents(ctx context.Context) ([]GetShortEventResponse, error)
	GetEventWithDetails(ctx context.Context, id uuid.UUID) (*GetEventResponse, error)
	CreateEvent(ctx context.Context, req *CreateEventRequest, creatorID uuid.UUID) (*GetEventResponse, error)

	GetAllCategories(ctx context.Context) ([]*GetCategoryResponse, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*GetCategoryResponse, error)

	AddParticipant(ctx context.Context, userID, eventID uuid.UUID) (*GetParticipantResponse, error)
}

type service struct {
	log             *slog.Logger
	eventRepo       EventRepository
	categoryRepo    CategoryRepository
	participantRepo ParticipantRepository
	userProvider    providers.UserProvider
}

func NewService(
	log *slog.Logger,
	eventRepo EventRepository,
	categoryRepo CategoryRepository,
	participantRepo ParticipantRepository,
	userProvider providers.UserProvider,
) Service {
	return &service{
		log:             log,
		eventRepo:       eventRepo,
		categoryRepo:    categoryRepo,
		participantRepo: participantRepo,
		userProvider:    userProvider,
	}
}

func (s *service) GetShortEvents(ctx context.Context) ([]GetShortEventResponse, error) {
	events, err := s.eventRepo.GetAll(ctx)
	if err != nil {
		s.log.Error("failed to get all events", "error", err)
		return nil, err
	}

	result := make([]GetShortEventResponse, 0)
	for _, event := range events {
		creator, err := s.getParticipantByUserID(ctx, event.CreatorID)
		if err != nil {
			return nil, err
		}

		category, err := s.GetCategoryByID(ctx, event.CategoryID)
		if err != nil {
			return nil, err
		}

		participants, err := s.getParticipantsByEventID(ctx, event.ID)
		if err != nil {
			s.log.Error("failed to get event participants", "id", event.CategoryID, "error", err)
			return nil, err
		}
		participantsCount := len(participants)

		dto := EventToGetShortResponse(&event, participantsCount, creator, category)

		result = append(result, *dto)
	}

	return result, nil
}

func (s *service) GetEventWithDetails(ctx context.Context, id uuid.UUID) (*GetEventResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get event by id", "id", id, "error", err)
		return nil, err
	}

	creator, err := s.getParticipantByUserID(ctx, event.CreatorID)
	if err != nil {
		return nil, err
	}

	category, err := s.GetCategoryByID(ctx, event.CategoryID)
	if err != nil {
		return nil, err
	}

	participants, err := s.getParticipantsByEventID(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	result := EventToGetResponse(event, creator, category, participants)

	return result, nil
}

func (s *service) CreateEvent(ctx context.Context, req *CreateEventRequest, creatorID uuid.UUID) (*GetEventResponse, error) {
	model, err := CreateEventRequestToModel(req, creatorID)

	event, err := s.eventRepo.Create(ctx, model)
	if err != nil {
		s.log.Error("failed to create event", "error", err)
		return nil, err
	}

	creator, err := s.getParticipantByUserID(ctx, event.CreatorID)
	if err != nil {
		return nil, err
	}

	category, err := s.GetCategoryByID(ctx, event.CategoryID)
	if err != nil {
		return nil, err
	}

	participants, err := s.getParticipantsByEventID(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	result := EventToGetResponse(event, creator, category, participants)

	return result, nil
}

func (s *service) GetAllCategories(ctx context.Context) ([]*GetCategoryResponse, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		s.log.Error("failed to get categories by id", "error", err)
		return nil, err
	}

	var result []*GetCategoryResponse
	for _, category := range categories {
		categoryDTO := CategoryToGetResponse(category)
		result = append(result, categoryDTO)
	}

	return result, nil
}

func (s *service) GetCategoryByID(ctx context.Context, id uuid.UUID) (*GetCategoryResponse, error) {
	category, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get category by id", "category_id", id, "error", err)
		return nil, err
	}

	result := CategoryToGetResponse(category)

	return result, nil
}

func (s *service) AddParticipant(ctx context.Context, userID, eventID uuid.UUID) (*GetParticipantResponse, error) {
	participant := &Participant{
		UserID:  userID,
		EventID: eventID,
	}

	err := s.participantRepo.Create(ctx, participant)
	if err != nil {

	}

	result, err := s.getParticipantByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user participant by userID", "user_id", userID)
	}

	return result, nil
}

func (s *service) getParticipantsByEventID(ctx context.Context, eventID uuid.UUID) ([]GetParticipantResponse, error) {
	participants, err := s.participantRepo.GetAllByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	result := make([]GetParticipantResponse, 0, len(participants))
	for _, participant := range participants {
		dto, err := s.getParticipantByUserID(ctx, participant.UserID)
		if err != nil {
			s.log.Error("failed to get user participant by userID", "user_id", participant.UserID)
			return nil, err
		}

		result = append(result, *dto)
	}

	return result, nil
}

func (s *service) getParticipantByUserID(ctx context.Context, userID uuid.UUID) (*GetParticipantResponse, error) {
	user, err := s.userProvider.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Error("failed to fetch user from user provider", "user_id", userID)
		return nil, fmt.Errorf("не удалось получить участника: %w", err)
	}

	result := UserInfoToGetParticipantResponse(user)

	return result, nil
}
