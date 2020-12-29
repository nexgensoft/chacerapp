package server

import (
	"context"

	"github.com/chacerapp/apiserver/name"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/chacerapp/apiserver/store"
	"github.com/golang/protobuf/ptypes/empty"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (s *server) ListRooms(ctx context.Context, req *serverpb.ListRoomsRequest) (*serverpb.ListRoomsResponse, error) {
	// Validate the pagination request
	pageInfo, err := s.validatePageableRequest(req)
	if err != nil {
		return nil, err
	}

	colors, err := s.store.ListRooms(ctx, req.Parent, store.WithPageInfo(pageInfo), store.WithPageSize(req.PageSize))
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

	return &serverpb.ListRoomsResponse{
		Rooms:         colors,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *server) CreateRoom(ctx context.Context, req *serverpb.CreateRoomRequest) (*serverpb.Room, error) {
	// Check that the provided room configuration is valid
	if err := validateRoom(req.Room, true); err != nil {
		return nil, err
	}

	// Validate the parent is accurate by looking up the location
	if location, err := s.store.GetLocation(ctx, req.Parent); err != nil {
		return nil, err
	} else if location == nil {
		return nil, errNotFound
	}

	// Set the name of the room based on the provided ID
	req.Room.Name = name.BuildRelativeName(req.Parent, name.CollectionRooms, req.RoomId)

	room, err := s.store.CreateRoom(ctx, req.Room)
	if err != nil {
		return nil, err
	} else if room == nil {
		return nil, errAlreadyExists
	}

	return room, nil
}

func (s *server) GetRoom(ctx context.Context, req *serverpb.GetRoomRequest) (*serverpb.Room, error) {
	if _, _, _, err := name.ParseRoom(req.Name); err != nil {
		return nil, err
	}

	room, err := s.store.GetRoom(ctx, req.Name)
	if err != nil {
		return nil, err
	} else if room == nil {
		return nil, errNotFound
	}

	return room, nil
}

func (s *server) UpdateRoom(ctx context.Context, req *serverpb.UpdateRoomRequest) (*serverpb.Room, error) {
	if err := validateRoom(req.Room, false); err != nil {
		return nil, err
	}

	if room, err := s.store.UpdateRoom(ctx, req.Room, store.WithUpdateMask(req.UpdateMask)); err != nil {
		return nil, err
	} else if room == nil {
		return nil, errNotFound
	} else {
		return room, nil
	}
}

func (s *server) DeleteRoom(ctx context.Context, req *serverpb.DeleteRoomRequest) (*empty.Empty, error) {
	if _, _, _, err := name.ParseRoom(req.Name); err != nil {
		return nil, err
	}

	if room, err := s.store.DeleteRoom(ctx, req.Name); err != nil {
		return nil, err
	} else if room == nil {
		return nil, errNotFound
	} else {
		return &empty.Empty{}, nil
	}
}

// validateRoom will validate the data specified in a room. The create
// argument should be used to indicate if the room is being created.
func validateRoom(room *serverpb.Room, create bool) error {
	path := field.NewPath("room")
	if room == nil {
		return convertErrorList(field.ErrorList{
			field.Required(path, "room is required"),
		})
	}

	// create a new list for our errors
	var errs field.ErrorList
	errs = append(errs, validateLabels(path, room)...)
	errs = append(errs, validateAnnotations(path, room)...)
	errs = append(errs, validateDisplayName(path, room)...)
	errs = append(errs, validateDescription(path, room)...)

	// Specific checks to perform only when updating a resource
	if !create {
		// Verify the room name is valid
		if room.Name == "" {
			errs = append(errs, field.Required(path.Child("name"), "name is required"))
		} else if _, _, _, err := name.ParseRoom(room.Name); err != nil {
			errs = append(errs, field.Invalid(path.Child("name"), room.Name, err.Error()))
		}
	}

	return convertErrorList(errs)
}
