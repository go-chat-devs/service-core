package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MemberRole string

const (
	MemberRole_Admin   MemberRole = "admin"
	MemberRole_Default MemberRole = "default"
)

type GroupChatMember struct {
	ID int

	ChatUID uuid.UUID
	UserUID uuid.UUID
	Role    MemberRole
}

func (m *GroupChatMember) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.ChatUID, &m.UserUID, &m.Role)
}
