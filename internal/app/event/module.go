package event

import (
	"log/slog"

	"github.com/RuLap/meetly-api/meetly/internal/pkg/providers"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	EventRepo       EventRepository
	CategoryRepo    CategoryRepository
	ParticipantRepo ParticipantRepository
	Service         Service
	Handler         Handler
}

func NewModule(log *slog.Logger, pool *pgxpool.Pool, userProvider providers.UserProvider) *Module {
	eventRepo := NewEventRepository(pool)
	categoryRepo := NewCategoryRepository(pool)
	participantRepo := NewParticipantRepository(pool)

	service := NewService(log, eventRepo, categoryRepo, participantRepo, userProvider)
	handler := NewHandler(service)

	return &Module{
		EventRepo:       eventRepo,
		CategoryRepo:    categoryRepo,
		ParticipantRepo: participantRepo,
		Service:         service,
		Handler:         *handler,
	}
}
