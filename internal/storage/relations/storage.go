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

func (s *Storage) Insert(ctx context.Context, userUID, friendUID uuid.UUID, added_at time.Time) error {
	const sql = `INSERT INTO core.relations(user_uid, friend_uid, added_at) VALUES($1, $2, $3)`
	_, err := s.db.Exec(ctx, sql, userUID, friendUID, added_at)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(ctx context.Context, userUID, friendUID uuid.UUID) (*models.Relation, error) {
	const sql = `SELECT * FROM core.relations WHERE user_uid=$1 AND friend_uid=$2`
	row := s.db.QueryRow(ctx, sql, userUID, friendUID)
	res, err := scanner.Row(row, models.RelationFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectAll(ctx context.Context, userUID uuid.UUID) ([]*models.Relation, error) {
	const sql = `SELECT * FROM core.relations WHERE user_uid=$1`
	rows, err := s.db.Query(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows(rows, models.RelationFactory)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(ctx context.Context, userUID, friendUID uuid.UUID) error {
	const sql = `DELETE FROM core.relations WHERE user_uid=$1 AND friend_uid=$2`
	_, err := s.db.Exec(ctx, sql, userUID, friendUID)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}

func (s *Storage) DeleteAll(ctx context.Context, userUID uuid.UUID) error {
	const sql = `DELETE FROM core.relations WHERE user_uid=$1`
	_, err := s.db.Exec(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
