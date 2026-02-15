package basemessage

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

var tag = tagger.Tagger("storage-baseMessage")

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
	timestamp time.Time,
	messageType models.MesssageType,
) error {
	const sql = `INSERT INTO messages(chat_uid, timestamp, type) VALUES ($1,$2,$3,$4)`
	_, err := s.db.Exec(ctx, sql, chatUID, timestamp, messageType)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) GetOne(
	ctx context.Context,
	uid uuid.UUID,
) (*models.BaseMessage, error) {
	const sql = "SELECT * FROM messages WHERE uid=$1"
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row[*models.BaseMessage](row)
	if err != nil {
		slog.Error(tag("get one error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetMany(ctx context.Context, chatUID uuid.UUID) ([]*models.BaseMessage, error) {
	const sql = "SELECT * FROM messages WHERE chat_uid=$1"
	rows, err := s.db.Query(ctx, sql, chatUID)
	if err != nil {
		slog.Error(tag("get many query error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows[*models.BaseMessage](rows)
	if err != nil {
		slog.Error(tag("get many scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	uid uuid.UUID,
) error {
	const sql = "DELETE * FROM messages WHERE uid=$1"
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
