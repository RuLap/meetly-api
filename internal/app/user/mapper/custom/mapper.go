package mapper

import (
	"fmt"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/dto"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/models"
	"time"
)

func SaveRequestToUser(dto *dto.SaveUserRequest) (*models.User, error) {
	model := models.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		AvatarUrl: dto.AvatarUrl,
	}

	switch dto.Gender {
	case string(models.GenderMale):
		model.Gender = models.GenderMale
	case string(models.GenderFemale):
		model.Gender = models.GenderFemale
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

func UserToGetResponse(model *models.User) *dto.GetUserResponse {
	return &dto.GetUserResponse{
		FirstName: model.FirstName,
		LastName:  model.LastName,
		BirthDate: model.BirthDate.Format(time.DateOnly),
		Gender:    string(model.Gender),
		AvatarUrl: model.AvatarUrl,
	}
}
