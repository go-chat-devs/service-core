package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ImageMessage struct {
	UID       uuid.UUID
	MessageID int
	FileUID   uuid.UUID
	From      *uuid.UUID
}

func (m *ImageMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.UID, &m.MessageID, &m.FileUID, &m.From)
}

func ImageMessageFactory() *ImageMessage {
	return &ImageMessage{}
}
