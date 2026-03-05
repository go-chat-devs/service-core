package userservice

import (
	"context"

	corev1 "github.com/go-chat-devs/proto-core-x-gateway/gen/go/core"
	"github.com/go-chat-devs/service-core/internal/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService interface {
	CreateUser(
		ctx context.Context,
		userUid uuid.UUID,
		username *string,
		avatarUid *uuid.UUID,
	) (err error)
	ChangeUsername(
		ctx context.Context,
		userUid uuid.UUID,
		newUsername *string,
	) (err error)
	ChangeUserAvatar(
		ctx context.Context,
		userUid uuid.UUID,
		newAvatarUid *uuid.UUID,
	) (err error)
	AddFriend(
		ctx context.Context,
		userUid uuid.UUID,
		friendUid uuid.UUID,
	) (err error)
	DeleteFriend(
		ctx context.Context,
		userUid uuid.UUID,
		friendUid uuid.UUID,
	) (err error)
}

type ServerApi struct {
	corev1.UnimplementedUserServiceServer
	userService UserService
}

func RegisterServerApi(gRPC *grpc.Server, userService UserService) {
	corev1.RegisterUserServiceServer(gRPC, &ServerApi{userService: userService})
}

func (s *ServerApi) CreateUser(ctx context.Context, req *corev1.CreateUserRequest) (*corev1.CreateUserResponse, error) {
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument,"invalid user UID")
	}

	username := utils.UnwrapStringValue(req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument,"invalid username")
	}
	var avatarUIDPtr *uuid.UUID
	if av := req.GetAvatarUid(); av != nil && av.Value != "" {
		avatarUID, err := uuid.Parse(av.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid avatar UID")
		}
		avatarUIDPtr = &avatarUID
	}

	if err != nil {
		return nil, err
	}
	err = s.userService.CreateUser(ctx, useruid, username, avatarUIDPtr)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &corev1.CreateUserResponse{},nil
}


func (s *ServerApi) ChangeUsername(ctx context.Context,req *corev1.ChangeUsernameRequest) (*corev1.ChangeUsernameResponse,error){
	userUid,err := uuid.Parse(req.GetUserUid())
	if err != nil{
		return nil,status.Error(codes.InvalidArgument,"invalid user UID")
	}
	newUsername := utils.UnwrapStringValue(req.GetNewUsername())
	err = s.userService.ChangeUsername(ctx,userUid,newUsername)
	if err != nil{
		return  nil, status.Error(codes.Internal,"internal error")
	}
	return &corev1.ChangeUsernameResponse{},nil

}


func (s *ServerApi) ChangeUserAvatar(ctx context.Context, req *corev1.ChangeUserAvatarRequest) (*corev1.ChangeUserAvatarResponse,error){
	useruid,err := uuid.Parse(req.GetUserUid())
	if err != nil{
		return nil, status.Error(codes.InvalidArgument,"error while uuid parsing")
	}
	var avatarUIDPtr *uuid.UUID
	if av := req.GetNewAvatarUid(); av != nil && av.Value != "" {
		avatarUID, err := uuid.Parse(av.Value)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid new avatar UID")
		}
		avatarUIDPtr = &avatarUID
	}
	err = s.userService.ChangeUserAvatar(ctx, useruid,avatarUIDPtr)
	
	if err != nil{
		return nil,status.Error(codes.InvalidArgument,"invalid user UID")
	}
	return &corev1.ChangeUserAvatarResponse{},nil
}



func (s *ServerApi) AddFriend(ctx context.Context, req *corev1.AddFriendRequest) (*corev1.AddFriendResponse,error){
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil{
		return nil, status.Error(codes.InvalidArgument,"error while uuid parsing")
	}
	frienduid, err :=uuid.Parse(req.GetFriendUid())
	if err != nil{
		return nil, status.Error(codes.InvalidArgument,"error while uuid parsing")
	}
	err = s.userService.AddFriend(ctx,useruid,frienduid)
	if err != nil{
		return nil,status.Error(codes.Internal,"internal error")
	}
	return &corev1.AddFriendResponse{},nil
}




func (s *ServerApi) DeleteFriend(ctx context.Context, req *corev1.DeleteFriendRequest) (*corev1.DeleteFriendResponse,error){
	useruid, err := uuid.Parse(req.GetUserUid())
	if err != nil{
		return nil, status.Error(codes.InvalidArgument,"error while uuid parsing")
	}
	frienduid, err :=uuid.Parse(req.GetFriendUid())
	if err != nil{
		return nil, status.Error(codes.InvalidArgument,"error while uuid parsing")
	}
	err = s.userService.DeleteFriend(ctx,useruid,frienduid)
	if err != nil{
		return nil,status.Error(codes.Internal,"internal error")
	}
	return &corev1.DeleteFriendResponse{},nil
}