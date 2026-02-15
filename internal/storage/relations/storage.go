package relations

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-chat-devs/service-core/internal/models"
	"github.com/go-chat-devs/service-core/internal/scanner"
	"github.com/go-chat-devs/service-core/internal/storage/db"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var tag = tagger.Tagger("storage-relations")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) Insert(ctx context.Context, userUID, friendUID uuid.UUID, timestamp time.Time) error {
	const sql = `INSERT INTO relations(user_uid, friend_uid, timestamp) VALUES($1, $2, $3)`
	_, err := s.db.Exec(ctx, sql, userUID, friendUID, timestamp)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(ctx context.Context, userUID, friendUID uuid.UUID) (*models.Relation, error) {
	const sql = `SELECT * FROM relations WHERE user_uid=$1 AND friend_uid=$2`
	row := s.db.QueryRow(ctx, sql, userUID, friendUID)
	res, err := scanner.Row[*models.Relation](row)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectAll(ctx context.Context, userUID string) ([]*models.Relation, error) {
	const sql = `SELECT * FROM relations WHERE user_uid=$1`
	rows, err := s.db.Query(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows[*models.Relation](rows)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(ctx context.Context, userUID, friendUID string) error {
	const sql = `DELETE FROM relations WHERE user_uid=$1 AND friend_uid=$2`
	_, err := s.db.Exec(ctx, sql, userUID, friendUID)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}

func (s *Storage) DeleteAll(ctx context.Context, userUID string) error {
	const sql = `DELETE FROM relations WHERE user_uid=$1`
	_, err := s.db.Exec(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
