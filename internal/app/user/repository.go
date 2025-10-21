package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, req *User) (*User, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool}
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, first_name, last_name, birth_date, gender, avatar_url
		FROM users WHERE id = $1`

	var user User
	err := r.pool.QueryRow(ctx, query, id).
		Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.BirthDate,
			&user.Gender,
			&user.AvatarUrl,
		)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить пользователя")
	}

	return &user, nil
}

func (r *repository) Update(ctx context.Context, req *User) (*User, error) {
	query := `
		UPDATE users
		SET
			first_name = $2,
			last_name = $3,
			birth_date = $4,
			gender = $5,
			avatar_url = $6
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query, req.FirstName, req.LastName, req.BirthDate, req.Gender, req.AvatarUrl)
	if err != nil {
		return nil, fmt.Errorf("не удалось сохрнить пользователя")
	}

	return req, nil
}
