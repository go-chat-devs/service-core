package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TextMessage struct {
	ID         int
	MessageUID uuid.UUID
	Content    string
	Changed    *time.Time
}

func (m *TextMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.MessageUID, &m.Content, &m.Changed)
}

func TextMessageFactory() *TextMessage {
	return &TextMessage{}
}
