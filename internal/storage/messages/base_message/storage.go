package basemessage

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

func (s *Storage) InsertMessage(
	ctx context.Context,
	uid uuid.UUID,
	chatUID uuid.UUID,
	timestamp int64,
	senderUID uuid.UUID,
	typemessage string,
) error {
	sql := `INSERT INTO messages 
			(uid,chat_uid, timestamp, sender_uid,type) 
			VALUES ($1,$2,$3,$4,$5);`
	_, err := s.db.Exec(ctx, sql, uid, chatUID, timestamp, senderUID, typemessage)
	if err != nil {
		slog.Error(tag(""))
	}
	return err
}

func (s *Storage) GetMessage(
	ctx context.Context,
	uid uuid.UUID,
) (*models.BaseMessage, error) {
	sql := "SELECT * FROM messages WHERE uid=$1;"
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row[*models.BaseMessage](row)
	if err != nil {
		slog.Error(tag("Storage base messages error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetMessages(
	ctx context.Context,
	uids []uuid.UUID,
) ([]*models.BaseMessage, error) {
	sql := "SELECT * FROM messages WHERE uid=ANY($1);"
	rows, err := s.db.Query(ctx, sql, uids)
	if err != nil {
		slog.Error(tag("Storage base messages error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows[*models.BaseMessage](rows)
	if err != nil {
		slog.Error(tag("Storage base messages error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) DeleteMessage(
	ctx context.Context,
	uid uuid.UUID,
) error {
	sql := "DELETE * FROM messages WHERE uid=$1"
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("Storage base messages error: %v", err))
	}
	return err
}
