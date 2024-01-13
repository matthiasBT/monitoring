package adapters

import (
	"context"
	"crypto/rsa"
	"encoding/json"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCReportAdapter struct {
	Logger     logging.ILogger
	Jobs       chan []*common.Metrics
	ServerAddr string
	HMACKey    []byte
	CryptoKey  *rsa.PublicKey
	Retrier    utils.Retrier
}

func NewGRPCReportAdapter(
	logger logging.ILogger,
	serverAddr string,
	retrier utils.Retrier,
	hmacKey []byte,
	cryptoKey *rsa.PublicKey,
	workerNum uint,
) *GRPCReportAdapter {
	adapter := GRPCReportAdapter{
		Logger:     logger,
		ServerAddr: serverAddr,
		Retrier:    retrier,
		HMACKey:    hmacKey,
		CryptoKey:  cryptoKey,
		Jobs:       make(chan []*common.Metrics, workerNum),
	}
	var i uint
	for i = 0; i < workerNum; i++ {
		go func() {
			for {
				data := <-adapter.Jobs
				if err := adapter.report(data); err != nil {
					logger.Errorf("Failed to report: %v", err)
				}
			}
		}()
	}
	return &adapter
}

func (r *GRPCReportAdapter) ReportBatch(batch []*common.Metrics) error {
	r.Jobs <- batch
	return nil
}

func (r *GRPCReportAdapter) report(payload []*common.Metrics) error {
	var (
		err  error
		meta []string
	)
	//if r.CryptoKey != nil {
	//	payload, err = encryptData(payload, r.CryptoKey)
	//	if err != nil {
	//		return err
	//	}
	//} // TODO: add a separate function

	if addr, err := getLocalIP(); err != nil {
		return err
	} else {
		meta = append(meta, "X-Real-IP", addr)
	}
	if hash, err := r.getHMACHeader(payload); err != nil {
		return err
	} else if hash != "" {
		meta = append(meta, "HashSHA256", hash)
	}

	req := utils.HTTPMultipleMetricsToGRPC(payload)

	f := func() (any, error) {
		conn, err := grpc.Dial(r.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		c := pb.NewMonitoringClient(conn)
		md := metadata.Pairs(meta...)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		res, err := c.MassUpdateMetrics(ctx, req)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	_, err = r.Retrier.RetryChecked(context.Background(), f, utils.CheckConnectionError)
	if err != nil {
		return err
	}
	r.Logger.Info("Success")
	return nil
}

func (r *GRPCReportAdapter) getHMACHeader(payload []*common.Metrics) (string, error) {
	binary, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	if hash, err := hashData(binary, r.HMACKey); err != nil {
		return "", err
	} else if hash != "" {
		return hash, nil
	}
	return "", nil
}
