package relations

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-chat-devs/service-core/internal/models"
	"github.com/go-chat-devs/service-core/internal/scanner"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tag = tagger.Tagger("storage-relations")

type Storage struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func New(ctx context.Context, pool *pgxpool.Pool) *Storage {
	return &Storage{
		ctx:  ctx,
		pool: pool,
	}
}

func (s *Storage) InsertRelation(userUID, friendUID string) error {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		slog.Error(tag("tx begin error: %v", err))
		return err
	}
	defer tx.Rollback(s.ctx)

	const sql = `INSERT INTO relations(user_uid, friend_uid, timestamp) VALUES($1, $2, $3)`
	timestamp := time.Now().Unix()

	if _, err := tx.Exec(s.ctx, sql, userUID, friendUID, timestamp); err != nil {
		return err
	}

	if err := tx.Commit(s.ctx); err != nil {
		slog.Error(tag("tx commit error: %v", err))
		return err
	}
	return nil
}

func (s *Storage) GetRelation(userUID, friendUID string) (*models.Relation, bool) {
	const sql = `SELECT id, user_uid, friend_uid, timestamp FROM relations WHERE user_uid=$1 AND friend_uid=$2`
	row := s.pool.QueryRow(s.ctx, sql, userUID, friendUID)
	res, err := scanner.Row[*models.Relation](row)
	if err != nil {
		return nil, false
	}
	return res, true
}
