package imagemessages

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

var tag = tagger.Tagger("storage-image-message")

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
	fileUID uuid.UUID,
	userUID uuid.UUID,
) error {
	const sql = "INSERT INTO image_messages(message_uid, file_uid, user_uid) VALUES($1, $2)"
	_, err := s.db.Exec(ctx, sql, messageUID, fileUID, userUID)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(
	ctx context.Context,
	messageUID uuid.UUID,
) (*models.ImageMessage, error) {
	const sql = "SELECT * FROM image_messages WHERE message_uid=$1"
	row := s.db.QueryRow(ctx, sql, messageUID)
	res, err := scanner.Row(row, models.ImageMessageFactory)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectMany(
	ctx context.Context,
	messageUIDs []uuid.UUID,
) ([]*models.ImageMessage, error) {
	const sql = "SELECT * FROM image_messages WHERE message_uid=ANY($1)"
	rows, err := s.db.Query(ctx, sql, messageUIDs)
	if err != nil {
		slog.Error(tag("select many query error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows(rows, models.ImageMessageFactory)
	if err != nil {
		slog.Error(tag("select many scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	uid uuid.UUID,
) error {
	const sql = "DELETE * FROM image_messages WHERE message_uid=$1"
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}
