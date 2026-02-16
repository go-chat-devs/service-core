package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ImageMessage struct {
	UID       uuid.UUID
	MessageID int
	FileUID   uuid.UUID
	UserUID   *uuid.UUID
}

func (m *ImageMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.UID, &m.MessageID, &m.FileUID, &m.UserUID)
}

func ImageMessageFactory() *ImageMessage {
	return &ImageMessage{}
}
