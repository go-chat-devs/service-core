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
	messageID int,
	text string,
	userUID uuid.UUID,
	changed_at *time.Time,
) error {
	const sql = "INSERT INTO text_messages(message_id, text, user_uid, changed_at) VALUES($1, $2, $3)"
	_, err := s.db.Exec(ctx, sql, messageID, text, userUID, changed_at)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(
	ctx context.Context,
	uid uuid.UUID,
) (*models.TextMessage, error) {
	const sql = "SELECT * FROM text_messages WHERE uid=$1"
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row(row, models.TextMessageFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectMany(
	ctx context.Context,
	uids []uuid.UUID,
) ([]*models.TextMessage, error) {
	const sql = "SELECT * FROM text_messages WHERE uid=ANY($1)"
	rows, err := s.db.Query(ctx, sql, uids)
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
	uid uuid.UUID,
	newText string,
	changed_at time.Time,
) error {
	const sql = "UPDATE text_messages SET text=$1, changed_at=$2 WHERE uid=$3"
	_, err := s.db.Exec(ctx, sql, newText, changed_at, uid)
	if err != nil {
		slog.Error(tag("update text error: %v", err))
	}
	return err
}
