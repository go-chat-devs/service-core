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
	senderUID *uuid.UUID,
	messageType models.MessageType,
	sentAt time.Time,
) (messageUID uuid.UUID, err error) {
	const sql = `INSERT INTO core.messages(chat_uid, sender_uid, type, sent_at) VALUES($1, $2, $3, $4) RETURNING uid`
	err = s.db.QueryRow(ctx, sql, chatUID, senderUID, messageType, sentAt).Scan(&messageUID)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return
}

func (s *Storage) Select(
	ctx context.Context,
	uid uuid.UUID,
) (*models.BaseMessage, error) {
	const sql = "SELECT * FROM core.messages WHERE uid=$1"
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row(row, models.BaseMessageFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectTypeAll(ctx context.Context, chatUID uuid.UUID, messageType models.MessageType) ([]*models.BaseMessage, error) {
	const sql = "SELECT * FROM core.messages WHERE chat_uid=$1 AND type=$2"
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
	const sql = "SELECT * FROM core.messages WHERE chat_uid=$1"
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
	uid uuid.UUID,
) error {
	const sql = "DELETE * FROM core.messages WHERE uid=$1"
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}

func (s *Storage) DeleteAll(
	ctx context.Context,
	chatUID uuid.UUID,
) error {
	const sql = "DELETE * FROM core.messages WHERE chat_uid=$1"
	_, err := s.db.Exec(ctx, sql, chatUID)
	if err != nil {
		slog.Error(tag("delete all error: %v", err))
	}
	return err
}
