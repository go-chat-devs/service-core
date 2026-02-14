package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chat-devs/service-core/internal/storage/db"
	basemessage "github.com/go-chat-devs/service-core/internal/storage/messages/base_message"
	imagemessage "github.com/go-chat-devs/service-core/internal/storage/messages/image_message"
	textmessage "github.com/go-chat-devs/service-core/internal/storage/messages/text_message"
	"github.com/go-chat-devs/service-core/internal/storage/relations"
	"github.com/go-chat-devs/service-core/internal/storage/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db db.DBTX

	Users *users.Storage
	Relations *relations.Storage
	BaseMessages *basemessage.Storage
	TextMessage *textmessage.Storage
	ImageMessage *imagemessage.Storage
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
		Relations: relations.New(pool),
		BaseMessages: basemessage.New(pool),
		TextMessage: textmessage.New(pool),
		ImageMessage: imagemessage.New(pool),
	}, nil
}
