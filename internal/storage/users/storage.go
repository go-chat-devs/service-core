package users

import (
	"context"
	"log/slog"

	"github.com/go-chat-devs/service-core/internal/models"
	"github.com/go-chat-devs/service-core/internal/scanner"
	"github.com/go-chat-devs/service-core/internal/storage/db"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var tag = tagger.Tagger("storage-users")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) Insert(ctx context.Context, userUID uuid.UUID, username *string, avatarUid *uuid.UUID) error {
	const sql = "INSERT INTO core.users(uid, username, avatar_uid) VALUES($1, $2, $3)"
	_, err := s.db.Exec(ctx, sql, userUID, username, avatarUid)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(ctx context.Context, userUID uuid.UUID) (*models.User, error) {
	const sql = "SELECT * FROM core.users WHERE uid = $1"
	row := s.db.QueryRow(ctx, sql, userUID)
	user, err := scanner.Row(row, models.UserFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return user, nil
}

func (s *Storage) SelectUsername(ctx context.Context, username string) (*models.User, error) {
	const sql = "SELECT * FROM core.users WHERE username = $1"
	row := s.db.QueryRow(ctx, sql, username)
	user, err := scanner.Row(row, models.UserFactory)
	if err != nil {
		slog.Error(tag("select username error: %v", err))
		return nil, err
	}
	return user, nil
}

func (s *Storage) UpdateUsername(ctx context.Context, userUID uuid.UUID, newUsername *string) error {
	const sql = "UPDATE core.users SET username = $1 WHERE uid = $2"
	_, err := s.db.Exec(ctx, sql, newUsername, userUID)
	if err != nil {
		slog.Error(tag("update username error: %v", err))
	}
	return err
}

func (s *Storage) UpdateAvatar(ctx context.Context, userUID uuid.UUID, avatarUID *uuid.UUID) error {
	const sql = "UPDATE core.users SET avatar_uid = $1 WHERE uid = $2"
	_, err := s.db.Exec(ctx, sql, avatarUID, userUID)
	if err != nil {
		slog.Error(tag("update avatar error: %v", err))
	}
	return err
}

func (s *Storage) DeleteUser(ctx context.Context, userUID uuid.UUID) error {
	const sql = "DELETE FROM core.users WHERE uid=$1"

	_, err := s.db.Exec(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("Delete User error: %v", err))
	}

	return err
}
