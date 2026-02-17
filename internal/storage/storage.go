package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-chat-devs/service-core/internal/storage/chats"
	"github.com/go-chat-devs/service-core/internal/storage/db"
	groupchatmembers "github.com/go-chat-devs/service-core/internal/storage/group_chat_members"
	groupchats "github.com/go-chat-devs/service-core/internal/storage/group_chats"
	basemessages "github.com/go-chat-devs/service-core/internal/storage/messages/base_messages"
	imagemessages "github.com/go-chat-devs/service-core/internal/storage/messages/image_messages"
	textmessages "github.com/go-chat-devs/service-core/internal/storage/messages/text_messages"
	"github.com/go-chat-devs/service-core/internal/storage/relations"
	"github.com/go-chat-devs/service-core/internal/storage/users"
	"github.com/go-chat-devs/service-core/internal/tagger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tag = tagger.Tagger("storage")

type Storage struct {
	db *pgxpool.Pool

	users            *users.Storage
	relations        *relations.Storage
	baseMessages     *basemessages.Storage
	textMessage      *textmessages.Storage
	imageMessage     *imagemessages.Storage
	chats            *chats.Storage
	groupChats       *groupchats.Storage
	groupChatMembers *groupchatmembers.Storage
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

		users:            users.New(pool),
		relations:        relations.New(pool),
		baseMessages:     basemessages.New(pool),
		textMessage:      textmessages.New(pool),
		imageMessage:     imagemessages.New(pool),
		chats:            chats.New(pool),
		groupChats:       groupchats.New(pool),
		groupChatMembers: groupchatmembers.New(pool),
	}, nil
}

func (s *Storage) CreateUser(ctx context.Context, userUID uuid.UUID, username *string, avatarUID *uuid.UUID) error {
	slog.Debug("CreateUser begin")
	defer slog.Debug("CreateUser end")
	return s.users.Insert(ctx, userUID, username, avatarUID)
}

func (s *Storage) ChangeUserUsername(ctx context.Context, userUID uuid.UUID, newUsername *string) error {
	slog.Debug("ChangeUserUsername begin")
	defer slog.Debug("ChangeUserUsername end")
	return s.users.UpdateUsername(ctx, userUID, newUsername)
}

func (s *Storage) ChangeUserAvatar(ctx context.Context, userUID uuid.UUID, newAvatarUID *uuid.UUID) error {
	slog.Debug("ChangeUserAvatar begin")
	defer slog.Debug("ChangeUserAvatar end")
	return s.users.UpdateAvatar(ctx, userUID, newAvatarUID)
}

func (s *Storage) AddFriend(ctx context.Context, userUID, friendUID uuid.UUID) (err error) {
	slog.Debug("AddFriend begin")
	defer slog.Debug("AddFriend end")
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) (err error) {
		users := s.users.WithTX(tx)
		relations := s.relations.WithTX(tx)

		friend, err := users.Select(ctx, friendUID)
		if err != nil {
			return err
		}
		err = relations.Insert(ctx, userUID, friend.UserUID, time.Now())
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *Storage) DeleteFriend(ctx context.Context, userUID uuid.UUID, friendUID uuid.UUID) error {
	slog.Debug("DeleteFriend begin")
	defer slog.Debug("DeleteFriend end")
	return s.relations.Delete(ctx, userUID, friendUID)
}

func (s *Storage) CreateChat(ctx context.Context, userUID, friendUID uuid.UUID) error {
	slog.Debug("CreateChat begin")
	defer slog.Debug("CreateChat end")
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) (err error) {
		Chats := s.chats.WithTX(tx)

		_, err = Chats.SelectUsers(ctx, [2]uuid.UUID{userUID, friendUID})
		if err == nil {
			return errors.New("chat already exists")
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		err = Chats.Insert(ctx, [2]uuid.UUID{userUID, friendUID})
		if err != nil {
			slog.Error(tag("chat insert error: %v", err))
			return err
		}
		return nil
	})
}
func (s *Storage) DeleteChat(ctx context.Context, userUID, friendUID uuid.UUID) error {
	slog.Debug("DeleteChat begin")
	defer slog.Debug("DeleteChat end")
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		// TODO
		return nil
	})
}
func (s *Storage) CreateGroupChat(ctx context.Context, userUID uuid.UUID, title string) error
func (s *Storage) ChangeGroupChatTitle(ctx context.Context, userUID uuid.UUID, chatUID uuid.UUID, newTitle string) error
func (s *Storage) ChangeGroupChatAvatar(ctx context.Context, userUID uuid.UUID, chatUID, avatarUID uuid.UUID) error
func (s *Storage) DeleteGroupChat(ctx context.Context, userUID, chatUID uuid.UUID) error
func (s *Storage) SendTextMessage(ctx context.Context, userUID, chatUID uuid.UUID, text string) error
func (s *Storage) SendImageMessage(ctx context.Context, userUID, chatUID, fileUID uuid.UUID) error
func (s *Storage) ChangeMessageText(ctx context.Context, userUID, messageUID uuid.UUID, newText string) error
func (s *Storage) DeleteMessage(ctx context.Context, userUID, messageUID uuid.UUID) error {
	slog.Debug("DeleteChat begin")
	defer slog.Debug("DeleteChat end")
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		// TODO
		return nil
	})
}
