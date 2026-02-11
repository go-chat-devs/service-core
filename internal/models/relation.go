package models

import "github.com/jackc/pgx/v5"

type Relation struct {
	ID int

	UserUID   string `json:"user_uid"`
	FriendUID string `json:"friend_uid"`
	Timestamp int64  `json:"timestamp"`
}

func (r *Relation) FromRow(row pgx.Row) error {
	return row.Scan(&r.ID, &r.UserUID, &r.FriendUID, &r.Timestamp)
}
