# go-protonats [![PkgGoDev](https://pkg.go.dev/badge/xiam.li/go-protonats)](https://pkg.go.dev/xiam.li/go-protonats)

This is a protoc plugin that generates Go server and client code for NATS microservices.

Go-ProtoNATS uses the shared [protonats](https://github.com/d0x7/protonats) package and is protocol compatible with the [Java implementation](https://github.com/d0x7/java-protonats).

Prior experience with Protobuf is greatly recommended, especially to understand how the package and imports work.

> [!IMPORTANT]
> Please note that the structure has greatly been changed since the last release.
> This repository now only contains the protoc-gen-go-nats compiler plugin, but not the runtime code.
> For that, please refer to [protonats](https://github.com/d0x7/protonats) and import that module instead.

## Installation

You already need to have the protoc compiler along the Go protobuf plugin installed on your system.
After that, you can go ahead and install this plugin using the following command:

```shell
go install xiam.li/go-protonats/cmd/protoc-gen-go-nats@latest
```

To check if the installation was successful, you can run:

```shell
protoc-gen-go-nats -v
```

## Usage

Upon installation, you should create a protobuf file that contains a service, similar to how gRPC servers work.
An example protobuf file might look like this:

```protobuf
syntax = "proto3";
package your.package;
option go_package = "github.com/user/repo/pb;pb";

service HelloWorldService {
    rpc HelloWorld(HelloWorldRequest) returns (HelloWorldResponse);
}

message HelloWorldRequest {
    string name = 1;
}

message HelloWorldResponse {
    string message = 1;
}
```

To generate the Go code for this service, run the following command.
This command expects your proto file in a directory called `pb` in your project.

```shell
protoc -I pb --go_out=pb --go_opt=paths=source_relative --go-nats_out=pb --go-nats_opt=paths=source_relative pb/hello_world.proto
```

This obviously requires the protoc compiler to be installed on your system
and also having the go protobuf plugin installed, so that besides the code
regarding NATS can be generated, the messages and everything else can also be generated.

Now you can use the generated code to create a NATS server and client:

```go
package main
import (
    "fmt"
    "github.com/nats-io/nats.go"
    "github.com/user/repo/pb"
)

type serviceImpl struct {}

func (s *serviceImpl) HelloWorld(req *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	msg := fmt.Sprintf("Hello, %s!", req.GetName())
	return &pb.HelloWorldResponse{Message: msg}, nil
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	pb.NewHelloWorldServiceNATSServer(nc, &serviceImpl{})
}
```

Client:

```go
package main
import (
	"github.com/nats-io/nats.go"
    "github.com/user/repo/pb"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	cli := pb.NewHelloWorldServiceNATSClient(nc)

	// List all instances currently connected
	instances, err := cli.ListInstances()

	// Get stats from all instances
	stats, err := cli.Stats()

	// Or, from a specific instance:
	stats, err := cli.Stats(pb.WithInstanceID(instances[0].ID))

	// Obviously, you can also call your defined service methods:
	response, err := cli.HelloWorld(&pb.HelloWorldRequest{Name: "John Doe"})

	// And again for a specific instance, instead of the default load balanced distribution:
	response, err := cli.HelloWorld(&pb.HelloWorldRequest{Name: "John Doe"}, pb.WithInstanceID(instances[0].ID))
}
```

### Special handling for empty requests/responses

When specifying an RPC method that uses either or both the [`google/protobuf/empty.proto`](https://protobuf.dev/reference/protobuf/google.protobuf/#empty) type, that method will not generate a parameter to be passed as request/response, depending on how the RPC is defined.

For example, the following protobuf file will generate the following method signature:

```protobuf
syntax = "proto3";
package your.package;
import "google/protobuf/empty.proto";
option go_package = "github.com/user/repo/pb;pb";

service HelloWorldService {
    rpc NoRequest(google.protobuf.Empty) returns (HelloWorldResponse);
    rpc NoResponse(HelloWorldRequest) returns (google.protobuf.Empty);
    rpc NoRequestNoResponse(google.protobuf.Empty) returns (google.protobuf.Empty);
}
```

```go
// Client
type HelloWorldServiceNATSClient interface {
    NoRequest(opts ...CallOption) (*HelloWorldResponse, error)
    NoResponse(req *HelloWorldResponse, opts ...CallOption) (error)
    NoRequestNoResponse(opts ...CallOption) (error)
    // [...]
}
// Server
type HelloWorldServiceNATSServer interface {
    NoRequest() (*HelloWorldResponse, error)
    NoResponse(req *HelloWorldResponse) (error)
    NoRequestNoResponse() (error)
    // [...]
}
```

### Broadcasting

If you want to broadcast a message to all instances of a service, you can set the `protonats.broadcast` option to true in the method definition.
This will generate a method in the client that will broadcast the message to all instances of the service.

The issue with that is, it requires the client to have the `protonats.proto` imported, but can be easily done, by appending this to your protoc generation command:

```shell
-I$(go list -m -f '{{ .Dir }}' xiam.li/protonats)/proto
```

This takes the local directory of the `protonats` module and adds it as a import path for proto, so that it can find the `protonats.proto` file in there.
You can now use it like this:

```protobuf
syntax = "proto3";
package your.package;
import "google/protobuf/empty.proto";
import "protonats.proto";
option go_package = "github.com/user/repo/pb;pb";

service BroadcastingService {
  rpc ABroadcastingMethod(HelloWorldRequest) returns (HelloWorldResponse) {
    option (protonats.broadcast) = true;
  }
  rpc FanOut(HelloWorldRequest) returns (google.protobuf.Empty) {
    option (protonats.broadcast) = true;
  }
  rpc FanIn(google.protobuf.Empty) returns (HelloWorldResponse) {
    option (protonats.broadcast) = true;
  }
  rpc VeryEmptyMethod(google.protobuf.Empty) returns (google.protobuf.Empty) {
    option (protonats.broadcast) = true;
  }
}
```

When using broadcast, it will also honour your usage of `google.protobuf.Empty`, so that these methods won't generate a parameter to be passed as request/response, depending on how the RPC is defined.
Although there are opts available for these methods, the only one used is the timeout and passing an instance id does nothing.

You can use it like this on the server side:

```go
type impl struct {
	srv micro.Service
	err bool
}

func (i *impl) ABroadcastingMethod(req *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.Test{Message: "Hello " + req.GetName() + ", welcome to a broadcasting response from " + i.srv.Info().ID}, nil
}

func (i *impl) FanOut(req *pb.HelloWorldRequest) error {
	// Do something with the request and then return to confirm the request was received
	return nil
}

func (i *impl) FanIn() (*pb.HelloWorldResponse, error) {
	// Do something and then return a response
	return &pb.HelloWorldResponse{Message: "Hello FanIn response from " + i.srv.Info().ID}, nil
}

func (i *impl) VeryEmptyMethod() error {
	// Do something
	return nil
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	srvImpl := &impl{}
	srv := broadcast.NewBroadcastingServiceNATSServer(nc, srvImpl)
	srvImpl.srv = srv

	log.Println("Server started")
	for {
		if !srv.Stopped() {
			time.Sleep(time.Second)
		} else {
			return
		}
	}
}
```

And on the client:

```go
nc, err := nats.Connect(nats.DefaultURL)
if err != nil {
    log.Fatalf("Failed to connect to NATS: %v", err)
}
cli := broadcast.NewBroadcastingServiceNATSClient(nc)
responses, serviceErrs, err := cli.ABroadcastingMethod(&pb.HelloWorldRequest{Name: "John Doe"})
if err != nil {
    log.Fatalf("Failed to call ABroadcastingMethod: %v", err)
}
if serviceErrs != nil {
    log.Println("Some services returned an error:")
    for i, serviceErr := range serviceErrs {
        log.Printf("Service Error %d/%d: %v\n", i+1, len(serviceErrs), serviceErr)
    }
}
log.Printf("Broadcast Responses")
for i, response := range responses {
    log.Printf("Response %d/%d: %v\n", i+1, len(responses), response)
}
// Or similar-ish for methods using the empty type
serviceErrs, errs = cli.VeryEmptyMethod()
responses, serviceErrs, errs = cli.FanIn()
serviceErrs, errs = cli.FanOut(&pb.HelloWorldRequest{Name: "John Doe"})
```

### Instance identifier

If you need the instance id of the current service on the server side, you could either call `.Info().ID` on the returned `micro.Service`, but your service implementation can also just implement the service-specfic generated `[ServiceName]ServiceId` interface.
For example, if your service is defined as `HelloWorldService` in the protobuf file, the following interface would be generated:

```go
type HelloWorldServiceId interface {
	SetHelloWorldServiceId(string)
}
```

Therefore, you can just implement this interface on your service implementation optionally and upon initialization, the method would be called once, so you can store it in a field inside the service implementation, or similar:

```go
type helloWorldImpl struct {
	id string
}

func (i *helloWorldImpl) SetHelloWorldServiceId(s string) {
	i.id = s
}
```

### Consensus Integration

If you use a consensus algorithm like Raft, you can use the `protonats.consensus_Target` option to mark methods to be used only by the leader or follower.
These methods will be generated onto a separate interface, which is composited onto the main service interface.
By default, the normal `NewYourServiceNATSServer` method will still register all methods, regardless of it the target is leader or follower, but you can use the specialized `NewYourServiceNATSLeaderServer` or `NewYourServiceNATSFollowerServer` methods to only register a server for either methods - or you can use the normal `[...]NATSServer` method and pass either a `protonats.WithoutLeaderFns()` or `protonats.WithoutFollowerFns()` to disable the registration of these, but still allow for the
normal methods to be registered.

Methods marked with a consensus can still use the broadcasting flag, which will for example make a call to that method broadcast to all followers, instead of only one follower.

To mark methods with a consensus target, use the `protonats.consensus_target` option in the method definition:

```protobuf
service ConsensusService {
  // The CurrentSnapshot method will only be called on the leader
  rpc CurrentSnapshot(google.protobuf.Empty) returns (Snapshot) {
    option (protonats.consensus_target) = LEADER;
  }
  // And the ApplyChange will be sent to all followers, because it's also a broadcast
  rpc ApplyChange(Snapshot) returns (google.protobuf.Empty) {
    option (protonats.consensus_target) = FOLLOWER;
    option (protonats.broadcast) = true;
  }
}
```

```go
// Full implementation, normal methods, leader methods and follower methods
_ = consensus.NewConsensusServiceNATSServer(conn, impl)

// Normal implementation with follower methods, but leader methods unimplemented
_ = consensus.NewConsensusServiceNATSServer(conn, impl, protonats.WithoutLeaderFns())

// Normal implementation with leader methods, but follower methods unimplemented
_ = consensus.NewConsensusServiceNATSServer(conn, impl, protonats.WithoutFollowerFns())

// Follower-only implementation - only those marked as follower methods will be registered 
// Notice the use of the [...]NATSFollowerServer interface instead the broader [...]NATSServer interface
_ = consensus.NewConsensusServiceNATSFollowerServer(conn, impl)

// Leader-only implementation - only those marked as leader methods will be registered
// Notice the use of the [...]NATSLeaderServer interface instead the broader [...]NATSServer interface
_ = consensus.NewConsensusServiceNATSLeaderServer(conn, impl)
```

### Custom Errors

You can also send custom errors to the client, but for that you need to add this package to your project:

```shell
go get xiam.li/protonats
```

Then, you can use the `protonats.ServerError` type to send custom errors to the client:

```go
// In any method of your service implementation, do the following
// Or, if you want to return a custom error:
return nil, protonats.NewServerErr("400", "Unknown Name")

// Or, you can also wrap an existing error for more detailed information:
return nil, protonats.WrapServerErr(err, "500", "Failed to query database")

// You can also send custom headers using this method:
serverErr := protonats.NewServerErr("400", "Unknown Name")
serverErr.AddHeader("err-details", "Username is not in the database")
return nil, serverErr
```

On the client side they are received as `ServiceError` (Important: ServiceError, not ServerError).

```go
_, err := cli.HelloWorld(&pb.HelloWorldRequest{Name: "John Doe"})
if err != nil {
    serviceErr, isSrvErr := protonats.AsServiceError(err)
    if isSrvErr {
        fmt.Printf("Got a service error with code %s: %s\n", serviceErr.Code, serviceErr.Description)
    } else {
        fmt.Println("Other different error, usually networking related or an issue with unmarshalling the response")
    }
}
```

You can also use `protonats.IsServiceError(err)` to check if an error is a ServiceError.

There's also an `Details` field in the ServiceError struct, but that's only used when
the server, instead of returning a proper ServerError, only returns a generic error.
In that case, the result from that error's `Error()` will end up in the `Details` field.

### Streaming

Streaming is not yet supported, but is planned for the future.
It'll probably be implemented along with better timeout handling,
that will come with keepalive messages and therefore also allow streaming.
