// Code generated by protoc-gen-auth. DO NOT EDIT!
// source: auth.proto
package example

import (
	"github.com/appootb/protobuf/go/permission"
	"github.com/appootb/protobuf/go/service"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = permission.TokenLevel_NONE_TOKEN
var _ = service.UnaryServerInterceptor

var _levelExample = map[string]permission.TokenLevel{
	"/example.example/Test1": permission.TokenLevel_LOW_TOKEN,
}

// Register scoped server.
func RegisterExampleScopeServer(auth service.Authenticator, impl service.Implementor, srv ExampleServer) error {
	// Register service required token level.
	auth.RegisterServiceTokenLevel(_levelExample)

	// Register scoped gRPC server.
	for _, grpc := range impl.GetScopedGRPCServer(permission.VisibleScope_DEFAULT_SCOPE) {
		RegisterExampleServer(grpc, srv)
	} // No gateway generated.
	return nil
}

var _levelExampleB = map[string]permission.TokenLevel{
	"/example.Example_b/Test2": permission.TokenLevel_MIDDLE_TOKEN,
	"/example.Example_b/TestA": permission.TokenLevel_INNER_TOKEN,
}

// Register scoped server.
func RegisterExampleBScopeServer(auth service.Authenticator, impl service.Implementor, srv ExampleBServer) error {
	// Register service required token level.
	auth.RegisterServiceTokenLevel(_levelExampleB)

	// Register scoped gRPC server.
	for _, grpc := range impl.GetScopedGRPCServer(permission.VisibleScope_ALL_SCOPES) {
		RegisterExampleBServer(grpc, srv)
	}
	// Register scoped gateway handler.
	return impl.RegisterGateway(permission.VisibleScope_ALL_SCOPES, RegisterExampleBHandler)
}
