package test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	natstest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
)

type NatsInstance struct {
	server *server.Server
	*nats.Conn
}

func newNATS(t *testing.T) *NatsInstance {
	port, err := GetFreePort()
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	opts := natstest.DefaultTestOptions
	opts.Port = port
	testServer := natstest.RunServer(&opts)
	if testServer.ReadyForConnections(1*time.Second) != true {
		t.Fatalf("Failed to start NATS server")
	}
	testClient, err := nats.Connect(fmt.Sprintf("nats://127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	if !testClient.HeadersSupported() {
		t.Fatalf("Headers not supported; please upgrade to a newer version of NATS.")
	}
	return &NatsInstance{
		server: testServer,
		Conn:   testClient,
	}
}

func (s *NatsInstance) Stop() {
	s.server.Shutdown()
	s.Conn.Close()
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer func(l *net.TCPListener) {
				_ = l.Close()
			}(l)
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
