package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Chat struct {
	UID          uuid.UUID
	UserUID_low  *uuid.UUID
	UserUID_high *uuid.UUID
	CreatedAt    time.Time
}

func (c *Chat) FromRow(row pgx.Row) error {
	return row.Scan(&c.UID, &c.UserUID_low, &c.UserUID_high, &c.CreatedAt)
}

func ChatFactory() *Chat {
	return &Chat{}
}
