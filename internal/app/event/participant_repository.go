package event

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ParticipantRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Participant, error)
	GetAllByEventID(ctx context.Context, eventID uuid.UUID) ([]Participant, error)
	Create(ctx context.Context, model *Participant) error
}

type participantRepository struct {
	pool *pgxpool.Pool
}

func NewParticipantRepository(pool *pgxpool.Pool) ParticipantRepository {
	return &participantRepository{pool}
}

func (r *participantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Participant, error) {
	query := `
		SELECT user_id, event_id
		FROM participants
		WHERE user_id = $1
	`

	var participant Participant
	err := r.pool.QueryRow(ctx, query, userID).Scan(&participant.UserID, &participant.EventID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить участника события по UserID: %w", err)
	}

	return &participant, err
}

func (r *participantRepository) GetAllByEventID(ctx context.Context, eventID uuid.UUID) ([]Participant, error) {
	query := `
		SELECT user_id, event_id
		FROM participants
		WHERE event_id = $1
	`

	rows, err := r.pool.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить участников события: %w", err)
	}

	result := make([]Participant, 0)
	for rows.Next() {
		var participant Participant
		err := rows.Scan(&participant.UserID, &participant.EventID)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить участника события по UserID: %w", err)
		}
	}

	return result, nil
}

func (r *participantRepository) Create(ctx context.Context, model *Participant) error {
	query := `
		INSERT INTO participants
		values ($1, $2)
	`
	_, err := r.pool.Exec(ctx, query, model.UserID, model.EventID)
	if err != nil {
		if isUniqueConstraintError(err) {
			return fmt.Errorf("пользователь уже учавствует в событии: %w", err)
		}
		return fmt.Errorf("не удалось добавить участника к событию: %w", err)
	}

	return nil
}
