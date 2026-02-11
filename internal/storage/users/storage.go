package users

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// TODO: Croshka
type Storage struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}
