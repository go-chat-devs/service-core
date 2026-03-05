package messageservice

import (
	"context"

	corev1 "github.com/go-chat-devs/proto-core-x-gateway/gen/go/core"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageService interface {
	SendChatTextMessage(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
		text string,
	) (sentAt *timestamppb.Timestamp, err error)

	SendGroupChatTextMessage(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
		text string,
	) (sentAt *timestamppb.Timestamp, err error)

	SendChatImageMessage(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
		fileuid uuid.UUID,
	) (sent_at *timestamppb.Timestamp, err error)

	SendGroupImageMessage(
		ctx context.Context,
		useruid uuid.UUID,
		chatuid uuid.UUID,
		fileuid uuid.UUID,
	) (sentAt *timestamppb.Timestamp, err error)

	ChangeMessageText(
		ctx context.Context,
		useruid uuid.UUID,
		messageuid uuid.UUID,
		newcontent string,
	) (updatedAt *timestamppb.Timestamp, err error)

	DeleteMessage(
		ctx context.Context,
		useruid uuid.UUID,
		messageuid uuid.UUID,
	) (err error)
}

type ServerApi struct {
	corev1.UnimplementedMessageServiceServer
	messageService MessageService
}

func RegisterServerApi(gRPC *grpc.Server, messageService MessageService) {
	corev1.RegisterMessageServiceServer(gRPC, &ServerApi{messageService: messageService})
}

func (s *ServerApi) SendChatTextMessage(ctx context.Context, req *corev1.SendChatTextMessageRequest) (*corev1.SendChatTextMessageResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	text := req.GetText()
	sentAt, err := s.messageService.SendChatTextMessage(ctx, useruid, chatuid, text)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &corev1.SendChatTextMessageResponse{SentAt: sentAt}, nil
}

func (s *ServerApi) SendGroupChatTextMessage(ctx context.Context, req *corev1.SendGroupChatTextMessageRequest) (*corev1.SendGroupChatTextMessageResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	text := req.GetText()
	sentAt, err := s.messageService.SendGroupChatTextMessage(ctx, useruid, chatuid, text)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &corev1.SendGroupChatTextMessageResponse{SentAt: sentAt}, nil
}

func (s *ServerApi) SendChatImageMessage(ctx context.Context, req *corev1.SendChatImageMessageRequest) (*corev1.SendChatImageMessageResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	fileuid, err := uuid.Parse(req.GetFileUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	sentAt, err := s.messageService.SendChatImageMessage(ctx, useruid, chatuid, fileuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &corev1.SendChatImageMessageResponse{SentAt: sentAt}, nil
}

func (s *ServerApi) SendGroupImageMessage(ctx context.Context, req *corev1.SendGroupImageMessageRequest) (*corev1.SendGroupImageMessageResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	chatuid, err := uuid.Parse(req.GetChatUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	fileuid, err := uuid.Parse(req.GetFileUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	sentAt, err := s.messageService.SendGroupImageMessage(ctx, useruid, chatuid, fileuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &corev1.SendGroupImageMessageResponse{SentAt: sentAt}, nil
}

func (s *ServerApi) ChangeMessageText(ctx context.Context, req *corev1.ChangeMessageTextRequest) (*corev1.ChangeMessageTextResponse, error) {
	messageuid, err := uuid.Parse(req.GetMessageUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	newContent := req.GetNewContent()
	updatedAt, err := s.messageService.ChangeMessageText(ctx, useruid, messageuid, newContent)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &corev1.ChangeMessageTextResponse{UpdatedAt: updatedAt}, nil
}

func (s *ServerApi) DeleteMessage(ctx context.Context, req *corev1.DeleteMessageRequest) (*corev1.DeleteMessageResponse, error) {
	messageuid, err := uuid.Parse(req.GetMessageUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "error while parsing user uid")
	}
	err = s.messageService.DeleteMessage(ctx, useruid, messageuid)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &corev1.DeleteMessageResponse{}, nil
}
