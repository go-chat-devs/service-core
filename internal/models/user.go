package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID int

	UserUID   uuid.UUID
	Username  *string
	AvatarUID *uuid.UUID
}

func (u *User) FromRow(row pgx.Row) error {
	return row.Scan(&u.ID, &u.UserUID, &u.Username, &u.AvatarUID)
}

func UserFactory() *User {
	return &User{}
}
