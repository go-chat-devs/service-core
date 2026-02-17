package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/go-chat-devs/service-core/internal/models"
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
	return s.users.Insert(ctx, userUID, username, avatarUID)
}

func (s *Storage) ChangeUserUsername(ctx context.Context, userUID uuid.UUID, newUsername *string) error {
	return s.users.UpdateUsername(ctx, userUID, newUsername)
}

func (s *Storage) ChangeUserAvatar(ctx context.Context, userUID uuid.UUID, newAvatarUID *uuid.UUID) error {
	return s.users.UpdateAvatar(ctx, userUID, newAvatarUID)
}

func (s *Storage) AddFriend(ctx context.Context, userUID, friendUID uuid.UUID) (err error) {
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
	return s.relations.Delete(ctx, userUID, friendUID)
}

func (s *Storage) CreateChat(ctx context.Context, userUID, friendUID uuid.UUID) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) (err error) {
		Chats := s.chats.WithTX(tx)

		_, err = Chats.SelectUsers(ctx, [2]uuid.UUID{userUID, friendUID})
		if err == nil {
			return errors.New("chat already exists")
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		_, err = Chats.Insert(ctx, [2]uuid.UUID{userUID, friendUID})
		if err != nil {
			slog.Error(tag("chat insert error: %v", err))
			return err
		}
		return nil
	})
}
func (s *Storage) DeleteChat(ctx context.Context, userUID, friendUID uuid.UUID) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		Chat := s.chats.WithTX(tx)
		chat, err := Chat.SelectUsers(ctx, [2]uuid.UUID{userUID, friendUID})
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `DELETE FROM core.base_messages WHERE chat_uid=$1`, chat.UID); err != nil {
			return err
		}
		return nil
	})
}
func (s *Storage) CreateGroupChat(ctx context.Context, userUID uuid.UUID, title string, bio *string, avatar_uid *uuid.UUID) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		GroupChats := s.groupChats.WithTX(tx)
		GroupChatMemebers := s.groupChatMembers.WithTX(tx)
		chatUID, err := GroupChats.Insert(ctx, title, bio, avatar_uid, time.Now())
		if err != nil {
			return err
		}
		return GroupChatMemebers.Insert(ctx, chatUID, userUID, models.MemberRole_Admin)
	})
}
func (s *Storage) ChangeGroupChatTitle(ctx context.Context, userUID uuid.UUID, chatUID uuid.UUID, newTitle string) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		GroupChats := s.groupChats.WithTX(tx)
		GroupChatMemebers := s.groupChatMembers.WithTX(tx)
		member, err := GroupChatMemebers.Select(ctx, chatUID, userUID)
		if err != nil {
			return err
		}
		if member.Role != models.MemberRole_Admin {
			return errors.New("only admins can change title")
		}
		return GroupChats.UpdateTitle(ctx, chatUID, newTitle)
	})
}
func (s *Storage) ChangeGroupChatAvatar(ctx context.Context, userUID, chatUID uuid.UUID, avatarUID *uuid.UUID) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		GroupChats := s.groupChats.WithTX(tx)
		GroupChatMemebers := s.groupChatMembers.WithTX(tx)
		member, err := GroupChatMemebers.Select(ctx, chatUID, userUID)
		if err != nil {
			return err
		}
		if member.Role != models.MemberRole_Admin {
			return errors.New("only admins can change avatar")
		}
		return GroupChats.UpdateAvatar(ctx, chatUID, avatarUID)
	})
}
func (s *Storage) DeleteGroupChat(ctx context.Context, userUID, chatUID uuid.UUID) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		// BaseMessages := s.baseMessages.WithTX(tx)
		// GroupChats := s.groupChats.WithTX(tx)
		GroupChatMemebers := s.groupChatMembers.WithTX(tx)
		member, err := GroupChatMemebers.Select(ctx, chatUID, userUID)
		if err != nil {
			return err
		}
		if member.Role != models.MemberRole_Admin {
			return errors.New("only admins can delete group chat")
		}
		return nil
	})
}
func (s *Storage) SendTextMessage(ctx context.Context, userUID, chatUID uuid.UUID, text string) error
func (s *Storage) SendImageMessage(ctx context.Context, userUID, chatUID, fileUID uuid.UUID) error
func (s *Storage) ChangeMessageText(ctx context.Context, userUID, messageUID uuid.UUID, newText string) error
func (s *Storage) DeleteMessage(ctx context.Context, userUID, messageUID uuid.UUID) error {
	return db.Transaction(ctx, s.db, func(tx pgx.Tx) error {
		Messages := s.baseMessages.WithTX(tx)
		msg, err := Messages.Select(ctx, messageUID)
		if err != nil {
			return err
		}
		if msg.SenderUID == nil || *msg.SenderUID != userUID {
			return errors.New("wrong sender")
		}
		Messages.Delete(ctx, messageUID)
		return nil
	})
}
