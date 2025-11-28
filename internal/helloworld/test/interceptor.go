package shared

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/go-protonats/internal/helloworld"
	"xiam.li/protonats/go/protonats"
)

func OTelClientInterceptor(ctx context.Context, info *helloworld.MethodInfo, req, reply proto.Message, header nats.Header, invoker helloworld.ClientInvoker, opts ...protonats.CallOption) error {
	value := ctx.Value("traceID")
	if value == nil {
		fmt.Println("traceID not found in context")
		return invoker(ctx, info, req, reply, header, opts...)
	}
	traceID, ok := value.(string)
	if !ok {
		fmt.Println("traceID is not a string")
		return invoker(ctx, info, req, reply, header, opts...)
	}

	fmt.Printf("OTelClientInterceptor: traceID=%s, method=%s\n", traceID, info.Method)
	header.Set("traceID", traceID)
	//headers := nats.Header{}
	//headers.Set("traceID", traceID)
	//ctx = helloworld.NewOutgoingContext(ctx, headers)

	err := invoker(ctx, info, req, reply, header, opts...)
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
		ctx = context.WithValue(ctx, "traceID", traceID)
		fmt.Printf("OTelServerInterceptor: traceID=%s, method=%s\n", traceID, info.Method)
	}
	fmt.Printf("MethodInfo: Subject=%s, Service=%s, Method=%s\n", info.Subject, info.Service, info.Method)

	return handler(ctx, req)
}
