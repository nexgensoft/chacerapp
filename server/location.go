package server

import (
	"context"

	"github.com/chacerapp/apiserver/name"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/chacerapp/apiserver/store"
	"github.com/golang/protobuf/ptypes/empty"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (s *server) ListLocations(ctx context.Context, req *serverpb.ListLocationsRequest) (*serverpb.ListLocationsResponse, error) {
	// Validate the pagination request
	pageInfo, err := s.validatePageableRequest(req)
	if err != nil {
		return nil, err
	}

	colors, err := s.store.ListLocations(ctx, req.Parent, store.WithPageInfo(pageInfo), store.WithPageSize(req.PageSize))
	if err != nil {
		return nil, err
	}

	var nextPageToken string
	// The next page token should only be generated when the number
	// of results being returned is equal to the page size. The lack
	// of a next page token is used to determine if a next page exists.
	if len(colors) == int(req.PageSize) {
		nextPageToken, err = s.store.GenerateNextPageToken(pageInfo, req.PageSize)
		if err != nil {
			return nil, err
		}
	}

	return &serverpb.ListLocationsResponse{
		Locations:     colors,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *server) CreateLocation(ctx context.Context, req *serverpb.CreateLocationRequest) (*serverpb.Location, error) {
	// Check that the provided location configuration is valid
	if err := validateCreateLocation(req.Location); err != nil {
		return nil, err
	}

	// Validate the parent is accurate by looking up the account
	account, err := s.store.GetAccount(ctx, req.Parent)
	if err != nil {
		return nil, err
	} else if account == nil {
		return nil, errNotFound
	}

	// Set the name of the location based on the display name
	req.Location.Name = name.BuildRelativeName(req.Parent, name.CollectionLocations, slugify(req.Location.DisplayName))

	location, err := s.store.CreateLocation(ctx, req.Location)
	if err != nil {
		return nil, err
	} else if location == nil {
		return nil, errAlreadyExists
	}

	return location, nil
}

func (s *server) GetLocation(ctx context.Context, req *serverpb.GetLocationRequest) (*serverpb.Location, error) {
	if _, _, err := name.ParseLocation(req.Name); err != nil {
		return nil, err
	}

	location, err := s.store.GetLocation(ctx, req.Name)
	if err != nil {
		return nil, err
	} else if location == nil {
		return nil, errNotFound
	}

	return location, nil
}

func (s *server) UpdateLocation(ctx context.Context, req *serverpb.UpdateLocationRequest) (*serverpb.Location, error) {
	if err := validateCreateLocation(req.Location); err != nil {
		return nil, err
	}

	if location, err := s.store.UpdateLocation(ctx, req.Location, store.WithUpdateMask(req.UpdateMask)); err != nil {
		return nil, err
	} else if location == nil {
		return nil, errNotFound
	} else {
		return location, nil
	}
}

func (s *server) DeleteLocation(ctx context.Context, req *serverpb.DeleteLocationRequest) (*empty.Empty, error) {
	if _, _, err := name.ParseLocation(req.Name); err != nil {
		return nil, err
	}

	if location, err := s.store.DeleteLocation(ctx, req.Name); err != nil {
		return nil, err
	} else if location == nil {
		return nil, errNotFound
	} else {
		return &empty.Empty{}, nil
	}
}

func validateCreateLocation(location *serverpb.Location) error {
	path := field.NewPath("location")
	if location == nil {
		return convertErrorList(field.ErrorList{
			field.Required(path, "location is required"),
		})
	}

	// create a new list for our errors
	var errs field.ErrorList
	if len(location.Description) > 255 {
		errs = append(errs, field.Invalid(path.Child("description"), location.Description, "description must not be longer than 255 characters"))
	}

	return convertErrorList(errs)
}
