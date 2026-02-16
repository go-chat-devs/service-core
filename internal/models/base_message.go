package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MesssageType string

const (
	Text  MesssageType = "text"
	Image MesssageType = "image"
)

type BaseMessage struct {
	ID int

	UID     uuid.UUID
	ChatUID uuid.UUID
	SentAt  time.Time
	Type    MesssageType
}

func (m *BaseMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.UID, &m.ChatUID, &m.Type, &m.SentAt)
}

func BaseMessageFactory() *BaseMessage {
	return &BaseMessage{}
}
