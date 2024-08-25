package protoc_plugin

import "github.com/nats-io/nats.go/micro"

type StatsHandler interface {
	Stats(endpoint *micro.Endpoint) any
}

type DoneHandler interface {
	Done(service micro.Service)
}

type ErrHandler interface {
	Err(service micro.Service, natsErr *micro.NATSError)
}
