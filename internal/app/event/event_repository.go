package event

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	GetAll(ctx context.Context) ([]Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Event, error)
	Create(ctx context.Context, model *Event) (*Event, error)
}

type eventRepository struct {
	pool *pgxpool.Pool
}

func NewEventRepository(pool *pgxpool.Pool) EventRepository {
	return &eventRepository{pool}
}

func (r *eventRepository) GetAll(ctx context.Context) ([]Event, error) {
	query := `
		SELECT id, creator_id, category_id, title, description,	latitude,
			longitude, address, starts_at, ends_at, is_public, created_at
		FROM events
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить события: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0)
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.CreatorID,
			&event.CategoryID,
			&event.Title,
			&event.Description,
			&event.Latitude,
			&event.Longitude,
			&event.Address,
			&event.StartsAt,
			&event.EndsAt,
			&event.IsPublic,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить событие: %w", err)
		}

	}

	return events, nil
}

func (r *eventRepository) GetByID(ctx context.Context, id uuid.UUID) (*Event, error) {
	query := `
		SELECT id, creator_id, category_id, title, description,	latitude,
			longitude, address, starts_at, ends_at, is_public, created_at
		FROM events
		WHERE id = $1
	`

	var event Event
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.CreatorID,
		&event.CategoryID,
		&event.Title,
		&event.Description,
		&event.Latitude,
		&event.Longitude,
		&event.Address,
		&event.StartsAt,
		&event.EndsAt,
		&event.IsPublic,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить событие: %w", err)
	}

	return &event, nil
}

func (r *eventRepository) Create(ctx context.Context, model *Event) (*Event, error) {
	query := `
		INSERT INTO events
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		model.CreatorID,
		model.CategoryID,
		model.Title,
		model.Description,
		model.Latitude,
		model.Longitude,
		model.Address,
		model.StartsAt,
		model.EndsAt,
		model.IsPublic,
	).Scan(
		&model.ID,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil, fmt.Errorf("событие уже существует: %w", err)
		}
		return nil, fmt.Errorf("не удалось создать событие: %w", err)
	}

	return model, nil
}

func isUniqueConstraintError(err error) bool {
	if pgErr, ok := err.(*pgx.PgError); ok {
		return pgErr.Code == "23505"
	}
	return false
}
