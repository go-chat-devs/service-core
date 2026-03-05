package chatservice

import (
	"context"

	corev1 "github.com/go-chat-devs/proto-core-x-gateway/gen/go/core"
	"github.com/go-chat-devs/service-core/internal/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatService interface {
	CreateChat(
		ctx context.Context,
		useruid uuid.UUID,
		frienduid uuid.UUID,
	) (err error)
	DeleteChat(
		ctx context.Context,
		useruid uuid.UUID,
		frienduid uuid.UUID,
	) (err error)
	CreateGroupChat(
		ctx context.Context,
		useruid uuid.UUID,
		title string,
		bio *string,
		avataruid *uuid.UUID,
	) (err error)
	ChangeGroupChatTitle(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
		newtitle string,
	) (err error)
	ChangeGroupChatAvatar(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
		avataruid *uuid.UUID,
	) (err error)
	DeleteGroupChat(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
	) (err error)
}

type ServerApi struct {
	corev1.UnimplementedChatServiceServer
	chatService ChatService
}

func RegisterServerApi(gRPC *grpc.Server, chatService ChatService) {
	corev1.RegisterChatServiceServer(gRPC, &ServerApi{chatService: chatService})
}

func (s *ServerApi) CreateChat(ctx context.Context, req *corev1.CreateChatRequest) (*corev1.CreateChatResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	frienduid, err := uuid.Parse(req.GetFriendUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}

	err = s.chatService.CreateChat(ctx,useruid,frienduid)
	if err != nil{
		return nil, status.Error(codes.Internal,"internal server error")
	}
	return &corev1.CreateChatResponse{}, nil
}


func (s *ServerApi) DeleteChat(ctx context.Context, req *corev1.DeleteChatRequest) (*corev1.DeleteChatResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	frienduid, err := uuid.Parse(req.GetFriendUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}

	err = s.chatService.DeleteChat(ctx,useruid,frienduid)
	if err != nil{
		return nil, status.Error(codes.Internal,"internal server error")
	}
	return &corev1.DeleteChatResponse{}, nil
}




func (s *ServerApi) CreateGroupChat(ctx context.Context, req *corev1.CreateGroupChatRequest) (*corev1.CreateGroupChatResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	title := req.GetTitle()
	bio := utils.UnwrapStringValue(req.GetBio())
	var avatarUIDPtr *uuid.UUID
	if av := req.GetAvatarUid(); av != nil && av.Value != "" {
		avatarUID, err := uuid.Parse(av.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid avatar UID")
		}
		avatarUIDPtr = &avatarUID
	}
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}

	err = s.chatService.CreateGroupChat(ctx,useruid,title,bio,avatarUIDPtr)
	if err != nil{
		return nil, status.Error(codes.Internal,"internal server error")
	}
	return &corev1.CreateGroupChatResponse{}, nil
}


func (s *ServerApi) ChangeGroupChatTitle(ctx context.Context, req *corev1.ChangeGroupChatTitleRequest) (*corev1.ChangeGroupChatTitleResponse,error){
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	newtitle := req.GetNewTitle()
	err = s.chatService.ChangeGroupChatTitle(ctx,useruid,chatuid,newtitle)
	if err != nil{
		return nil, status.Error(codes.Internal,"internal server error")
	}
	return &corev1.ChangeGroupChatTitleResponse{},nil
}

func (s *ServerApi) ChangeGroupChatAvatar(ctx context.Context, req *corev1.ChangeGroupChatAvatarRequest) (*corev1.ChangeGroupChatAvatarResponse,error){
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	var avatarUIDPtr *uuid.UUID
	if av := req.GetAvatarUid(); av != nil && av.Value != "" {
		avatarUID, err := uuid.Parse(av.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid avatar UID")
		}
		avatarUIDPtr = &avatarUID
	}
	err = s.chatService.ChangeGroupChatAvatar(ctx,useruid,chatuid,avatarUIDPtr)
	if err != nil{
		return nil, status.Error(codes.Internal,"internal server error")
	}
	return &corev1.ChangeGroupChatAvatarResponse{},nil
}


func (s *ServerApi) DeleteGroupChat(ctx context.Context, req *corev1.DeleteGroupChatRequest) (*corev1.DeleteGroupChatResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}
	
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while uuid parsing")
	}

	err = s.chatService.DeleteGroupChat(ctx,useruid,chatuid)
	if err != nil{
		return nil, status.Error(codes.Internal,"internal server error")
	}
	return &corev1.DeleteGroupChatResponse{}, nil
}


