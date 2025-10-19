package repository

import (
	"context"
	"fmt"
	"github.com/RuLap/meetly-api/meetly/internal/app/user/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, req *models.User) (*models.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, first_name, last_name, birth_date, gender, avatar_url
		FROM users WHERE id = $1`

	var user models.User
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

func (r *userRepository) Update(ctx context.Context, req *models.User) (*models.User, error) {
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
