package groupchats

import (
	"context"
	"log/slog"
	"time"

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

func (s *Storage) Insert(ctx context.Context,
	title string,
	bio *string,
	avatarUID *uuid.UUID,
	created_at time.Time,
) error {
	const sql = `INSERT INTO core.group_chats(title, bio, avatar_uid, created_at) VALUES($1, $2, $3, $4)`
	_, err := s.db.Exec(ctx, sql, title, bio, avatarUID, created_at)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) UpdateTitle(ctx context.Context, uid uuid.UUID, title string) error {
	const sql = `UPDATE core.group_chats SET title=$1 WHERE uid=$2`
	_, err := s.db.Exec(ctx, sql, title, uid)
	if err != nil {
		slog.Error(tag("update title error: %v", err))
	}
	return err
}

func (s *Storage) UpdateBIO(ctx context.Context, uid uuid.UUID, bio *string) error {
	const sql = `UPDATE core.group_chats SET bio=$1 WHERE uid=$2`
	_, err := s.db.Exec(ctx, sql, bio, uid)
	if err != nil {
		slog.Error(tag("update bio error: %v", err))
	}
	return err
}

func (s *Storage) UpdateAvatar(ctx context.Context, uid uuid.UUID, avatarUID *uuid.UUID) error {
	const sql = `UPDATE core.group_chats SET avatar_uid=$1 WHERE uid=$2`
	_, err := s.db.Exec(ctx, sql, avatarUID, uid)
	if err != nil {
		slog.Error(tag("update avatar error: %v", err))
	}
	return err
}

func (s *Storage) Delete(ctx context.Context, uid string) error {
	const sql = `DELETE FROM core.group_chats WHERE uid=$1`
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
