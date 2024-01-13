package servergrpc

import (
	"context"
	"log"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type LoggingInterceptor struct {
	Logger logging.ILogger
}

func (l LoggingInterceptor) Interceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	l.Logger.Infof("Served: %s %v\n", info.FullMethod, duration)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("Error response Code: %v", st.Code())
	}
	return resp, err
}
