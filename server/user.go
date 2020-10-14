package server

import (
	"context"

	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *server) ListUsers(ctx context.Context, req *serverpb.ListUsersRequest) (*serverpb.ListUsersResponse, error) {
	return &serverpb.ListUsersResponse{}, nil
}

func (s *server) CreateUser(ctx context.Context, req *serverpb.CreateUserRequest) (*serverpb.User, error) {
	return &serverpb.User{}, nil
}

func (s *server) GetUser(ctx context.Context, req *serverpb.GetUserRequest) (*serverpb.User, error) {
	return &serverpb.User{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *serverpb.DeleteUserRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *serverpb.UpdateUserRequest) (*serverpb.User, error) {
	return &serverpb.User{}, nil
}
