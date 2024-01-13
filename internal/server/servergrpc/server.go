package servergrpc

import (
	"context"
	"encoding/json"
	"errors"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
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

func NewServer(logger logging.ILogger, storage entities.Storage) *Server {
	return &Server{
		Logger:  logger,
		Storage: storage,
	}
}

func (s *Server) Ping(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	if err := s.Storage.Ping(ctx); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return new(pb.Empty), nil
}

func (s *Server) GetMetric(ctx context.Context, req *pb.Metrics) (*pb.Metrics, error) {
	metrics := utils.GRPCMetricToHTTP(req)
	if err := metrics.Validate(false); err != nil {
		return nil, utils.WrapInvalidMetricError(err)
	}
	result, err := s.Storage.Get(ctx, metrics)
	if err != nil {
		var code codes.Code
		if errors.Is(err, common.ErrUnknownMetric) {
			code = codes.NotFound
		} else {
			code = codes.Internal
		}
		return nil, status.Errorf(code, err.Error())
	}
	return utils.HTTPMetricToGRPC(result), nil
}

func (s *Server) UpdateMetric(ctx context.Context, req *pb.Metrics) (*pb.Metrics, error) {
	metrics := utils.GRPCMetricToHTTP(req)
	if err := metrics.Validate(true); err != nil {
		return nil, utils.WrapInvalidMetricError(err)
	}
	result, err := s.Storage.Add(ctx, metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return utils.HTTPMetricToGRPC(result), nil
}

func (s *Server) MassUpdateMetrics(ctx context.Context, req *pb.MetricsArray) (*pb.Empty, error) {
	batch := utils.GRPCMultipleMetricsToHTTP(req)
	return s.massUpdate(ctx, batch)
}

func (s *Server) MassUpdateMetricsEncrypted(ctx context.Context, req *pb.EncryptedMetricsArray) (*pb.Empty, error) {
	var batch []*common.Metrics
	if err := json.Unmarshal(req.Metrics, &batch); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return s.massUpdate(ctx, batch)
}

func (s *Server) GetAllMetrics(ctx context.Context, req *pb.Empty) (*pb.MetricsArray, error) {
	var batch map[string]*common.Metrics
	var err error
	batch, err = s.Storage.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var values []*common.Metrics
	for _, v := range batch {
		values = append(values, v)
	}
	return utils.HTTPMultipleMetricsToGRPC(values), nil
}

func (s *Server) massUpdate(ctx context.Context, batch []*common.Metrics) (*pb.Empty, error) {
	for _, metrics := range batch {
		if err := metrics.Validate(true); err != nil {
			return nil, utils.WrapInvalidMetricError(err)
		}
	}
	if err := s.Storage.AddBatch(ctx, batch); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return new(pb.Empty), nil
}
