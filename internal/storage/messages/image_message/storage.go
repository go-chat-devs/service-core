package imagemessage

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

var tag = tagger.Tagger("storage-imageMessage")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) InsertMessage(
	ctx context.Context,
	uid uuid.UUID,
	fileUID uuid.UUID,
) error {
	sql := "INSERT INTO image_messages (message_uid,file_uid) VALUES ($1,$2);"
	_, err := s.db.Exec(ctx, sql, uid, fileUID)
	if err != nil {
		slog.Error(tag("Storage image messages error: %v", err))
	}
	return err
}

func (s *Storage) GetMessage(
	ctx context.Context,
	uid uuid.UUID,
) (*models.ImageMessage, error) {
	sql := "SELECT * FROM image_messages WHERE message_uid=$1;"
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row[*models.ImageMessage](row)
	if err != nil {
		slog.Error(tag("Storage image messages error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetMessages(
	ctx context.Context,
	uids []uuid.UUID,
) ([]*models.ImageMessage, error) {
	sql := "SELECT * FROM image_messages WHERE message_uid=ANY($1);"
	rows, err := s.db.Query(ctx, sql, uids)
	if err != nil {
		slog.Error(tag("Storage image messages error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows[*models.ImageMessage](rows)
	if err != nil {
		slog.Error(tag("Storage image messages error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) DeleteMessage(
	ctx context.Context,
	uid uuid.UUID,
) error {
	sql := "DELETE * FROM image_messages WHERE message_uid=$1"
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("Storage image messages error: %v", err))
	}
	return err
}
