package shared

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"xiam.li/protonats/go/protonats"
)

var missingTracingId = protonats.NewServerErr("1234", "missing tracing id")

func TimingInterceptor(name string) protonats.UnaryClientInterceptor {
	return func(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, header nats.Header, invoker protonats.UnaryInvoker, opts ...protonats.CallOption) error {
		fmt.Printf("TimingInterceptor (%s): method=%s\n", name, info.Method)
		start := time.Now()
		err := invoker(ctx, info, req, reply, header, opts...)
		duration := time.Since(start)
		fmt.Printf("ProtoNATS (%s): %s, duration: %s, err: %v\n", name, info.Method, duration, err)
		return err
	}
}

func TimingClientInterceptor(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, header nats.Header, invoker protonats.UnaryInvoker, opts ...protonats.CallOption) error {
	fmt.Printf("TimingClientInterceptor: method=%s\n", info.Method)
	start := time.Now()
	err := invoker(ctx, info, req, reply, header, opts...)
	duration := time.Since(start)
	fmt.Printf("ProtoNATS: %s, duration: %s, err: %v\n", info.Method, duration, err)
	return err
}

func OTelClientInterceptor(ctx context.Context, info *protonats.MethodInfo, req, reply proto.Message, header nats.Header, invoker protonats.UnaryInvoker, opts ...protonats.CallOption) error {
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

func OTelServerInterceptor(ctx context.Context, req proto.Message, info *protonats.MethodInfo, handler protonats.UnaryHandler) (proto.Message, error) {
	headers := protonats.HeadersFromContext(ctx)
	traceID := headers.Get("traceID")
	if traceID == "" {
		fmt.Println("OTelServerInterceptor: traceID not found in headers")
		return nil, missingTracingId
	} else {
		ctx = context.WithValue(ctx, "traceID", traceID)
		fmt.Printf("OTelServerInterceptor: traceID=%s, method=%s\n", traceID, info.Method)
	}
	fmt.Printf("MethodInfo: Subject=%s, Service=%s, Method=%s\n", info.Subject, info.Service, info.Method)

	return handler(ctx, req)
}
