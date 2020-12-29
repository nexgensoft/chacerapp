package server

import (
	"context"

	"github.com/chacerapp/apiserver/server/serverpb"
)

func (*server) ListTemplates(ctx context.Context, req *serverpb.ListTemplatesRequest) (*serverpb.ListTemplatesResponse, error) {
	return &serverpb.ListTemplatesResponse{}, nil
}

func (*server) CreateTemplate(ctx context.Context, req *serverpb.CreateTemplateRequest) (*serverpb.Template, error) {
	return &serverpb.Template{}, nil
}
