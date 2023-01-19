package grpc

import (
	"context"
	"github.com/ElioLi-1991/app-runner/transport"
	"google.golang.org/grpc"
	"net"
	"net/url"
	"time"
)

const (
	defaultNetwork = "tcp"
	defaultTimeOut = 1 * time.Second
	defaultAddress = ":0"
)

var (
	_ transport.Service    = (*Server)()
	_ transport.EndPointer = (*Server)()
)

type ServerOption func(o *Server)

// NetWork Set NetWork config
func NetWork(network string) ServerOption {
	return func(o *Server) {
		o.network = network
	}
}

// Address Set Server address config
func Address(address string) ServerOption {
	return func(o *Server) {
		o.address = address
	}
}

// TimeOut set request timeout config
func TimeOut(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.timeout = timeout
	}
}

// NewServer creates a new grpc server
func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:     ctx,
		network: defaultNetwork,
		address: defaultAddress,
		timeout: defaultTimeOut,
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

type Server struct {
	*grpc.Server
	ctx     context.Context
	lis     net.Listener
	network string
	address string
	timeout time.Duration
}

func (s *Server) Start(ctx context.Context) error {
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) EndPoint() (*url.URL, error) {
	return nil, nil
}

func (s *Server) listenAndEndpoint() error {
	return nil
}
