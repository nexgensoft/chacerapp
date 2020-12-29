package server_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/chacerapp/apiserver/server"
	"github.com/chacerapp/apiserver/store"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v10"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	_ "github.com/lib/pq"
	"github.com/stretchr/objx"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/anypb"
)

var opt = godog.Options{
	Output:   colors.Colored(os.Stdout),
	Format:   "progress",
	Paths:    []string{"../features"},
	NoColors: true,
	Tags:     "",
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

type serverFeature struct {
	responseError error
	server        *grpc.Server
	listener      net.Listener
	clientConn    grpc.ClientConnInterface
	request       interface{}
	response      interface{}
	nextPageToken string
	ctx           context.Context
	db            *sql.DB
}

func TestMain(m *testing.M) {
	flag.Parse()

	txdb.Register(
		"txdb",
		"postgres",
		"postgres://root@localhost:26257/chacerapp_tests?sslmode=disable",
		txdb.SavePointOption(nil),
	)

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, opt)

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func (f *serverFeature) aJSONgRPCRequest(protoMessageType string, jsonMessage *messages.PickleStepArgument_PickleDocString) error {
	t := proto.MessageType(protoMessageType)
	if t == nil {
		return fmt.Errorf("unknown protobuf message %q", protoMessageType)
	}

	// Grab a new instance of the proto message
	request := reflect.New(t.Elem()).Interface()
	if err := jsonpb.Unmarshal(strings.NewReader(jsonMessage.Content), request.(proto.Message)); err != nil {
		return fmt.Errorf("failed to unmarshal message into %q: %v", protoMessageType, err)
	}
	f.request = request
	return nil
}

func (f *serverFeature) callingTheRPC(method string) error {
	// clear out entries from previous calls
	f.response = nil
	f.responseError = nil

	// Split the full name into its parts
	nameParts := strings.Split(method, "/")
	serviceDescr, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(nameParts[0]))
	if err != nil {
		return fmt.Errorf("unable to find service descriptor for method %v: %v", method, err)
	}
	// Grab the method descriptor
	methodDescr := serviceDescr.(protoreflect.ServiceDescriptor).Methods().ByName(protoreflect.Name(nameParts[1]))
	// Grab the response message descriptor so we can use the type information
	responseMessageName := string(methodDescr.(protoreflect.MethodDescriptor).Output().FullName())
	// Grab a new instance of the proto response message. This should be guaranteed to be
	// registered in the Proto registry since it came from the method descriptor
	f.response = reflect.New(proto.MessageType(responseMessageName).Elem()).Interface()
	// Copy the request so that any modifications to the request object
	// by the RPC do not effect the object we have stored.
	request := proto.Clone(f.request.(proto.Message))
	// Invoke the API call
	f.responseError = f.clientConn.Invoke(f.ctx, method, request, f.response)

	return nil
}

func (f *serverFeature) iWillReceiveAnErrorWithCode(stringCode string) error {
	expectedCode := new(codes.Code)
	if err := expectedCode.UnmarshalJSON([]byte(stringCode)); err != nil {
		return fmt.Errorf("invalid gRPC code: %v", err)
	}

	if s, ok := status.FromError(f.responseError); !ok {
		return fmt.Errorf("error was not able to be converted to a gRPC status: %v", f.responseError)
	} else if s.Code() != *expectedCode {
		return fmt.Errorf("received error code %q, expected code %q", s.Code().String(), expectedCode.String())
	}
	return nil
}

// Verifies that the response error contains specific error details. This will perform strict
// validation so that if any unknown error details are included in the error, the assertion will
// fail.
func (f *serverFeature) theErrorDetailsWillBeForTheFollowingFields(details *godog.Table) error {
	// Verify the response is a gRPC error
	errStatus, ok := status.FromError(f.responseError)
	if !ok {
		return fmt.Errorf("error was not able to be converted to a gRPC status: %v", f.responseError)
	}

	// Verify the error contains a BadRequest error details
	var badRequest *errdetails.BadRequest
	for _, detail := range errStatus.Details() {
		if v, ok := detail.(*errdetails.BadRequest); ok {
			badRequest = v
			break
		}
	}
	if badRequest == nil {
		return fmt.Errorf("response error did not contain a bad request error detail")
	}

	// Simple check to see if the number of field errors we're expecting matches the actual
	if len(details.Rows) != len(badRequest.FieldViolations) {
		return fmt.Errorf("expected %d field violations, got %d: %+v", len(details.Rows), len(badRequest.FieldViolations), badRequest.FieldViolations)
	}

	// Check each expected entry and compare to the actual field errors
	for _, row := range details.Rows {
		// Try to find a matching field error
		var fieldViolation *errdetails.BadRequest_FieldViolation
		for i, violation := range badRequest.FieldViolations {
			if violation.Field == row.Cells[0].Value {
				fieldViolation = badRequest.FieldViolations[i]
				break
			}
		}

		if fieldViolation == nil {
			return fmt.Errorf("unable to find field violation for field %q", row.Cells[0].Value)
		}

		if fieldViolation.Description != row.Cells[1].Value {
			return fmt.Errorf("expected error description to be %q, got %q", row.Cells[1].Value, fieldViolation.Description)
		}
	}
	return nil
}

func (f *serverFeature) iWillReceiveASuccessfulResponse() error {
	if f.responseError != nil {
		return fmt.Errorf(
			"expected a successful response, but received an error: %v: %v",
			f.responseError,
			status.Convert(f.responseError).Details(),
		)
	}
	return nil
}

func (f *serverFeature) theResponseValueWillBe(path, expected string) error {
	var jsonMarshaler = &jsonpb.Marshaler{EmitDefaults: true}
	// Turn the response to json to more easily get a map[string]interface{}
	r, _ := jsonMarshaler.MarshalToString(f.response.(proto.Message))
	actual := objx.MustFromJSON(r).Get(path).String()
	if actual != expected {
		return fmt.Errorf("expected '%s' to be '%s', got '%s'", path, expected, actual)
	}
	return nil
}

func (f *serverFeature) theResponseValueWillHaveLength(path string, expectedLen int) error {
	var jsonMarshaler = &jsonpb.Marshaler{}
	r, _ := jsonMarshaler.MarshalToString(f.response.(proto.Message))
	actual := objx.MustFromJSON(r).Get(path).Data()

	actualLen := 0
	if a, ok := actual.([]interface{}); ok {
		actualLen = len(a)
	}
	if actualLen != expectedLen {
		return fmt.Errorf("expected length of %s to be %d but it was %d", path, expectedLen, actualLen)
	}

	return nil
}

func (f *serverFeature) stashingTheNextPageTokenFromTheResponse() error {
	if t, ok := f.response.(interface{ GetNextPageToken() string }); !ok {
		return fmt.Errorf("the response does satisfy the next page token interface")
	} else if t.GetNextPageToken() == "" {
		return fmt.Errorf("the response does not contain a next page token")
	} else {
		f.nextPageToken = t.GetNextPageToken()
		return nil
	}
}

func (f *serverFeature) usingTheStashedNextPageToken() error {
	// verify the request satisfies the page token interface
	if _, ok := f.request.(interface{ GetPageToken() string }); !ok {
		return fmt.Errorf("the response does satisfy the next page token interface")
	}
	var jsonMarshaler = &jsonpb.Marshaler{}
	r, _ := jsonMarshaler.MarshalToString(f.request.(proto.Message))
	updated := objx.MustFromJSON(r).Set("pageToken", f.nextPageToken).MustJSON()
	request := reflect.ValueOf(f.request).Interface()
	if err := jsonpb.Unmarshal(strings.NewReader(updated), request.(proto.Message)); err != nil {
		return fmt.Errorf("failed to update the page token in the request: %v", err)
	}

	f.request = request
	return nil
}

// Takes a Path to a seed file in a format of a list of google.protobuf.Any resources
// that should be created against the API service.
//
// The seed file is expected to be in the following format:
//
// {
// 	 "resources": [
//     {
//       "@type": "chacerapp.v1.CreateLocationRequest",
//       "parent": "accounts/secondary",
//       "location": {
//         "displayName": "Secondary",
//         "description": "This is my secondary location"
//       },
//       "location_id": "secondary"
//     }
// 	 ]
// }
//
// This file would look for a gRPC method that accepts the "chacerapp.v1.CreateLocationRequest"
// message as input and invoke the endpoint with the given request. If the endpoint the provided
// resource cannot be found or if an error is turned from the endpoint, the entire step will fail.
func (f *serverFeature) dataLoadedFromTheSeedFile(fileName string) error {
	handle, err := os.Open("../features/" + fileName)
	if err != nil {
		return fmt.Errorf("failed to open seed file: %v", err)
	}

	// Read in the resources to create
	rawResources, err := ioutil.ReadAll(handle)
	if err != nil {
		return fmt.Errorf("failed to read contents of seed file: %v", err)
	}

	return f.resourcesCreatedFromJSON(rawResources)
}

//
// {
// 	 "resources": [
//     {
//       "@type": "chacerapp.v1.CreateLocationRequest",
//       "parent": "accounts/secondary",
//       "location": {
//         "displayName": "Secondary",
//         "description": "This is my secondary location"
//       },
//       "location_id": "secondary"
//     }
// 	 ]
// }
//
// This file would look for a gRPC method that accepts the "chacerapp.v1.CreateLocationRequest"
// message as input and invoke the endpoint with the given request. If the endpoint the provided
// resource cannot be found or if an error is turned from the endpoint, the entire step will fail.
func (f *serverFeature) dataSeededFromJSONBlob(blob *messages.PickleStepArgument_PickleDocString) error {
	return f.resourcesCreatedFromJSON([]byte(blob.Content))
}

func (f *serverFeature) resourcesCreatedFromJSON(resourcesJSON []byte) error {
	// Define a new type so we can unmarshal the seed file
	type seedResources struct {
		Resources []json.RawMessage `json:"resources"`
	}
	seedResource := &seedResources{}
	if err := json.Unmarshal(resourcesJSON, seedResource); err != nil {
		return fmt.Errorf("failed to unmarshal seed file: %v", err)
	}

	// Loop through each Resource and convert to an Any type
	for _, r := range seedResource.Resources {
		// Get the Request message based on the type specified in the message
		var any anypb.Any
		if err := jsonpb.UnmarshalString(string(r), &any); err != nil {
			return fmt.Errorf("failed to unmarshal resource into google.protobuf.Any: %v", err)
		}
		request, err := any.UnmarshalNew()
		if err != nil {
			return fmt.Errorf("failed to get request message from seed file: %v", err)
		}

		// Find the RPC method that accepts the message as an input parameter.
		var method protoreflect.MethodDescriptor
		protoregistry.GlobalFiles.RangeFiles(func(file protoreflect.FileDescriptor) bool {
			// Go through all of the registered services
			for i := 0; i < file.Services().Len(); i++ {
				srv := file.Services().Get(i)
				for j := 0; j < srv.Methods().Len(); j++ {
					serviceMethod := srv.Methods().Get(j)
					// Check if the current method input matches the same name as the request object
					if serviceMethod.Input().FullName() == request.ProtoReflect().Descriptor().FullName() {
						method = serviceMethod
						return false
					}
				}
			}
			return true
		})
		if method == nil {
			return fmt.Errorf("could not find a method in the protoregistry that accepts the message type %v", any.MessageName())
		}

		// Grab the response message descriptor so we can use the type information
		responseMessageName := string(method.Output().FullName())
		// Grab a new instance of the proto response message. This should be guaranteed to be
		// registered in the Proto registry since it came from the method descriptor
		response := reflect.New(proto.MessageType(responseMessageName).Elem()).Interface()
		// Build the fully qualified name of the method that should be invoked
		methodName := string(method.Parent().FullName()) + "/" + string(method.Name())

		// Call the endpoint with the request object ignoring the response object
		if err := f.clientConn.Invoke(f.ctx, methodName, request, response); err != nil {
			return fmt.Errorf("failed to call %v for type %v: %v", method, any.GetTypeUrl(), err)
		}
	}

	return nil
}

func (f *serverFeature) registerSteps(suite *godog.Suite) {
	suite.Step(`^a JSON "([^"]*Request)"$`, f.aJSONgRPCRequest)
	suite.Step(`^calling the "([^"]*)" RPC$`, f.callingTheRPC)
	suite.Step(`^I will receive an error with code ("[^"]*")$`, f.iWillReceiveAnErrorWithCode)
	suite.Step(`^the BadRequest error details will be for the following fields$`, f.theErrorDetailsWillBeForTheFollowingFields)
	suite.Step(`^I will receive a successful response$`, f.iWillReceiveASuccessfulResponse)
	suite.Step(`^the response value "([^"]*)" will be "([^"]*)"$`, f.theResponseValueWillBe)
	suite.Step(`^the response value "([^"]*)" will have a length of (\d+)$`, f.theResponseValueWillHaveLength)
	suite.Step(`^stashing the next page token from the response$`, f.stashingTheNextPageTokenFromTheResponse)
	suite.Step(`^using the stashed next page token$`, f.usingTheStashedNextPageToken)
	suite.Step(`^data loaded from the seed file "([^"]*)"$`, f.dataLoadedFromTheSeedFile)
	suite.Step(`^these resources are created:$`, f.dataSeededFromJSONBlob)
}

func FeatureContext(s *godog.Suite) {
	var err error

	feature := &serverFeature{}
	feature.registerSteps(s)

	s.BeforeSuite(func() {
		m, err := migrate.New("file://../migrations", "cockroachdb://root@localhost:26257/chacerapp_tests?sslmode=disable")
		if err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to migrate database: %v", err)
		}
	})

	s.BeforeScenario(func(*messages.Pickle) {
		feature.listener, err = net.Listen("tcp", ":33000")
		if err != nil {
			log.Fatalf("failed to create tcp listener: %v", err)
		}

		feature.clientConn, err = grpc.Dial(feature.listener.Addr().String(), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("failed to create client connection: %v", err)
		}
		feature.ctx = context.Background()
		feature.db, err = sql.Open("txdb", "postgres://root@localhost:26257/chacerapp_tests?sslmode=disable")
		if err != nil {
			log.Fatalf("failed to open new database connection: %v", err)
		}

		// Create a new gRPC server to run tests against
		feature.server = server.NewGRPCServer(store.New(
			feature.db,
			store.NewPaginator([]byte("my-super-secure-test-secret-3234")),
		))
		// Start the server in the background
		go feature.server.Serve(feature.listener)
	})

	s.AfterScenario(func(*messages.Pickle, error) {
		feature.listener.Close()
		feature.server.Stop()
		feature.db.Close()
	})
}
