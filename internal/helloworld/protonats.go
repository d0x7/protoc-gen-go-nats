package helloworld

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

// ============================================================================
// Context & Metadata (Header) Helpers
// ============================================================================

type headerKey struct{}

// FromContext extracts NATS headers from the context (Server side)
func HeadersFromContext(ctx context.Context) nats.Header {
	if h, ok := ctx.Value(headerKey{}).(nats.Header); ok {
		return h
	}
	return nil
}

// NewContextWithHeaders creates a context containing NATS headers (Server side mostly)
func NewContextWithHeaders(ctx context.Context, h nats.Header) context.Context {
	return context.WithValue(ctx, headerKey{}, h)
}

// For Client side: We need a way to pass headers that should be SENT
type outgoingHeaderKey struct{}

// NewOutgoingContext adds headers that should be sent with the request
func NewOutgoingContext(ctx context.Context, h nats.Header) context.Context {
	return context.WithValue(ctx, outgoingHeaderKey{}, h)
}

func HeadersFromOutgoingContext(ctx context.Context) nats.Header {
	if h, ok := ctx.Value(outgoingHeaderKey{}).(nats.Header); ok {
		return h
	}
	return nil
}

// ============================================================================
// Interceptor Types (Modeled after gRPC)
// ============================================================================

// ClientInvoker is the function that sends the actual request
type ClientInvoker func(ctx context.Context, method string, req, reply proto.Message, opts ...protonats.CallOption) error

// ClientInterceptor wraps the client call
type ClientInterceptor func(ctx context.Context, method string, req, reply proto.Message, cc *nats.Conn, invoker ClientInvoker, opts ...protonats.CallOption) error

// ServerHandler is the wrapper around the actual user implementation
type ServerHandler func(ctx context.Context, req proto.Message) (proto.Message, error)

// ServerInterceptor wraps the server handling
type ServerInterceptor func(ctx context.Context, req proto.Message, info *MethodInfo, handler ServerHandler) (proto.Message, error)

// MethodInfo provides info about the RPC to the interceptor
type MethodInfo struct {
	Subject string // e.g. "service.HelloWorldService.HelloWorld" or "service.HelloWorldService.HelloWorld.<instanceID>"
	Service string // HelloWorldService
	Method  string // HelloWorld
	//FullMethod string // e.g. "service.HelloWorldService.HelloWorld"
}
