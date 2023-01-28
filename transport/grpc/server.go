package grpc

import (
	"context"
	"net"
	"net/url"
	"time"

	"github.com/elioli1991/app-infra/abstract"
	logger "github.com/elioli1991/app-logger"

	"github.com/elioli1991/app-runner/internal/host"
	"github.com/elioli1991/app-runner/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	defaultNetwork = "tcp"
	defaultTimeOut = 1 * time.Second
	defaultAddress = ":0"
)

var (
	_ transport.Service    = (*Server)(nil)
	_ transport.EndPointer = (*Server)(nil)
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

func Options(opts ...grpc.ServerOption) ServerOption {
	return func(o *Server) {
		o.grpcOpts = opts
	}
}

// OpenHealth use grpc default health service
func OpenHealth() ServerOption {
	return func(o *Server) {
		o.openHealth = true
	}
}

// Logger set logger
func Logger(l abstract.Logger) ServerOption {
	return func(o *Server) {
		o.logger = l
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(o *Server) {
		o.lis = lis
	}
}

// NewServer creates a new grpc server by server options
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: defaultNetwork,
		address: defaultAddress,
		timeout: defaultTimeOut,
		health:  health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}
	if srv.logger == nil {
		srv.logger = logger.GetLogger()
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(srv.unaryInterceptor...),
		grpc.ChainStreamInterceptor(srv.streamInterceptor...),
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	if srv.openHealth {
		grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	}
	reflection.Register(srv.Server)
	return srv
}

type Server struct {
	*grpc.Server
	ctx               context.Context
	lis               net.Listener
	network           string
	address           string
	timeout           time.Duration
	logger            abstract.Logger
	endpoint          *url.URL
	health            *health.Server
	openHealth        bool
	unaryInterceptor  []grpc.UnaryServerInterceptor
	streamInterceptor []grpc.StreamServerInterceptor
	grpcOpts          []grpc.ServerOption
	adminClean        func()
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.listenAndEndpoint(); err != nil {
		return err
	}
	s.ctx = ctx
	s.logger.Infof("[gRpc] server listening on %s", s.lis.Addr().String())
	s.health.Resume()
	return s.Serve(s.lis)
}

func (s *Server) Stop(ctx context.Context) error {
	if s.adminClean != nil {
		s.adminClean()
	}
	s.health.Shutdown()
	s.GracefulStop()
	s.logger.Info("[gRPC] server stopping")
	return nil
}

// EndPoint return a real address
func (s *Server) EndPoint() (*url.URL, error) {
	if err := s.listenAndEndpoint(); err != nil {
		return nil, err
	}
	return s.endpoint, nil
}

// listenAndEndpoint
func (s *Server) listenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	if s.endpoint == nil {
		addr, err := host.Extract(s.address, s.lis)
		if err != nil {
			return err
		}
		s.endpoint = &url.URL{Scheme: "grpc", Host: addr}
	}
	return nil
}
