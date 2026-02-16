package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TextMessage struct {
	UID       uuid.UUID
	MessageID int
	Text      string
	UserUID   *uuid.UUID
	Changed   *time.Time
}

func (m *TextMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.UID, &m.MessageID, &m.Text, &m.UserUID, &m.Changed)
}

func TextMessageFactory() *TextMessage {
	return &TextMessage{}
}
