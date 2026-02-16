package chats

import (
	"context"
	"log/slog"

	"github.com/go-chat-devs/service-core/internal/models"
	"github.com/go-chat-devs/service-core/internal/scanner"
	"github.com/go-chat-devs/service-core/internal/storage/db"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/go-chat-devs/service-core/internal/utils"
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
	const sql = `INSERT INTO core.chats(user_uid_low, uiser_uid_high) VALUES($1, $2)`
	userUID_low, userUID_high := utils.SortUUID(users[0], users[1])
	_, err := s.db.Exec(ctx, sql, userUID_low, userUID_high)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(ctx context.Context, uid uuid.UUID) (*models.Chat, error) {
	const sql = `SELECT * FROM core.chats WHERE uid = $1`
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row(row, models.ChatFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectUsers(ctx context.Context, users [2]uuid.UUID) (*models.Chat, error) {
	const sql = `SELECT * FROM core.chats WHERE user_uid_low=$1 AND user_uid_high=$2`
	userUID_low, userUID_high := utils.SortUUID(users[0], users[1])
	row := s.db.QueryRow(ctx, sql, userUID_low, userUID_high)
	res, err := scanner.Row(row, models.ChatFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectAll(ctx context.Context, userUID uuid.UUID) ([]*models.Chat, error) {
	const sql = `SELECT * FROM core.chats WHERE user_uid_low = $1 OR user_uid_high = $1`
	rows, err := s.db.Query(ctx, sql, userUID)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows(rows, models.ChatFactory)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(ctx context.Context, uid uuid.UUID) error {
	const sql = `DELETE FROM core.chats WHERE uid=$1`
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
