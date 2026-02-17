package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GroupChat struct {
	UID       uuid.UUID
	Title     string
	BIO       *string
	AvatarUID *uuid.UUID
	CreatedAt time.Time
}

func (gc *GroupChat) FromRow(row pgx.Row) error {
	return row.Scan(
		&gc.UID,
		&gc.Title,
		&gc.BIO,
		&gc.AvatarUID,
		&gc.CreatedAt,
	)
}

func GroupChatFactory() *GroupChat {
	return &GroupChat{}
}
