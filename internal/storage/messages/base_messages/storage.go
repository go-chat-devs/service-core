package basemessages

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

var tag = tagger.Tagger("storage-base-message")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) Insert(
	ctx context.Context,
	chatUID uuid.UUID,
	sent_at time.Time,
	messageType models.MesssageType,
) error {
	const sql = `INSERT INTO messages(chat_uid, sent_at, type) VALUES($1, $2, $3, $4)`
	_, err := s.db.Exec(ctx, sql, chatUID, sent_at, messageType)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(
	ctx context.Context,
	id int,
) (*models.BaseMessage, error) {
	const sql = "SELECT * FROM messages WHERE id=$1"
	row := s.db.QueryRow(ctx, sql, id)
	res, err := scanner.Row(row, models.BaseMessageFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectTypeAll(ctx context.Context, chatUID uuid.UUID, messageType models.MesssageType) ([]*models.BaseMessage, error) {
	const sql = "SELECT * FROM messages WHERE chat_uid=$1 AND type=$2"
	rows, err := s.db.Query(ctx, sql, chatUID, messageType)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows(rows, models.BaseMessageFactory)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectAll(ctx context.Context, chatUID uuid.UUID) ([]*models.BaseMessage, error) {
	const sql = "SELECT * FROM messages WHERE chat_uid=$1"
	rows, err := s.db.Query(ctx, sql, chatUID)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows(rows, models.BaseMessageFactory)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	id int,
) error {
	const sql = "DELETE * FROM messages WHERE id=$1"
	_, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
