package server

import (
	"github.com/chacerapp/apiserver/store"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// The maximum page size that should be used by default.
const maxGlobalPageSize int32 = 500

func (s *server) validatePageableRequest(req interface {
	GetPageSize() int32
	GetPageToken() string
}) (store.PageInfo, error) {
	var errs field.ErrorList
	if req.GetPageSize() > maxGlobalPageSize || req.GetPageSize() < 0 {
		errs = append(errs, field.Invalid(field.NewPath("page_size"), req.GetPageSize(), "page_size must be between 0 and 500 inclusive"))
	}

	var parent, filter, order string
	if r, ok := req.(interface{ GetParent() string }); ok {
		parent = r.GetParent()
	}
	if r, ok := req.(interface{ GetFilter() string }); ok {
		filter = r.GetFilter()
	}
	if r, ok := req.(interface{ GetOrder() string }); ok {
		order = r.GetOrder()
	}

	// Create a default page info when the page token is not provided from a previous request
	if req.GetPageToken() == "" {
		return store.PageInfo{
			EndCursor:  "0",
			Filter:     filter,
			Order:      order,
			RequestKey: parent,
		}, nil
	}

	token, err := s.store.ParsePageToken(req.GetPageToken())
	if err != nil {
		if s, ok := status.FromError(err); ok {
			errs = append(errs, field.Invalid(field.NewPath("page_token"), req.GetPageToken(), s.Message()))
		} else {
			// Non gRPC error means that we got an internal error, return immediately
			return store.PageInfo{}, err
		}
	}

	if token.RequestKey != parent {
		errs = append(errs, field.Invalid(field.NewPath("parent"), parent, "parent must not be changed during pagination query"))
	}
	if token.Filter != filter {
		errs = append(errs, field.Invalid(field.NewPath("filter"), filter, "filter must not be changed during pagination query"))
	}
	if token.Order != order {
		errs = append(errs, field.Invalid(field.NewPath("order"), order, "order must not be changed during pagination query"))
	}

	if len(errs) > 0 {
		return store.PageInfo{}, convertErrorList(errs)
	}

	return token, nil
}
