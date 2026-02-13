package models

import "github.com/jackc/pgx/v5"

type User struct {
	ID        int64
	UserUID   string
	Username  string
	AvatarUID string
}

func (u *User) FromRow(row pgx.Row) error {
	return row.Scan(&u.ID, &u.UserUID, &u.Username, &u.AvatarUID)
}
