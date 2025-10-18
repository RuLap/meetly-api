package repository

import (
	"context"
	"fmt"
	"github.com/RuLap/meetly-api/meetly/internal/app/category/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*models.Category, error)
}

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{pool}
}

func (r *categoryRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	query := `
		SELECT id, name
		FROM categories
	`

	var categories []*models.Category
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить категории: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("не удалось получить категории: %w", err)
		}

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("не удалось получить категории: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
	`

	var category models.Category
	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
	)

	if err != nil {
		return nil, fmt.Errorf("не удалось получить категорию: %w", err)
	}

	return &category, nil
}
