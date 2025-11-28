package helloworld

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

// ============================================================================
// Context & Metadata (Header) Helpers - THE FIX
// ============================================================================

type headerKey struct{}         // For Server (Incoming)
type outgoingHeaderKey struct{} // For Client (Outgoing)

// --- Server Side (Incoming) ---

func HeadersFromContext(ctx context.Context) nats.Header {
	if h, ok := ctx.Value(headerKey{}).(nats.Header); ok {
		return h
	}
	return nil
}

func NewContextWithHeaders(ctx context.Context, h nats.Header) context.Context {
	return context.WithValue(ctx, headerKey{}, h)
}

// --- Client Side (Outgoing) ---

func HeadersFromOutgoingContext(ctx context.Context) nats.Header {
	if h, ok := ctx.Value(outgoingHeaderKey{}).(nats.Header); ok {
		return h
	}
	return nil
}

// WithOutgoingHeader merges new headers into the context safely.
// It does NOT overwrite existing headers from previous interceptors.
func WithOutgoingHeader(ctx context.Context, key, value string) context.Context {
	// 1. Get existing
	existing := HeadersFromOutgoingContext(ctx)

	// 2. Clone/Create (To avoid mutating the map in the parent context reference)
	newHeaders := make(nats.Header)
	if existing != nil {
		for k, v := range existing {
			newHeaders[k] = v // Deep copy the slice if you want to be 100% safe, but usually slice ref is fine here
		}
	}

	// 3. Add new
	newHeaders.Add(key, value)

	// 4. Return new context
	return context.WithValue(ctx, outgoingHeaderKey{}, newHeaders)
}

// WithOutgoingHeaders merges a whole map
func WithOutgoingHeaders(ctx context.Context, h nats.Header) context.Context {
	existing := HeadersFromOutgoingContext(ctx)
	newHeaders := make(nats.Header)

	if existing != nil {
		for k, v := range existing {
			for _, val := range v {
				newHeaders.Add(k, val)
			}
		}
	}
	for k, v := range h {
		for _, val := range v {
			newHeaders.Add(k, val)
		}
	}

	return context.WithValue(ctx, outgoingHeaderKey{}, newHeaders)
}

// ============================================================================
// Interceptor Types (Modeled after gRPC)
// ============================================================================

// ClientInvoker is the function that sends the actual request
// ClientInvoker: removed 'method' string, replaced with MethodInfo
type ClientInvoker func(ctx context.Context, info *MethodInfo, req, reply proto.Message, header nats.Header, opts ...protonats.CallOption) error

// ClientInterceptor wraps the client call
// 1. Removed nats.Conn (unnecessary noise)
// 2. Replaced 'method' string with *MethodInfo (Consistency)
// 3. Kept 'opts' because interceptors might want to ADD retry options dynamically.
type ClientInterceptor func(ctx context.Context, info *MethodInfo, req, reply proto.Message, header nats.Header, invoker ClientInvoker, opts ...protonats.CallOption) error

// ServerHandler is the wrapper around the actual user implementation
type ServerHandler func(ctx context.Context, req proto.Message) (proto.Message, error)

// ServerInterceptor wraps the server handling
type ServerInterceptor func(ctx context.Context, req proto.Message, info *MethodInfo, handler ServerHandler) (proto.Message, error)

// MethodInfo provides info about the RPC. Used by BOTH Client and Server interceptors.
type MethodInfo struct {
	Subject string // e.g. "service.HelloWorldService.HelloWorld" or "service.HelloWorldService.HelloWorld.<instanceID>"
	Service string // HelloWorldService
	Method  string // HelloWorld
	//FullMethod string // e.g. "service.HelloWorldService.HelloWorld"
}
