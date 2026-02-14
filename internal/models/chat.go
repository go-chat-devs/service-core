package models

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type Chat struct {
	ID int

	UID          uuid.UUID `json:"uid"`
	UserUID_low  uuid.UUID `json:"user_uid_low"`
	UserUID_high uuid.UUID `json:"user_uid_high"`
}

func (c *Chat) FromRow(row pgx.Row) error {
	return row.Scan(&c.ID, &c.UID, &c.UserUID_low, &c.UserUID_high)
}
