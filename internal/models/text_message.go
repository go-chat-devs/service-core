package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TextMessage struct {
	ID int

	MessageUID uuid.UUID
	Text       string
	From       *uuid.UUID
	Changed    *time.Time
}

func (m *TextMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.MessageUID, &m.Text, &m.From, &m.Changed)
}

func TextMessageFactory() *TextMessage {
	return &TextMessage{}
}
