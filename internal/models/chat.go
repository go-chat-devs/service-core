package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Chat struct {
	ID int

	UID          uuid.UUID
	UserUID_low  *uuid.UUID
	UserUID_high *uuid.UUID
}

func (c *Chat) FromRow(row pgx.Row) error {
	return row.Scan(&c.ID, &c.UID, &c.UserUID_low, &c.UserUID_high)
}

func ChatFactory() *Chat {
	return &Chat{}
}
