package services

import (
	"context"
	"github.com/RuLap/meetly-api/meetly/internal/app/category/dto"
	mapper "github.com/RuLap/meetly-api/meetly/internal/app/category/mapper/custom"
	"github.com/RuLap/meetly-api/meetly/internal/app/category/repository"
	"log/slog"
)

type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]*dto.GetCategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*dto.GetCategoryResponse, error)
}

type categoryService struct {
	log          *slog.Logger
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(log *slog.Logger, categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{log: log, categoryRepo: categoryRepo}
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]*dto.GetCategoryResponse, error) {
	categories, err := s.categoryRepo.GetAllCategories(ctx)
	if err != nil {
		s.log.Warn("failed to get categories by id", "error", err)
		return nil, err
	}

	var result []*dto.GetCategoryResponse
	for _, category := range categories {
		categoryDTO := mapper.CategoryToGetResponse(category)
		result = append(result, categoryDTO)
	}

	return result, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*dto.GetCategoryResponse, error) {
	category, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		s.log.Warn("failed to get category by id", "category_id", id, "error", err)
		return nil, err
	}

	result := mapper.CategoryToGetResponse(category)

	return result, nil
}
