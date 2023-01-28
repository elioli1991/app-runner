package grpc

import (
	"context"
	"testing"

	logger "github.com/elioli1991/app-logger"
	pb "github.com/elioli1991/app-runner/internal/testdata"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	logger.Infof("grpc request name %v", request.GetName())
	return &pb.HelloReply{Message: "ok"}, nil
}

func TestServer_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			pb.RegisterGreeterServer(s, &server{})
			if err := s.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
