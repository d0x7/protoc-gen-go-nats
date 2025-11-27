package shared

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/go-protonats/internal/helloworld"
	"xiam.li/protonats/go/protonats"
)

func OTelClientInterceptor(ctx context.Context, method string, req, reply proto.Message, cc *nats.Conn, invoker helloworld.ClientInvoker, opts ...protonats.CallOption) error {
	value := ctx.Value("traceID")
	if value == nil {
		fmt.Println("traceID not found in context")
		return invoker(ctx, method, req, reply, opts...)
	}
	traceID, ok := value.(string)
	if !ok {
		fmt.Println("traceID is not a string")
		return invoker(ctx, method, req, reply, opts...)
	}

	fmt.Printf("OTelClientInterceptor: traceID=%s, method=%s\n", traceID, method)
	headers := nats.Header{}
	ctx = helloworld.NewOutgoingContext(ctx, headers)

	err := invoker(ctx, method, req, reply, opts...)
	if err != nil {
		fmt.Printf("OTelClientInterceptor: error=%v\n", err)
	} else {
		fmt.Println("OTelClientInterceptor: call successful")
	}
	return err
}

func OTelServerInterceptor(ctx context.Context, req proto.Message, info *helloworld.MethodInfo, handler helloworld.ServerHandler) (proto.Message, error) {
	headers := helloworld.HeadersFromContext(ctx)
	traceID := headers.Get("traceID")
	if traceID == "" {
		fmt.Println("OTelServerInterceptor: traceID not found in headers")
	} else {
		fmt.Printf("OTelServerInterceptor: traceID=%s, method=%s\n", traceID, info.FullMethod)
	}

	return handler(ctx, req)
}
