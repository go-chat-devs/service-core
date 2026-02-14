package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TextMessage struct {
	ID         int
	MessageUID uuid.UUID
	Text       string
	Changed int64
}

func (m *TextMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.MessageUID, &m.Text, &m.Changed)
}
