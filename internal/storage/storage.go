package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chat-devs/service-core/internal/storage/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db DB

	Users *users.Storage
}

func New(ctx context.Context) (*Storage, error) {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		slog.Error(fmt.Sprintf("error parsing connection config: %v", err))
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating connection pool: %v", err))
		return nil, err
	}

	return &Storage{
		db: pool,

		Users: users.New(pool),
	}, nil
}
