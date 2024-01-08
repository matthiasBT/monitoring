package servergrpc

import (
	"context"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedMonitoringServer
	Logger  logging.ILogger
	Storage entities.Storage
}

func (s *Server) Ping(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	if err := s.Storage.Ping(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return new(pb.Empty), nil
}

func NewServer(logger logging.ILogger, storage entities.Storage) *Server {
	return &Server{
		Logger:  logger,
		Storage: storage,
	}
}
