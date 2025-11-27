package helloworld

import (
	context "context"
	fmt "fmt"
	"time"

	nats_go "github.com/nats-io/nats.go"
	micro "github.com/nats-io/nats.go/micro"
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
	// The subject is hardcoded or generated
	const subject = "service.HelloWorldService.HelloWorld"
	const method = "HelloWorld"

	resp := &HelloWorldResponse{}

	// Define the "Final Invoker". This is the function that actually calls NATS.
	// It is the last link in the chain.
	invoker := func(ctx context.Context, method string, req, reply proto.Message, opts ...protonats.CallOption) error {
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
			Subject: subject,
			Data:    data,
			Header:  nats_go.Header{},
		}

		// IMPORTANT: Inject Headers from Context (e.g. set by OpenTelemetry)
		if outgoingHeaders := HeadersFromOutgoingContext(ctx); outgoingHeaders != nil {
			for k, v := range outgoingHeaders {
				for _, val := range v {
					msg.Header.Add(k, val)
				}
			}
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
		chain = func(currentCtx context.Context, currentMethod string, currentReq, currentReply proto.Message, currentOpts ...protonats.CallOption) error {
			return interceptor(currentCtx, currentMethod, currentReq, currentReply, c.nc, next, currentOpts...)
		}
	}

	// Execute the chain
	err := chain(ctx, method, req, resp, opts...)
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

func NewHelloWorldServiceNATSServer(nc *nats_go.Conn, server HelloWorldServiceNATSServer, opts ...ServerOption) micro.Service {
	// Parse options
	so := &serverOptions{}
	for _, opt := range opts {
		opt(so)
	}

	// (Existing impl initialization logic omitted for brevity, assuming impl.NewService handles basic setup)
	service, _, err := impl.NewService("HelloWorldService", nc, server)
	if err != nil {
		panic(err)
	}

	_newHelloWorldServiceServer(service, server, so.interceptors)
	return service
}

func _newHelloWorldServiceServer(service micro.Service, server HelloWorldServiceNATSServer, interceptors []ServerInterceptor) {

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
			Subject:    request.Subject(),
			FullMethod: "service.HelloWorldService.HelloWorld",
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
			// Check if it's a specific ServiceError, etc. (Existing logic)
			if protonats.IsServiceError(err) {
				// Log warning...
			}
			// Respond with error
			request.Error("500", "Internal Server Error", []byte(err.Error()))
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

	// Add Endpoint
	// Note: You might want to pass middlewares for validation here too, but that's separate
	service.AddEndpoint("HelloWorld", HelloWorldHandler)
}

//endregion
