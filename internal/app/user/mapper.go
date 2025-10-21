package user

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func SaveRequestToUser(dto *SaveUserRequest, id uuid.UUID) (*User, error) {
	model := User{
		ID:        id,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		AvatarUrl: dto.AvatarUrl,
	}

	switch dto.Gender {
	case string(GenderMale):
		model.Gender = GenderMale
	case string(GenderFemale):
		model.Gender = GenderFemale
	default:
		return nil, fmt.Errorf("неверное значение поля пол: %s", dto.Gender)
	}

	birthDate, err := time.Parse(time.DateOnly, dto.BirthDate)
	if err != nil {
		return nil, err
	}
	model.BirthDate = birthDate

	return &model, nil
}

func UserToGetResponse(model *User) *GetUserResponse {
	return &GetUserResponse{
		FirstName: model.FirstName,
		LastName:  model.LastName,
		BirthDate: model.BirthDate.Format(time.DateOnly),
		Gender:    string(model.Gender),
		AvatarUrl: model.AvatarUrl,
	}
}
