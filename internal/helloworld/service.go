package helloworld

import (
	context "context"
	fmt "fmt"
	"log/slog"
	"time"

	nats_go "github.com/nats-io/nats.go"
	micro "github.com/nats-io/nats.go/micro"
	"github.com/pkg/errors"
	proto "google.golang.org/protobuf/proto"

	impl "xiam.li/protonats/go/impl"
	protonats "xiam.li/protonats/go/protonats"
)

// region Client

// 1. Interface Change: Added context.Context
type HelloWorldServiceNATSClient interface {
	HelloWorld(ctx context.Context, req *HelloWorldRequest, opts ...protonats.CallOption) (*HelloWorldResponse, error)
}

type helloWorldServiceNATSClient struct {
	nc           *nats_go.Conn
	interceptors []ClientInterceptor
}

// 2. Client Constructor now accepts Interceptors
func NewHelloWorldServiceNATSClient(nc *nats_go.Conn, interceptors ...ClientInterceptor) HelloWorldServiceNATSClient {
	return &helloWorldServiceNATSClient{
		nc:           nc,
		interceptors: interceptors,
	}
}

// 3. The Implementation of the Client Method
func (c *helloWorldServiceNATSClient) HelloWorld(ctx context.Context, req *HelloWorldRequest, opts ...protonats.CallOption) (*HelloWorldResponse, error) {
	resp := &HelloWorldResponse{}

	info := &MethodInfo{
		Subject: "service.HelloWorldService.HelloWorld",
		Method:  "HelloWorld",
		Service: "HelloWorldService",
	}

	// 1. Initialize Headers ONCE
	// We extract what the user put in Context (if any) and create the map that will be used
	// for the rest of the chain.
	var headerMap nats_go.Header
	if ctxHeaders := HeadersFromOutgoingContext(ctx); ctxHeaders != nil {
		// We must copy once to avoid mutating the immutable Context value
		headerMap = make(nats_go.Header)
		for k, v := range ctxHeaders {
			headerMap[k] = v // Slice copy
		}
	} else {
		headerMap = make(nats_go.Header)
	}

	// Define the "Final Invoker". This is the function that actually calls NATS.
	// It is the last link in the chain.
	invoker := func(ctx context.Context, info *MethodInfo, req, reply proto.Message, headers nats_go.Header, opts ...protonats.CallOption) error {
		options := impl.ProcessCallOptions(opts...)

		// Logic to handle Timeouts: Context takes precedence, but we can fallback to options
		var cancel context.CancelFunc
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			// If context has no deadline, use the default from options (e.g. 5 seconds)
			timeout := options.GetTimeoutOr(5 * time.Second)
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		// Marshal the request
		data, err := proto.Marshal(req)
		if err != nil {
			return protonats.ErrMarshallingFailed
		}

		// Create the NATS Message
		msg := &nats_go.Msg{
			Subject: options.Subject(info.Subject),
			Data:    data,
			Header:  headers,
		}

		// Perform the Request using RequestMsgWithContext (Native Context Support!)
		// This handles premature return if ctx is cancelled.
		respMsg, err := c.nc.RequestMsgWithContext(ctx, msg)
		if err != nil {
			return err
		}

		// Check NATS-Micro Error Headers
		if errMsg, errCode := respMsg.Header.Get(micro.ErrorHeader), respMsg.Header.Get(micro.ErrorCodeHeader); len(errMsg) > 0 && len(errCode) > 0 {
			return protonats.ServiceError{Code: errCode, Description: errMsg, Details: string(respMsg.Data)}
		}

		// Unmarshal the response
		if err = proto.Unmarshal(respMsg.Data, reply); err != nil {
			return protonats.ErrUnmarshallingFailed
		}

		return nil
	}

	// Chain the interceptors
	// We wrap the invoker recursively
	chain := invoker
	for i := len(c.interceptors) - 1; i >= 0; i-- {
		interceptor := c.interceptors[i]
		// Capture loop variables
		next := chain
		chain = func(currentCtx context.Context, info *MethodInfo, currentReq, reply proto.Message, headers nats_go.Header, currentOpts ...protonats.CallOption) error {
			return interceptor(currentCtx, info, currentReq, reply, headers, next, currentOpts...)
		}
	}

	// Execute the chain
	err := chain(ctx, info, req, resp, headerMap, opts...)
	return resp, err
}

//endregion

// region Server

// 1. Interface Change: Added context.Context
type HelloWorldServiceNATSServer interface {
	HelloWorld(ctx context.Context, req *HelloWorldRequest) (*HelloWorldResponse, error)
}

// Options struct for Server (to hold interceptors)
type ServerOption func(*serverOptions)

type serverOptions struct {
	interceptors []ServerInterceptor
	// other options...
}

func WithServerInterceptors(interceptors ...ServerInterceptor) ServerOption {
	return func(o *serverOptions) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

type HelloWorldServiceId interface {
	SetHelloWorldServiceId(string)
}

func NewHelloWorldServiceNATSServer(nc *nats_go.Conn, server HelloWorldServiceNATSServer, opts ...ServerOption) micro.Service {
	service, options, err := impl.NewService("HelloWorldService", nc, server)
	if err != nil {
		panic(err) // TODO: Update this to proper error handling
	}
	so := &serverOptions{}
	for _, opt := range opts {
		opt(so)
	}
	if setId, ok := server.(HelloWorldServiceId); ok {
		setId.SetHelloWorldServiceId(service.Info().ID)
	}
	_newHelloWorldServiceServer(service, server, options, so.interceptors...)

	return service
}

func _newHelloWorldServiceServer(service micro.Service, server HelloWorldServiceNATSServer, opts *impl.ServerOpts, interceptors ...ServerInterceptor) {
	var err error
	// Define the NATS Micro Handler
	HelloWorldHandler := micro.HandlerFunc(func(request micro.Request) {
		// 1. Create the Context
		// We inherit from Background, but we could allow a base context to be passed in options
		ctx := context.Background()

		// 2. Extract Headers and put them into Context (Crucial for Tracing extraction!)
		// The request.Headers() contains the SpanID sent by client
		ctx = NewContextWithHeaders(ctx, nats_go.Header(request.Headers()))

		// 3. Unmarshal (We must do this before calling interceptors so they see the object)
		var req HelloWorldRequest
		if err := proto.Unmarshal(request.Data(), &req); err != nil {
			request.Error("560", "Failed to unmarshal proto message", []byte(err.Error()))
			return
		}

		// 4. Define the Method Info
		info := &MethodInfo{
			Subject: request.Subject(),
			Service: "HelloWorldService",
			Method:  "HelloWorld",
		}

		// 5. Define the "Final Handler"
		// This bridges the generic Interceptor world back to your specific Typed Implementation
		handler := func(ctx context.Context, req proto.Message) (proto.Message, error) {
			// CASTING: This is safe because we know 'req' is *HelloWorldRequest
			typedReq, ok := req.(*HelloWorldRequest)
			if !ok {
				return nil, fmt.Errorf("invalid request type: %T", req)
			}
			// Call the user's implementation
			return server.HelloWorld(ctx, typedReq)
		}

		// 6. Chain Interceptors
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := chain
			chain = func(currentCtx context.Context, currentReq proto.Message) (proto.Message, error) {
				return interceptor(currentCtx, currentReq, info, next)
			}
		}

		// 7. Execute
		resp, err := chain(ctx, &req)

		// 8. Handle Response / Error
		if err != nil {
			if protonats.IsServiceError(err) {
				slog.Warn("Server implementations should not return ServiceError, use go_nats.NewServerError instead", "error", err)
			}
			var serverErr protonats.ServerError
			if errors.As(err, &serverErr) {
				request.Error(serverErr.Code, serverErr.Description, serverErr.GetWrapped(), serverErr.GetOptHeaders())
			} else {
				request.Error("500", "Internal server error", []byte(err.Error()))
			}
			return
		}

		// Marshal Response
		data, err := proto.Marshal(resp)
		if err != nil {
			request.Error("560", "Failed to marshal response", []byte(err.Error()))
			return
		}
		request.Respond(data)
	})

	err = service.AddEndpoint("HelloWorld", HelloWorldHandler, opts.Subject("service.HelloWorldService.HelloWorld", ""))
	if err != nil {
		panic(err) // TODO: Update this to proper error handling
	}
	err = service.AddEndpoint("HelloWorld-Direct", HelloWorldHandler, opts.Subject("service.HelloWorldService.HelloWorld", service.Info().ID))
	if err != nil {
		panic(err) // TODO: Update this to proper error handling
	}
}

//endregion
