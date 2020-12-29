package server

import (
	"strings"

	"github.com/chacerapp/apiserver/server/serverpb"
	"github.com/chacerapp/apiserver/store"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	errNotFound      = status.Error(codes.NotFound, "not found")
	errAlreadyExists = status.Error(codes.AlreadyExists, "already exists")
)

// NewGRPCServer will create a new gRPC server
// with a default set of interceptors that should
// be used for the server.
func NewGRPCServer(storage store.Storage) *grpc.Server {
	// create a new RPC server
	rpcServer := &server{storage}
	// Create a new gRPC server
	svr := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInterceptor()),
	)

	// Register all of the services for this server
	serverpb.RegisterAccountsServer(svr, rpcServer)
	serverpb.RegisterIAMCredentialsServer(svr, rpcServer)
	serverpb.RegisterLocationsServer(svr, rpcServer)
	serverpb.RegisterMessengerServer(svr, rpcServer)
	serverpb.RegisterRoomsServer(svr, rpcServer)
	serverpb.RegisterTemplatesServer(svr, rpcServer)
	serverpb.RegisterUserManagerServer(svr, rpcServer)
	return svr
}

// APIServer exposes an interface for all available methods on the
// API server.
type APIServer interface {
	serverpb.AccountsServer
	serverpb.IAMCredentialsServer
	serverpb.LocationsServer
	serverpb.MessengerServer
	serverpb.RoomsServer
	serverpb.TemplatesServer
	serverpb.UserManagerServer
}

// New creates an API server
func New(storage store.Storage) APIServer {
	return &server{storage}
}

type server struct {
	store store.Storage
}

func convertErrorList(errs field.ErrorList) error {
	if len(errs) == 0 {
		return nil
	}

	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, len(errs))
	for i := range errs {
		fieldViolations[i] = &errdetails.BadRequest_FieldViolation{
			Field:       errs[i].Field,
			Description: errs[i].Detail,
		}
	}

	s, err := status.New(codes.InvalidArgument, "invalid argument").
		WithDetails(&errdetails.BadRequest{
			FieldViolations: fieldViolations,
		})
	if err != nil {
		return err
	}
	return s.Err()
}

func errFailedPrecondition(msg string) error {
	return status.Error(codes.FailedPrecondition, msg)
}

func slugify(str string) string {
	slugified := ""

	addedHyphen := false
	for _, char := range strings.ToLower(str) {
		// Allow the character if it's a-z
		if char >= 97 && char <= 122 {
			addedHyphen = false
			slugified += string(char)
		}

		// Allow 0-9 characters
		if char >= 48 && char <= 57 {
			addedHyphen = false
			slugified += string(char)
		}

		// Handle a "-" character
		if char == 45 && !addedHyphen {
			// Make sure we don't accidentally add a second hyphen
			addedHyphen = true
			slugified += string(char)
		}

		// Handle a space character
		if char == 32 && !addedHyphen {
			addedHyphen = true
			slugified += "-"
		}
	}

	return slugified
}

// Validate the labels of the resource are in the correct format.
func validateLabels(path *field.Path, labelsObj interface{ GetLabels() map[string]string }) field.ErrorList {
	var errs field.ErrorList
	return errs
}

// Validate the annotations of the resource are in the correct format.
func validateAnnotations(path *field.Path, annotationsObj interface{ GetAnnotations() map[string]string }) field.ErrorList {
	var errs field.ErrorList
	return errs
}

// Validate the Display Name of a resource is in the correct format.
func validateDisplayName(path *field.Path, name interface{ GetDisplayName() string }) field.ErrorList {
	if len(name.GetDisplayName()) > 255 {
		return field.ErrorList{field.Invalid(
			path.Child("displayName"),
			name.GetDisplayName(),
			"display name must not be longer than 255 characters",
		)}
	}
	return nil
}

// Validate the description of a resource.
func validateDescription(path *field.Path, desc interface{ GetDescription() string }) field.ErrorList {
	if len(desc.GetDescription()) > 255 {
		return field.ErrorList{field.Invalid(
			path.Child("description"),
			desc.GetDescription(),
			"description must not be longer than 255 characters",
		)}
	}
	return nil
}
