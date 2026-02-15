package groupchats

import (
	"context"
	"log/slog"

	"github.com/go-chat-devs/service-core/internal/storage/db"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var tag = tagger.Tagger("storage-group-chats")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) Insert(ctx context.Context, title string, bio *string, avatarUID *uuid.UUID) error {
	const sql = `INSERT INTO group_chats(title, bio, avatar_uid) VALUES($1, $2, $)`
	_, err := s.db.Exec(ctx, sql, title, bio, avatarUID)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Delete(ctx context.Context, uid string) error {
	const sql = `DELETE FROM group_chats WHERE uid=$1`
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
