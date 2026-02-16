package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MessageType string

const (
	Text  MessageType = "text"
	Image MessageType = "image"
)

type BaseMessage struct {
	ID      int
	ChatUID uuid.UUID
	SentAt  time.Time
	Type    MessageType
}

func (m *BaseMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.ChatUID, &m.Type, &m.SentAt)
}

func BaseMessageFactory() *BaseMessage {
	return &BaseMessage{}
}
