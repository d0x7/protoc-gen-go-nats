# protoc-gen-go-nats

This is a protoc plugin that generates Go server and client code for NATS microservices.

Prior experience with Protobuf is greatly recommended, especially to understand how the package and imports work.

## Installation

You already need to have the protoc compiler along the Go protobuf plugin installed on your system.
After that, you can go ahead and install this plugin using the following command:

```shell
go install xiam.li/go-nats/cmd/protoc-gen-go-nats@latest
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

### Custom Errors

You can also send custom errors to the client, but for that you need to add this package to your project:

```shell
go get xiam.li/go-nats
```

Then, you can use the `go_nats.ServerError` type to send custom errors to the client:

```go
// In any method of your service implementation, do the following
// Or, if you want to return a custom error:
return nil, go_nats.NewServerErr("400", "Unknown Name")

// Or, you can also wrap an existing error for more detailed information:
return nil, go_nats.WrapServerErr(err, "500", "Failed to query database")

// You can also send custom headers using this method:
serverErr := go_nats.NewServerErr("400", "Unknown Name")
serverErr.AddHeader("err-details", "Username is not in the database")
return nil, serverErr
```

### Streaming

Streaming is not yet supported, but is planned for the future.
It'll probably be implemented along with better timeout handling,
that will come with keepalive messages and therefore also allow streaming.
