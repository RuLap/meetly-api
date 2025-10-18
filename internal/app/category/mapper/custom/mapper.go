package mapper

import (
	"github.com/RuLap/meetly-api/meetly/internal/app/category/dto"
	"github.com/RuLap/meetly-api/meetly/internal/app/category/models"
)

func CategoryToGetResponse(model *models.Category) *dto.GetCategoryResponse {
	return &dto.GetCategoryResponse{
		ID:   model.ID.String(),
		Name: model.Name,
	}
}
