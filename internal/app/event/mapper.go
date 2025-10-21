package event

import (
	"github.com/RuLap/meetly-api/meetly/internal/pkg/providers"
	"github.com/google/uuid"
)

func EventToGetShortResponse(model *Event, participantsCount int, creator *GetParticipantResponse, category *GetCategoryResponse) *GetShortEventResponse {
	return &GetShortEventResponse{
		ID:                model.ID.String(),
		Title:             model.Title,
		Latitude:          model.Latitude,
		Longitude:         model.Longitude,
		StartsAt:          model.StartsAt,
		EndsAt:            model.EndsAt,
		IsPublic:          model.IsPublic,
		ParticipantsCount: participantsCount,
		Category:          *category,
		Creator:           *creator,
	}
}

func EventToGetResponse(model *Event, creator *GetParticipantResponse, category *GetCategoryResponse, participants []GetParticipantResponse) *GetEventResponse {
	return &GetEventResponse{
		ID:           model.ID.String(),
		Title:        model.Title,
		Description:  model.Description,
		Latitude:     model.Latitude,
		Longitude:    model.Longitude,
		Address:      model.Address,
		StartsAt:     model.StartsAt,
		EndsAt:       model.EndsAt,
		IsPublic:     model.IsPublic,
		CreatedAt:    model.CreatedAt,
		Creator:      *creator,
		Category:     *category,
		Participants: participants,
	}
}

func CreateEventRequestToModel(dto *CreateEventRequest, creatorID uuid.UUID) (*Event, error) {
	model := &Event{
		CreatorID:   creatorID,
		Title:       dto.Title,
		Description: dto.Description,
		Latitude:    dto.Latitude,
		Longitude:   dto.Longitude,
		Address:     dto.Address,
		StartsAt:    dto.StartsAt,
		EndsAt:      dto.EndsAt,
		IsPublic:    dto.IsPublic,
	}

	categoryID, err := uuid.Parse(dto.CategoryID)
	if err != nil {
		return nil, err
	}

	model.CategoryID = categoryID

	return model, nil
}

func CategoryToGetResponse(model *Category) *GetCategoryResponse {
	return &GetCategoryResponse{
		ID:   model.ID.String(),
		Name: model.Name,
	}
}

func UserInfoToGetParticipantResponse(user *providers.UserInfo) *GetParticipantResponse {
	return &GetParticipantResponse{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarUrl: user.AvatarUrl,
	}
}
