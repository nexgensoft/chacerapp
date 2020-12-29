package server

import (
	"context"

	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/golang/protobuf/ptypes/empty"
)

func (*server) ListMessages(ctx context.Context, req *serverpb.ListMessagesRequest) (*serverpb.ListMessagesResponse, error) {
	return &serverpb.ListMessagesResponse{}, nil
}

func (*server) SendMessage(ctx context.Context, req *serverpb.SendMessageRequest) (*serverpb.Message, error) {
	return &serverpb.Message{}, nil
}

func (*server) GenerateMessage(ctx context.Context, req *serverpb.GenerateMessageRequest) (*serverpb.Message, error) {
	return &serverpb.Message{}, nil
}

func (*server) CompleteMessage(ctx context.Context, req *serverpb.CompleteMessageRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (*server) CancelMessage(ctx context.Context, req *serverpb.CancelMessageRequest) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
