package relations

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-chat-devs/service-core/internal/models"
	"github.com/go-chat-devs/service-core/internal/scanner"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tag = tagger.Tagger("storage-relations")

type Storage struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{db: db}
}

func (s *Storage) InsertRelation(ctx context.Context, userUID, friendUID string) error {
	const sql = `INSERT INTO relations(user_uid, friend_uid, timestamp) VALUES($1, $2, $3)`
	timestamp := time.Now().Unix()
	_, err := s.db.Exec(ctx, sql, userUID, friendUID, timestamp)
	if err != nil {
		slog.Error(tag("InsertRelation error: %v", err))
	}
	return err
}

func (s *Storage) DeleteRelation(ctx context.Context, userUID, friendUID string) error {
	const sql = `DELETE FROM relations WHERE user_uid=$1 AND friend_uid=$2`
	_, err := s.db.Exec(ctx, sql, userUID, friendUID)
	if err != nil {
		slog.Error(tag("DeleteRelation error: %v", err))
	}
	return err
}

func (s *Storage) GetRelation(ctx context.Context, userUID, friendUID string) (*models.Relation, bool) {
	const sql = `SELECT id, user_uid, friend_uid, timestamp FROM relations WHERE user_uid=$1 AND friend_uid=$2`
	row := s.db.QueryRow(ctx, sql, userUID, friendUID)
	res, err := scanner.Row[*models.Relation](row)
	if err != nil {
		return nil, false
	}
	return res, true
}

func (s *Storage) GetRelations(ctx context.Context, userUID string) ([]*models.Relation, error) {
	const sql = `SELECT id, user_uid, friend_uid, timestamp FROM relations WHERE user_uid=$1`
	rows, err := s.db.Query(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("GetRelations query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows[*models.Relation](rows)
	if err != nil {
		slog.Error(tag("GetRelations scan error: %v", err))
		return nil, err
	}
	return res, nil
}
