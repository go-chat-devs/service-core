package chats

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

func (s *Storage) Insert(ctx context.Context, users [2]uuid.UUID) error {
	const sql = `INSERT INTO chats(user_uid_low, uiser_uid_high) VALUES($2, $3)`
	var userUID_low, userUID_high uuid.UUID
	if users[0].String() < users[1].String() {
		userUID_low, userUID_high = users[0], users[1]
	} else {
		userUID_low, userUID_high = users[1], users[0]
	}
	_, err := s.db.Exec(ctx, sql, userUID_low, userUID_high)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(ctx context.Context, uid uuid.UUID) (*models.Chat, error) {
	const sql = `SELECT * FROM chats WHERE uid = $1`
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row[*models.Chat](row)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectAll(ctx context.Context, userUID uuid.UUID) ([]*models.Chat, error) {
	const sql = `SELECT * FROM chats WHERE user_uid_low = $1 OR user_uid_high = $1`
	rows, err := s.db.Query(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows[*models.Chat](rows)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(ctx context.Context, uid uuid.UUID) error {
	const sql = `DELETE FROM chats WHERE uid=$1`
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
