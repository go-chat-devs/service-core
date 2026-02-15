package storage

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chat-devs/service-core/internal/storage/chats"
	"github.com/go-chat-devs/service-core/internal/storage/db"
	groupchatmembers "github.com/go-chat-devs/service-core/internal/storage/group_chat_members"
	groupchats "github.com/go-chat-devs/service-core/internal/storage/group_chats"
	basemessages "github.com/go-chat-devs/service-core/internal/storage/messages/base_messages"
	imagemessages "github.com/go-chat-devs/service-core/internal/storage/messages/image_messages"
	textmessages "github.com/go-chat-devs/service-core/internal/storage/messages/text_messages"
	"github.com/go-chat-devs/service-core/internal/storage/relations"
	"github.com/go-chat-devs/service-core/internal/storage/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db db.DBTX

	Users            *users.Storage
	Relations        *relations.Storage
	BaseMessages     *basemessages.Storage
	TextMessage      *textmessages.Storage
	ImageMessage     *imagemessages.Storage
	Chats            *chats.Storage
	GroupChats       *groupchats.Storage
	GroupChatMembers *groupchatmembers.Storage
}

func New(ctx context.Context) (*Storage, error) {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		slog.Error(fmt.Sprintf("error parsing connection config: %v", err))
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("error creating connection pool: %v", err))
		return nil, err
	}

	return &Storage{
		db: pool,

		Users:            users.New(pool),
		Relations:        relations.New(pool),
		BaseMessages:     basemessages.New(pool),
		TextMessage:      textmessages.New(pool),
		ImageMessage:     imagemessages.New(pool),
		Chats:            chats.New(pool),
		GroupChats:       groupchats.New(pool),
		GroupChatMembers: groupchatmembers.New(pool),
	}, nil
}
