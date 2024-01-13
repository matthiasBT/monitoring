package utils

import (
	"errors"

	"github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/proto"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func HTTPMetricToGRPC(result *entities.Metrics) *proto.Metrics {
	resp := new(proto.Metrics)
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

func GRPCMetricToHTTP(req *pb.Metrics) *entities.Metrics {
	metrics := new(entities.Metrics)
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

func GRPCMultipleMetricsToHTTP(grpcMetrics *pb.MetricsArray) []*entities.Metrics {
	var batch []*entities.Metrics
	for _, wrapped := range grpcMetrics.Objects {
		unwrapped := GRPCMetricToHTTP(wrapped)
		batch = append(batch, unwrapped)
	}
	return batch
}

func HTTPMultipleMetricsToGRPC(httpMetrics []*entities.Metrics) *pb.MetricsArray {
	var result []*pb.Metrics
	for _, unwrapped := range httpMetrics {
		wrapped := HTTPMetricToGRPC(unwrapped)
		result = append(result, wrapped)
	}
	arr := new(pb.MetricsArray)
	arr.Objects = result
	return arr
}

func WrapInvalidMetricError(err error) error {
	var code codes.Code
	switch {
	case errors.Is(err, entities.ErrInvalidMetricType) || errors.Is(err, entities.ErrInvalidMetricVal):
		code = codes.InvalidArgument
	case errors.Is(err, entities.ErrMissingMetricName):
		code = codes.NotFound
	default:
		code = codes.Internal
	}
	return status.Errorf(code, err.Error())
}
