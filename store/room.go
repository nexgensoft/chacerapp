package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chacerapp/apiserver/name"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lib/pq"
)

// Room provides a storage implementation for managing rooms within storage
type Room interface {
	// GetRoom will retrieve an Room by name from storage
	//
	// This function will return a nil Room when an Room does not
	// exist with the given name. An error will only be returned when
	// the Room failed to be retrieved.
	GetRoom(ctx context.Context, name string) (*serverpb.Room, error)
	ListRooms(ctx context.Context, parent string, opts ...ListOption) ([]*serverpb.Room, error)
	CreateRoom(ctx context.Context, Room *serverpb.Room) (*serverpb.Room, error)
	UpdateRoom(ctx context.Context, Room *serverpb.Room, opts ...UpdateOption) (*serverpb.Room, error)
	DeleteRoom(ctx context.Context, name string) (*serverpb.Room, error)
}

func (s *store) GetRoom(ctx context.Context, name string) (*serverpb.Room, error) {
	return doGetRoom(ctx, s.db, name)
}

func (s *store) ListRooms(ctx context.Context, parent string, opts ...ListOption) ([]*serverpb.Room, error) {
	options := getListOptions(opts...)

	counter := 1
	var queryParts []string
	var values []interface{}
	accountName, locationName, err := name.ParseLocation(parent)
	if err != nil {
		return nil, err
	}
	if accountName != "-" {
		queryParts = append(queryParts, fmt.Sprintf("account = $%d", counter))
		values = append(values, accountName)
		counter++
	}
	if locationName != "-" {
		queryParts = append(queryParts, fmt.Sprintf("location = $%d", counter))
		values = append(values, locationName)
	}

	// Build the query based on the parent given
	query := selectRoomBaseQuery
	// Filter the query if needed
	if len(queryParts) > 0 {
		query += fmt.Sprintf(" WHERE %s", strings.Join(queryParts, " AND "))
	}

	rows, err := s.db.Query(paginateQuery(query+" ORDER BY account, location, name", options.pageInfo, options.pageSize), values...)
	if err != nil {
		return nil, err
	}

	// Close the rows once we are done retrieving results
	defer rows.Close()

	var rooms []*serverpb.Room
	for rows.Next() {
		room, err := scanRoom(rows)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// CreateRoom will create a new room in storage
//
// Only settable fields are respected when creating an room. All other fields
// will be discarded or overwritten. The Quotas and Status of the returned room
// will be guaranteed to be set. If an room with the provided name already exists
// a nil room will be returned.
func (s *store) CreateRoom(ctx context.Context, room *serverpb.Room) (*serverpb.Room, error) {
	var newRoom *serverpb.Room

	accountName, locationName, roomName, err := name.ParseRoom(room.Name)
	if err != nil {
		return nil, err
	}

	// Run in a transaction so we can atomically check if the room already exists
	err = doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		// Check that the room doesn't already exists, when it does
		// then we should return without returning an room.
		if existing, err := doGetRoom(ctx, tx, room.Name); err != nil || existing != nil {
			return err
		}

		// Create the new room with all the defaults that should be set
		newRoom = &serverpb.Room{
			Name:        room.Name,
			CreateTime:  ptypes.TimestampNow(),
			SelfLink:    serviceName + room.Name,
			DisplayName: room.DisplayName,
			Description: room.Description,
		}

		created, err := ptypes.Timestamp(newRoom.CreateTime)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			roomInsertQuery,
			roomName,
			accountName,
			locationName,
			newRoom.DisplayName,
			newRoom.Description,
			created,
		)
		return err
	})

	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (s *store) UpdateRoom(ctx context.Context, room *serverpb.Room, opts ...UpdateOption) (*serverpb.Room, error) {
	var existing *serverpb.Room
	var err error

	options := getUpdateOptions(opts...)

	err = doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if existing, err = doGetRoom(ctx, tx, room.Name); err != nil || existing == nil {
			return err
		}

		merged, err := applyUpdateMask(existing, room, options.fieldMask)
		if err != nil {
			return err
		}

		mergedRoom := merged.(*serverpb.Room)
		// Override the values in the existing account
		existing.UpdateTime = ptypes.TimestampNow()
		existing.DisplayName = mergedRoom.DisplayName
		existing.Description = mergedRoom.Description

		updated, err := ptypes.Timestamp(existing.UpdateTime)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, updateRoomQuery, existing.DisplayName, existing.Description, updated, existing.Uid)
		return err
	})

	return existing, nil
}

func (s *store) DeleteRoom(ctx context.Context, name string) (*serverpb.Room, error) {
	var room *serverpb.Room

	// Run in a transaction so we can atomically check if the room already exists
	err := doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		// Check if the room exists
		if room, err = doGetRoom(ctx, tx, name); err != nil {
			return err
		} else if room == nil {
			// return nil here so we can indicate the room does not exist in the system
			return nil
		}

		_, err = tx.Exec(roomDeleteQuery, room.Uid)
		return err
	})

	if err != nil {
		return nil, err
	}

	return room, nil
}

func doGetRoom(ctx context.Context, query retriever, fullyQualifiedName string) (*serverpb.Room, error) {
	accountName, locationName, roomName, err := name.ParseRoom(fullyQualifiedName)
	if err != nil {
		return nil, err
	}

	rows := query.QueryRowContext(ctx, selectRoomBaseQuery+` WHERE account = $1 AND location = $2 AND name = $3`, accountName, locationName, roomName)
	room, err := scanRoom(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return room, nil
}

func scanRoom(scan scanner) (*serverpb.Room, error) {
	// Allocate all the variables we will need to scan
	var uid, roomName, account, location, displayName, description string
	var createdTime time.Time
	var updateTime pq.NullTime
	// Scan the row from the database
	if err := scan.Scan(&uid, &roomName, &account, &location, &displayName, &description, &createdTime, &updateTime); err != nil {
		return nil, err
	}

	created, err := ptypes.TimestampProto(createdTime)
	if err != nil {
		return nil, err
	}

	var updated *timestamp.Timestamp
	if updateTime.Valid {
		updated, err = ptypes.TimestampProto(updateTime.Time)
		if err != nil {
			return nil, err
		}
	}

	fqn := name.BuildRoom(account, location, roomName)

	return &serverpb.Room{
		Uid:         uid,
		SelfLink:    serviceName + fqn,
		Name:        fqn,
		CreateTime:  created,
		UpdateTime:  updated,
		DisplayName: displayName,
		Description: description,
	}, nil
}

const selectRoomBaseQuery = `
SELECT id, name, account, location, display_name, description, created_time, updated_time FROM room`

const roomInsertQuery = `
INSERT INTO room (name, account, location, display_name, description, created_time, updated_time)
VALUES ($1, $2, $3, $4, $5, $6, NULL)`

const roomDeleteQuery = `
DELETE FROM room WHERE id = $1`

const updateRoomQuery = `
UPDATE room SET display_name = $1, description = $2, updated_time = $3 WHERE id = $4`
