package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ImageMessage struct {
	ID int

	MessageUID uuid.UUID
	FileUID    uuid.UUID
	From       *uuid.UUID
}

func (m *ImageMessage) FromRow(row pgx.Row) error {
	return row.Scan(&m.ID, &m.MessageUID, &m.FileUID, &m.From)
}

func ImageMessageFactory() *ImageMessage {
	return &ImageMessage{}
}
