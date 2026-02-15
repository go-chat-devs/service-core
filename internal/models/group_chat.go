package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GroupChat struct {
	ID int

	UID       uuid.UUID
	Title     string
	BIO       *string
	AvatarUID *uuid.UUID
	Timestamp time.Time
}

func (gc *GroupChat) FromRow(row pgx.Row) error {
	return row.Scan(
		&gc.ID,
		&gc.UID,
		&gc.Title,
		&gc.BIO,
		&gc.AvatarUID,
		&gc.Timestamp,
	)
}
