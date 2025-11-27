package main

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"xiam.li/go-protonats/internal/helloworld"
	shared "xiam.li/go-protonats/internal/helloworld/test"
)

type helloWorldServiceImpl struct {
}

func (h *helloWorldServiceImpl) HelloWorld(ctx context.Context, req *helloworld.HelloWorldRequest) (*helloworld.HelloWorldResponse, error) {
	value := ctx.Value("traceID")
	if value != nil {
		if traceID, ok := value.(string); ok {
			fmt.Printf("helloWorldServiceImpl: traceID=%s", traceID)
		} else {
			fmt.Println("helloWorldServiceImpl: traceID is not a string")
		}
	} else {
		fmt.Println("helloWorldServiceImpl: traceID not found in context")
	}
	message := "Hello, " + req.Name + "!"
	return &helloworld.HelloWorldResponse{Message: message}, nil
}

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_ = helloworld.NewHelloWorldServiceNATSServer(conn, &helloWorldServiceImpl{}, helloworld.WithServerInterceptors(shared.OTelServerInterceptor))
	fmt.Println("Starting HelloWorldServiceNATSServer...")
	for {
	}
}
