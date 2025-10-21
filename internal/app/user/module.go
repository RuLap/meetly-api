package user

import (
	"log/slog"

	"github.com/RuLap/meetly-api/meetly/internal/pkg/providers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Repo    Repository
	Service Service
	Handler *Handler
}

func NewModule(log *slog.Logger, pool *pgxpool.Pool) *Module {
	repo := NewRepository(pool)
	service := NewService(log, repo)
	handler := NewHandler(service)

	return &Module{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}

func (m *Module) GetUserProvider() providers.UserProvider {
	return NewUserProvider(m.Service)
}
