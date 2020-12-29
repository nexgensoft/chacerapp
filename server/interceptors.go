package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/chacerapp/apiserver/server/serverpb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func authInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		name := strings.Split(info.FullMethod, "/")
		descr, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(name[1]))
		if err != nil {
			return nil, fmt.Errorf("unable to resolve method descriptor for endpoint %v: %v", info.FullMethod, err)
		}

		// Grab the descriptor for the RPC method that's being called
		methodDesc := descr.(protoreflect.ServiceDescriptor).Methods().ByName(protoreflect.Name(name[2]))

		var requiredPermissions []string
		// Grab the required permissions for the endpoint
		if proto.HasExtension(methodDesc.Options(), serverpb.E_RequiredPermissions) {
			requiredPermissions = proto.GetExtension(methodDesc.Options(), serverpb.E_RequiredPermissions).([]string)
		}

		if len(requiredPermissions) > 0 {
			// Add shit
		}

		return handler(ctx, req)
	}
}
