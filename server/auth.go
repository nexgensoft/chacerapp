package server

import (
	"context"

	"github.com/chacerapp/apiserver/server/serverpb"
)

// func (s *server) Register(ctx context.Context, req *serverpb.RegisterRequest) (*serverpb.RegisterResponse, error) {
// 	return &serverpb.RegisterResponse{}, nil
// }

// func (s *server) Login(ctx context.Context, req *serverpb.LoginRequest) (*serverpb.LoginResponse, error) {
// 	return &serverpb.LoginResponse{}, nil
// }

func (s *server) GenerateAccessToken(ctx context.Context, req *serverpb.GenerateAccessTokenRequest) (*serverpb.GenerateAccessTokenResponse, error) {
	return &serverpb.GenerateAccessTokenResponse{}, nil
}
