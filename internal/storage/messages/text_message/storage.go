package textmessage

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

var tag = tagger.Tagger("storage-textMessage")

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
	text string,
	changed int64,
) error  {
	sql := "INSERT INTO text_messages (message_uid,text,changed) VALUES ($1,$2,$3);"
	_, err := s.db.Exec(ctx, sql, uid,text,changed)
	if err != nil{
		slog.Error(tag("Storage text messages error: %v",err))
	}
	return err
}





func (s *Storage) GetMessage(
	ctx context.Context,
	uid uuid.UUID,
) (*models.TextMessage,error){
	sql := "SELECT * FROM text_messages WHERE message_uid=$1;"
	row := s.db.QueryRow(ctx, sql, uid)
	res, err := scanner.Row[*models.TextMessage](row)
	if err != nil{
		slog.Error(tag("Storage text messages error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetMessages(
	ctx context.Context,
	uids []uuid.UUID,
) ([]*models.TextMessage, error) {
	sql := "SELECT * FROM text_messages WHERE message_uid=ANY($1);"
	rows, err := s.db.Query(ctx, sql, uids)
	if err != nil {
		slog.Error(tag("Storage text messages error: %v", err))
		return nil, err
	}
	res, err := scanner.Rows[*models.TextMessage](rows)
	if err != nil {
		slog.Error(tag("Storage text messages error: %v", err))
		return nil, err
	}
	return res, nil
}


func (s *Storage) DeleteMessage(
	ctx context.Context,
	uid uuid.UUID,
) error {
	sql := "DELETE * FROM text_messages WHERE message_uid=$1"
	_, err := s.db.Exec(ctx, sql, uid)
	if err != nil {
		slog.Error(tag("Storage text messages error: %v", err))
	}
	return err
}



func (s *Storage) ChangeText(
	ctx context.Context,
	uid uuid.UUID,
	newText string,
) error{
	timestamp := time.Now().Unix()
	sql := "UPDATE text_messages SET text=$1, changed=$2 WHERE message_uid=$3;"
	_, err := s.db.Exec(ctx,sql,newText,timestamp,uid)
	if err != nil{
		slog.Error(tag("Storage text messages error: %v", err))
	}
	return err
}