package servergrpc

import (
	"context"
	"log"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TODO: implement other interceptors

func LoggingInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	logger := logging.SetupLogger() // TODO: find a way to pass the logger (via a closure?)
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	logger.Infof("Served: %s %v\n", info.FullMethod, duration)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Response Code: %v", st.Code())
	} else {
		logger.Infof("Response: %d bytes", 0) // TODO: implement
	}
	return resp, err
}

func GzipDecompressorInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val := md["content-encoding"]; len(val) > 0 && val[0] == "gzip" {
			newCtx := context.WithValue(ctx, grpc.UseCompressor(gzip.Name), true)
			return handler(newCtx, req)
		}
	}
	return handler(ctx, req)
}
