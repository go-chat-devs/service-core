package textmessages

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

var tag = tagger.Tagger("storage-text-message")

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
	messageUID uuid.UUID,
	content string,
	changedAt *time.Time,
) error {
	const sql = "INSERT INTO core.text_messages(message_uid, content, changed_at) VALUES($1, $2, $3)"
	_, err := s.db.Exec(ctx, sql, messageUID, content, changedAt)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(
	ctx context.Context,
	messageUID uuid.UUID,
) (*models.TextMessage, error) {
	const sql = "SELECT * FROM core.text_messages WHERE message_uid=$1"
	row := s.db.QueryRow(ctx, sql, messageUID)
	res, err := scanner.Row(row, models.TextMessageFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectMany(
	ctx context.Context,
	messageUIDs []uuid.UUID,
) ([]*models.TextMessage, error) {
	const sql = "SELECT * FROM core.text_messages WHERE message_uid=ANY($1)"
	rows, err := s.db.Query(ctx, sql, messageUIDs)
	if err != nil {
		slog.Error(tag("select many query error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows(rows, models.TextMessageFactory)
	if err != nil {
		slog.Error(tag("select many scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) UpdateText(
	ctx context.Context,
	messageUID uuid.UUID,
	newContent string,
	changed_at time.Time,
) error {
	const sql = "UPDATE core.text_messages SET content=$1, changed_at=$2 WHERE message_uid=$3"
	_, err := s.db.Exec(ctx, sql, newContent, changed_at, messageUID)
	if err != nil {
		slog.Error(tag("update text error: %v", err))
	}
	return err
}
