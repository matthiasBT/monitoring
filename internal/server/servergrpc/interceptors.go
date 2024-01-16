package servergrpc

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	Logger    logging.ILogger
	HMACKey   []byte
	Subnet    *net.IPNet
	CryptoKey *rsa.PrivateKey
}

func (i Interceptor) LoggingInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	i.Logger.Infof("Served: %s %v\n", info.FullMethod, duration)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Error response Code: %v. Message: %v", st.Code(), err.Error())
	}
	return resp, err
}

func (i Interceptor) HashCheckInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-Data-Hash")
		if len(values) > 0 {
			payload, err := toBinary(req)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			serverHash, err := utils.HashData(payload, i.HMACKey)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			if serverHash != values[0] {
				return nil, status.Errorf(codes.InvalidArgument, "payload hash mismatch")
			}
		}
	}
	return handler(ctx, req)
}

func (i Interceptor) HashWriteInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	payload, err := toBinary(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if hash, err := utils.HashData(payload, i.HMACKey); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	} else {
		md := metadata.Pairs("X-Data-Hash", hash)
		if err := grpc.SetTrailer(ctx, md); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return resp, nil
}

func (i Interceptor) SubnetCheckInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("X-Real-IP")
		if len(values) > 0 {
			var clientIP = net.ParseIP(values[0])
			if clientIP == nil {
				return nil, status.Errorf(codes.PermissionDenied, "Invalid X-Real-IP header value")
			}
			if !i.Subnet.Contains(clientIP) {
				return nil, status.Errorf(codes.PermissionDenied, "Untrusted client IP address")
			}
		} else {
			return nil, status.Errorf(codes.PermissionDenied, "Missing X-Real-IP header value")
		}
	}
	return handler(ctx, req)
}

func (i Interceptor) DecryptionInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if data, ok := req.(*pb.MetricsArray); ok && data.Metrics != nil {
		if plaintext, err := utils.Decrypt(data.Metrics, i.CryptoKey); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		} else {
			data.Metrics = plaintext
		}
	}
	return handler(ctx, req)
}

func toBinary(data any) ([]byte, error) {
	var (
		payload []byte
		err     error
	)
	if reqArr, ok := data.(*pb.MetricsArray); ok {
		if reqArr.Metrics == nil {
			metricsMultiple := utils.GRPCMultipleMetricsToHTTP(reqArr)
			payload, err = json.Marshal(metricsMultiple)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		} else {
			payload = reqArr.Metrics
		}
	} else if reqSingle, ok := data.(*pb.Metrics); ok {
		metricsSingle := utils.GRPCMetricToHTTP(reqSingle)
		payload, err = json.Marshal(metricsSingle)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	return payload, nil
}
