package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MesssageType string

const (
	Text  MesssageType = "text"
	Image MesssageType = "image"
)

type BaseMessage struct {
	ID          int
	UID         uuid.UUID
	ChatUID     uuid.UUID
	Timestamp   int64
	SenderUID   uuid.UUID
	TypeMessage MesssageType
}

func (m *BaseMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.UID, &m.ChatUID, &m.Timestamp, &m.SenderUID, &m.TypeMessage)
}
