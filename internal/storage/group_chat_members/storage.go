package groupchatmembers

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

var tag = tagger.Tagger("storage-group-chat-members")

type Storage struct {
	db db.DBTX
}

func New(db db.DBTX) *Storage {
	return &Storage{db: db}
}

func (s *Storage) WithTX(tx pgx.Tx) *Storage {
	return &Storage{db: tx}
}

func (s *Storage) Insert(ctx context.Context, chatUID, userUID uuid.UUID, role models.MemberRole) error {
	const sql = `INSERT INTO group_chat_members(chat_uid, user_uid, role) VALUES($1, $2, $)`
	_, err := s.db.Exec(ctx, sql, chatUID, userUID, role)
	if err != nil {
		slog.Error(tag("insert error: %v", err))
	}
	return err
}

func (s *Storage) Select(ctx context.Context, chatUID, userUID uuid.UUID) (*models.GroupChatMember, error) {
	const sql = "SELECT * FROM group_chat_members WHERE chat_uid=$1 AND user_uid=$2"
	row := s.db.QueryRow(ctx, sql, chatUID, userUID)
	res, err := scanner.Row[*models.GroupChatMember](row)
	if err != nil {
		slog.Error(tag("select error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) SelectAll(ctx context.Context, chatUID uuid.UUID) ([]*models.GroupChatMember, error) {
	const sql = "SELECT * FROM group_chat_members WHERE chat_uid=$1"
	rows, err := s.db.Query(ctx, sql, chatUID)
	if err != nil {
		slog.Error(tag("select all query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	res, err := scanner.Rows[*models.GroupChatMember](rows)
	if err != nil {
		slog.Error(tag("select all scan error: %v", err))
		return nil, err
	}
	return res, nil
}

func (s *Storage) UpdateRole(ctx context.Context, id int, newRole models.MemberRole) error {
	const sql = `UPDATE group_chat_members SET role=$1 WHERE id=$2`
	_, err := s.db.Exec(ctx, sql, newRole, id)
	if err != nil {
		slog.Error(tag("update role error: %v", err))
	}
	return err
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	const sql = `DELETE FROM group_chat_members WHERE id=$1`
	_, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		slog.Error(tag("delete error: %v", err))
	}
	return err
}

func (s *Storage) DeleteAll(ctx context.Context, chatUID uuid.UUID) error {
	const sql = `DELETE FROM group_chat_members WHERE chat_uid=$1`
	_, err := s.db.Exec(ctx, sql, chatUID)
	if err != nil {
		slog.Error(tag("delete all error: %v", err))
	}
	return err
}
