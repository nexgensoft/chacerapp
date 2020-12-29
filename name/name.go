package name

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	CollectionAccounts = "accounts"

	CollectionLocations = "locations"

	CollectionColors = "colors"

	CollectionDevices = "devices"

	CollectionTemplates = "templates"

	CollectionRooms = "rooms"

	CollectionMessage = "messages"
)

type parseOptions struct {
	allowWildcard bool
}

type ParseOption func(*parseOptions)

func AllowWildcard() ParseOption {
	return func(opts *parseOptions) {
		opts.allowWildcard = true
	}
}

func invalidNameError(collectionIDs ...string) error {
	return status.Error(codes.InvalidArgument, fmt.Sprintf("a valid name will be in the format of `%s/*`", strings.Join(collectionIDs, "/*/")))
}

func BuildRelativeName(parts ...string) string {
	return strings.Join(parts, "/")
}

func BuildAccount(account string) string {
	return BuildRelativeName(CollectionAccounts, account)
}

func BuildRoom(account, location, room string) string {
	return BuildRelativeName(CollectionAccounts, account, CollectionLocations, location, CollectionRooms, room)
}

func BuildLocation(account, location string) string {
	return BuildRelativeName(CollectionAccounts, account, CollectionLocations, location)
}

func ParseAccount(name string) (accountName string, err error) {
	parts, err := ParseRelativeName(name, CollectionAccounts)
	if err != nil {
		return "", err
	}

	return parts[0], nil
}

func ParseDevice(name string) (accountName, locationName, deviceName string, err error) {
	parts, err := ParseRelativeName(name, CollectionAccounts, CollectionLocations, CollectionDevices)
	if err != nil {
		return "", "", "", err
	}

	return parts[0], parts[1], parts[2], nil
}

func ParseColor(name string) (accountName, colorName string, err error) {
	parts, err := ParseRelativeName(name, CollectionAccounts, CollectionColors)
	if err != nil {
		return "", "", err
	}

	return parts[0], parts[1], nil
}

func ParseLocation(name string) (accountName, locationName string, err error) {
	parts, err := ParseRelativeName(name, CollectionAccounts, CollectionLocations)
	if err != nil {
		return "", "", err
	}

	return parts[0], parts[1], nil
}

func ParseRoom(name string) (accountName, locationName, roomName string, err error) {
	parts, err := ParseRelativeName(name, CollectionAccounts, CollectionLocations, CollectionRooms)
	if err != nil {
		return "", "", "", err
	}
	return parts[0], parts[1], parts[2], nil
}

// ValidResourceID will check if the provided ID can be used as a valid
// resource ID. A resource ID will be 4-63 characters and will only contain
// the characters "a-z", "0-9", and "-". A resource ID must not start with
// or end with a "-" and will not have any consecutive "-".
func ValidResourceID(id string) bool {
	// Validate the length of the ID.
	if len(id) < 4 || len(id) > 63 {
		return false
	}

	// Validate the ID does not start with or end with "-".
	if id[0] == '-' || id[len(id)-1] == '-' {
		return false
	}

	// Validate the characters of the resource ID.
	for i := range id {
		// Allow characters 0-9
		if id[i] >= '0' && id[i] <= '9' {
			continue
		}
		// Allow characters a-z
		if id[i] >= 'a' && id[i] <= 'z' {
			continue
		}
		// Only "-" should be allowed as long as it is
		// not consecutive.
		if id[i] != '-' || id[i-1] == '-' {
			return false
		}
	}
	return true
}

func ParseRelativeName(name string, collectionIDs ...string) ([]string, error) {
	numCollections := len(collectionIDs)

	// convert to an array of strings
	requiredCollections := make([]string, numCollections)
	for i := range collectionIDs {
		requiredCollections[i] = string(collectionIDs[i])
	}

	parts := strings.Split(name, "/")
	// A valid collection ID name MUST have twice the number
	// of parts as the number of collections minus one since
	// the last collection will not have a corresponding
	// resource ID.
	if len(parts) != (numCollections * 2) {
		return nil, invalidNameError(requiredCollections...)
	}

	// Since the last collection ID will not have a corresponding resource ID
	// we can make it one less than the length of the required collections.
	resourceIDs := make([]string, numCollections)

	// Check that each of the required collectionIDs are
	// present in the correct location
	for i := range requiredCollections {
		if requiredCollections[i] != parts[i*2] {
			return nil, invalidNameError(requiredCollections...)
		}

		// The last collection won't have a corresponding resource ID
		if i != numCollections {
			// validate the resource ID is valid
			if parts[i*2+1] != "-" && !ValidResourceID(parts[i*2+1]) {
				return nil, invalidNameError(requiredCollections...)
			}

			resourceIDs[i] = parts[i*2+1]
		}
	}

	return resourceIDs, nil
}
