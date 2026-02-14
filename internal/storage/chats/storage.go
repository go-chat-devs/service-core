package chats

import (
	"context"
	"log/slog"

	"github.com/go-chat-devs/service-core/internal/models"
	"github.com/go-chat-devs/service-core/internal/scanner"
	"github.com/go-chat-devs/service-core/internal/storage/db"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/jackc/pgx/v5"
)

var tag = tagger.Tagger("storage-chats")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) InsertChat(ctx context.Context, userUID1, userUID2 string) error {
	const sql = `INSERT INTO chats(user_uid_low, uiser_uid_high) VALUES($2, $3)`
	_, err := s.db.Exec(ctx, sql, userUID1, userUID2)
	if err != nil {
		slog.Error(tag("InsertRelation error: %v", err))
	}
	return err
}

func (s *Storage) DeleteChat(ctx context.Context, uid string) error {
	const sql = `DELETE FROM chats WHERE uid=$1`
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("DeleteChat error: %v", err))
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
