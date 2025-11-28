package main

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"xiam.li/go-protonats/internal/helloworld"
	shared "xiam.li/go-protonats/internal/helloworld/test"
)

func main() {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := helloworld.NewHelloWorldServiceNATSClient(conn, shared.TimingClientInterceptor, shared.OTelClientInterceptor, shared.TimingClientInterceptor)
	ctx := context.WithValue(context.Background(), "traceID", "12345")
	resp, err := client.HelloWorld(ctx, &helloworld.HelloWorldRequest{Name: "Client Test"}) // , protonats.WithInstanceID("SHOErVxlCWk93k7EGbjuKz"))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", resp.Message)
}
