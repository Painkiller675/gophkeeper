package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor client interceptor for authorization
type AuthInterceptor struct {
	token string
}

// NewAuthInterceptor creates a new AuthInterceptor interceptor
func NewAuthInterceptor(token string) *AuthInterceptor {
	return &AuthInterceptor{token: token}
}

// Unary returns client interceptor-func for unary gRPC requests' authorization
func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+interceptor.token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
