package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Relation struct {
	ID int

	UserUID   uuid.UUID
	FriendUID uuid.UUID
	Timestamp time.Time
}

func (r *Relation) FromRow(row pgx.Row) error {
	return row.Scan(&r.ID, &r.UserUID, &r.FriendUID, &r.Timestamp)
}
