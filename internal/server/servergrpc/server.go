package servergrpc

import (
	"context"
	"errors"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	metrics := unwrapMetrics(req)
	if err := metrics.Validate(false); err != nil {
		return nil, wrapInvalidMetricError(err)
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
	return wrapMetrics(result), nil
}

func (s *Server) UpdateMetric(ctx context.Context, req *pb.Metrics) (*pb.Metrics, error) {
	metrics := unwrapMetrics(req)
	if err := metrics.Validate(true); err != nil {
		return nil, wrapInvalidMetricError(err)
	}
	result, err := s.Storage.Add(ctx, metrics)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return wrapMetrics(result), nil
}

func (s *Server) MassUpdateMetrics(ctx context.Context, req *pb.MetricsArray) (*pb.Empty, error) {
	var batch []*common.Metrics
	for _, wrapped := range req.Objects {
		unwrapped := unwrapMetrics(wrapped)
		if err := unwrapped.Validate(true); err != nil {
			return nil, wrapInvalidMetricError(err)
		}
		batch = append(batch, unwrapMetrics(wrapped))
	}
	if err := s.Storage.AddBatch(ctx, batch); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return new(pb.Empty), nil
}

func (s *Server) GetAllMetrics(ctx context.Context, req *pb.Empty) (*pb.MetricsArray, error) {
	var batch map[string]*common.Metrics
	var err error
	batch, err = s.Storage.GetAll(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	var result []*pb.Metrics
	for _, unwrapped := range batch {
		wrapped := wrapMetrics(unwrapped)
		result = append(result, wrapped)
	}
	arr := new(pb.MetricsArray)
	arr.Objects = result
	return arr, nil
}

func wrapInvalidMetricError(err error) error {
	var code codes.Code
	switch {
	case errors.Is(err, common.ErrInvalidMetricType) || errors.Is(err, common.ErrInvalidMetricVal):
		code = codes.InvalidArgument
	case errors.Is(err, common.ErrMissingMetricName):
		code = codes.NotFound
	default:
		code = codes.Internal
	}
	return status.Errorf(code, err.Error())
}

func unwrapMetrics(req *pb.Metrics) *common.Metrics {
	metrics := new(common.Metrics)
	metrics.ID = req.Id
	metrics.MType = req.MType
	if req.Value != nil {
		metrics.Value = &req.Value.Value
	}
	if req.Delta != nil {
		metrics.Delta = &req.Delta.Value
	}
	return metrics
}

func wrapMetrics(result *common.Metrics) *pb.Metrics {
	resp := new(pb.Metrics)
	resp.Id = result.ID
	resp.MType = result.MType
	if result.Delta != nil {
		resp.Delta = &wrapperspb.Int64Value{Value: *result.Delta}
	}
	if result.Value != nil {
		resp.Value = &wrapperspb.DoubleValue{Value: *result.Value}
	}
	return resp
}
