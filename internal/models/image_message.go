package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ImageMessage struct {
	MessageUID uuid.UUID
	FileUID    uuid.UUID
}

func (m *ImageMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.MessageUID, &m.FileUID)
}

func ImageMessageFactory() *ImageMessage {
	return &ImageMessage{}
}
