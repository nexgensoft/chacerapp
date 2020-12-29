package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chacerapp/apiserver/name"
	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/lib/pq"
)

// Location provides a storage implementation for managing locations within storage
type Location interface {
	// GetLocation will retrieve an Location by name from storage
	//
	// This function will return a nil Location when an Location does not
	// exist with the given name. An error will only be returned when
	// the Location failed to be retrieved.
	GetLocation(ctx context.Context, name string) (*serverpb.Location, error)
	ListLocations(ctx context.Context, parent string, opts ...ListOption) ([]*serverpb.Location, error)
	CreateLocation(ctx context.Context, Location *serverpb.Location) (*serverpb.Location, error)
	UpdateLocation(ctx context.Context, Location *serverpb.Location, opts ...UpdateOption) (*serverpb.Location, error)
	DeleteLocation(ctx context.Context, name string) (*serverpb.Location, error)
}

func (s *store) GetLocation(ctx context.Context, name string) (*serverpb.Location, error) {
	return doGetLocation(ctx, s.db, name)
}

func (s *store) ListLocations(ctx context.Context, parent string, opts ...ListOption) ([]*serverpb.Location, error) {
	options := getListOptions(opts...)

	// Build the query based on the parent given
	baseQuery := locationSelectBaseQuery
	if accountName, err := name.ParseAccount(parent); err != nil {
		return nil, err
	} else if accountName != "-" {
		baseQuery += fmt.Sprintf(" WHERE account = '%s' ", accountName)
	}

	rows, err := s.db.Query(paginateQuery(baseQuery+" ORDER BY account, name", options.pageInfo, options.pageSize))
	if err != nil {
		return nil, err
	}

	// Close the rows once we are done retrieving results
	defer rows.Close()

	var locations []*serverpb.Location
	for rows.Next() {
		location, err := scanLocation(rows)
		if err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

// CreateLocation will create a new location in storage
//
// Only settable fields are respected when creating an location. All other fields
// will be discarded or overwritten. The Quotas and Status of the returned location
// will be guaranteed to be set. If an location with the provided name already exists
// a nil location will be returned.
func (s *store) CreateLocation(ctx context.Context, location *serverpb.Location) (*serverpb.Location, error) {
	var newLocation *serverpb.Location

	accountName, _, err := name.ParseLocation(location.Name)
	if err != nil {
		return nil, err
	}

	// Run in a transaction so we can atomically check if the location already exists
	err = doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		// Check that the location doesn't already exists, when it does
		// then we should return without returning an location.
		if existing, err := doGetLocation(ctx, tx, location.Name); err != nil || existing != nil {
			return err
		}

		// Create the new location with all the defaults that should be set
		newLocation = &serverpb.Location{
			Name:        location.Name,
			DisplayName: location.DisplayName,
			Description: location.Description,
			CreatedTime: ptypes.TimestampNow(),
			SelfLink:    serviceName + location.Name,
		}
		created, err := ptypes.Timestamp(newLocation.CreatedTime)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(
			ctx,
			locationInsertQuery,
			newLocation.Name,
			accountName,
			newLocation.DisplayName,
			newLocation.Description,
			created,
		)
		return err
	})

	if err != nil {
		return nil, err
	}

	return newLocation, nil
}

func (s *store) UpdateLocation(ctx context.Context, location *serverpb.Location, opts ...UpdateOption) (*serverpb.Location, error) {
	var existing *serverpb.Location
	var err error

	options := getUpdateOptions(opts...)

	err = doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		if existing, err = doGetLocation(ctx, tx, location.Name); err != nil || existing == nil {
			return err
		}

		merged, err := applyUpdateMask(existing, location, options.fieldMask)
		if err != nil {
			return err
		}

		mergedLocation := merged.(*serverpb.Location)
		// Override the values in the existing account
		existing.UpdatedTime = ptypes.TimestampNow()
		existing.DisplayName = mergedLocation.DisplayName
		existing.Description = mergedLocation.Description

		updated, err := ptypes.Timestamp(existing.UpdatedTime)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, locationUpdateQuery, existing.DisplayName, existing.Description, updated, existing.Name)
		return err
	})

	return existing, nil
}

func (s *store) DeleteLocation(ctx context.Context, name string) (*serverpb.Location, error) {
	var location *serverpb.Location

	// Run in a transaction so we can atomically check if the location already exists
	err := doTransaction(ctx, s.db, func(ctx context.Context, tx *sql.Tx) error {
		var err error
		// Check if the location exists
		if location, err = doGetLocation(ctx, tx, name); err != nil {
			return err
		} else if location == nil {
			// return nil here so we can indicate the location does not exist in the system
			return nil
		}

		_, err = tx.Exec(locationDeleteQuery, name)
		return err
	})

	if err != nil {
		return nil, err
	}

	return location, nil
}

func doGetLocation(ctx context.Context, query retriever, name string) (*serverpb.Location, error) {
	rows := query.QueryRowContext(ctx, locationSelectBaseQuery+` WHERE name = $1`, name)
	location, err := scanLocation(rows)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return location, nil
}

func scanLocation(scan scanner) (*serverpb.Location, error) {
	// Allocate all the variables we will need to scan
	var name, displayName, description string
	var createdTime time.Time
	var updatedTime pq.NullTime
	// Scan the row from the database
	if err := scan.Scan(&name, &displayName, &description, &createdTime, &updatedTime); err != nil {
		return nil, err
	}

	created, err := ptypes.TimestampProto(createdTime)
	if err != nil {
		return nil, err
	}

	var updated *timestamp.Timestamp
	if updatedTime.Valid {
		updated, err = ptypes.TimestampProto(updatedTime.Time)
		if err != nil {
			return nil, err
		}
	}

	return &serverpb.Location{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		SelfLink:    serviceName + name,
		CreatedTime: created,
		UpdatedTime: updated,
	}, nil
}

const locationSelectBaseQuery = `
SELECT name, display_name, description, created_time, updated_time FROM location`

const locationInsertQuery = `
INSERT INTO location (name, account, display_name, description, created_time, updated_time)
VALUES ($1, $2, $3, $4, $5, NULL)`

const locationDeleteQuery = `
DELETE FROM location WHERE name = $1`

const locationUpdateQuery = `
UPDATE location SET display_name = $1, description = $2, updated_time = $3 WHERE name = $4`
