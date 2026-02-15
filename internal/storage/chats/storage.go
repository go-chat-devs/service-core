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

func (s *Storage) GetChat(ctx context.Context, uid string) (*models.Chat, bool) {
	const sql = `SELECT * FROM chats WHERE uid = $1`
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row[*models.Chat](row)
	if err != nil {
		return nil, false
	}
	return res, true
}

func (s *Storage) GetChats(ctx context.Context, userUID string) ([]*models.Chat, error) {
	const sql = `SELECT * FROM chats WHERE user_uid_low = $1 OR user_uid_high = $1`
	rows, err := s.db.Query(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("GetChats query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows[*models.Chat](rows)
	if err != nil {
		slog.Error(tag("GetChats scan error: %v", err))
		return nil, err
	}
	return res, nil
}
