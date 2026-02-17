package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type MessageType string

const (
	MessageType_Text  MessageType = "text"
	MessageType_Image MessageType = "image"
)

type BaseMessage struct {
	UID       uuid.UUID
	ChatUID   uuid.UUID
	SenderUID *uuid.UUID
	Type      MessageType
	SentAt    time.Time
}

func (m *BaseMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.UID, &m.ChatUID, &m.SenderUID, &m.Type, &m.SentAt)
}

func BaseMessageFactory() *BaseMessage {
	return &BaseMessage{}
}
