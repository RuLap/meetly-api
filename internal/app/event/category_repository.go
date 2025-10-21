package event

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]*Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
}

type categoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{pool}
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]*Category, error) {
	query := `
		SELECT id, name
		FROM categories
	`

	var categories []*Category
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить категории: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
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

func (r *categoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	query := `
		SELECT id, name
		FROM categories
		WHERE id = $1
	`

	var category Category
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
