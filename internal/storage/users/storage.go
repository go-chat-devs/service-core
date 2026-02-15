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

func (s *Storage) InsertNewUser(ctx context.Context, userUid, avatarUid uuid.UUID, username string) error {

	sql := "INSERT INTO users (user_uid, username, avatar_uid) VALUES ($1,$2,$3);"
	_, err := s.db.Exec(ctx, sql, userUid, username, avatarUid)
	if err != nil {
		slog.Error(tag("InsertUser error: %v", err))

	}
	return err
}

func (s *Storage) DeleteUser(ctx context.Context, userUid uuid.UUID) error {
	sql := "DELETE FROM users WHERE user_uid=$1;"

	_, err := s.db.Exec(ctx, sql, userUid)
	if err != nil {
		slog.Error(tag("Delete User error: %v", err))
	}

	return err
}

func (s *Storage) SelectUserByUid(ctx context.Context, userUid uuid.UUID) (*models.User, error) {
	sql := "SELECT * FROM users WHERE user_uid = $1;"
	row := s.db.QueryRow(ctx, sql, userUid)
	user, err := scanner.Row[*models.User](row)
	if err != nil {
		slog.Error("failed to scan row to user")
		return nil, err
	}
	return user, nil
}

func (s *Storage) SelectUsersByUid(ctx context.Context, userUids []uuid.UUID) ([]*models.User, error) {
	sql := "SELECT * FROM users WHERE user_uid = ANY($1);"

	rows, err := s.db.Query(ctx, sql, userUids)
	if err != nil {
		slog.Error(tag("Select Users error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows[*models.User](rows)
	if err != nil {
		slog.Error(tag("Select Users error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) ChangeUsername(ctx context.Context, userUID uuid.UUID, newUsername string) error {
	sql := "UPDATE users SET username = $1 WHERE user_uid = $2;"
	_, err := s.db.Exec(ctx, sql, newUsername, userUID)
	if err != nil {
		slog.Error(tag("Change Username error: %v", err))
	}
	return err
}

func (s *Storage) DeleteUsername(ctx context.Context, userUID uuid.UUID) error {
	sql := "UPDATE users SET username = NULL WHERE user_uid = $1;"
	_, err := s.db.Exec(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("Delete Username error: %v", err))
	}
	return err
}

func (s *Storage) ChangeAvatar(ctx context.Context, userUID, avatarUID uuid.UUID) error {
	sql := "UPDATE users SET avatar_uid = $1 WHERE user_uid = $2;"
	_, err := s.db.Exec(ctx, sql, avatarUID, userUID)
	if err != nil {
		slog.Error(tag("Change Avatar error: %v", err))
	}

	return err
}

func (s *Storage) DeleteAvatar(ctx context.Context, userUID uuid.UUID) error {
	sql := "UPDATE users SET avatar_uid = NULL WHERE user_uid = $1;"
	_, err := s.db.Exec(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("Delete Avatar error: %v", err))
	}

	return err
}
