package go_nats

import "github.com/nats-io/nats.go/micro"

// Those handler interfaces can be implemented by the respective NATS server impl,
//as an alternative to setting a handler via an option when creating the server.

// StatsHandler is an interface that when implemented on the server, will
// be used directly instead of having to use the WithStatsHandler option.
type StatsHandler interface {
	Stats(endpoint *micro.Endpoint) any
}

// DoneHandler is an interface that when implemented on the server, will
// be used directly instead of having to use the WithDoneHandler option.
type DoneHandler interface {
	Done(service micro.Service)
}

// ErrHandler is an interface that when implemented on the server, will
// be used directly instead of having to use the WithErrorHandler option.
type ErrHandler interface {
	Err(service micro.Service, natsErr *micro.NATSError)
}
