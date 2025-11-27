package test

import (
	"testing"

	"github.com/nats-io/nats.go/micro"
	"xiam.li/protonats/go/impl"
)

type serviceImpl struct {
}

func (s *serviceImpl) Stats(*micro.Endpoint) any { return nil }

func (s *serviceImpl) Done(micro.Service) {}

func (s *serviceImpl) Err(micro.Service, *micro.NATSError) {}

func TestService(t *testing.T) {
	instance := newNATS(t)
	defer t.Cleanup(instance.Stop)
	service, opts, err := impl.NewService("test", instance.Conn, &serviceImpl{})
	if err != nil {
		t.Fatal(err)
	}
	if service == nil {
		t.Fatal("service is nil")
	}
	if opts == nil {
		t.Fatal("opts is nil")
	}
}
