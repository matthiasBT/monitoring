package adapters

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
		err    error
		meta   []string
		cipher []byte
		hash   string
	)

	encrypted := r.CryptoKey != nil
	if encrypted {
		cipher, err = r.encrypt(payload)
		if err != nil {
			return err
		}
		if hash, err = utils.HashData(cipher, r.HMACKey); err != nil {
			return err
		}
	} else if hash, err = r.getHMACHeader(payload); err != nil {
		return err
	}

	if addr, err := getLocalIP(); err != nil {
		return err
	} else {
		meta = append(meta, "X-Real-IP", addr)
	}
	if hash != "" {
		meta = append(meta, "HashSHA256", hash)
	}

	f := func() (any, error) {
		var (
			res any
			err error
		)
		conn, err := grpc.Dial(r.ServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		c := pb.NewMonitoringClient(conn)
		md := metadata.Pairs(meta...)
		ctx := metadata.NewOutgoingContext(context.Background(), md)
		if encrypted {
			req := new(pb.EncryptedMetricsArray)
			req.Metrics = cipher
			res, err = c.MassUpdateMetricsEncrypted(ctx, req)
		} else {
			req := utils.HTTPMultipleMetricsToGRPC(payload)
			res, err = c.MassUpdateMetrics(ctx, req)
		}
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	ctx := context.Background()
	_, err = r.Retrier.RetryChecked(ctx, f, utils.CheckConnectionError)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		r.Logger.Infof("Response metadata: %v", md)
	}
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			fmt.Println("Error was not a gRPC status error")
		} else {
			r.Logger.Infof("Error code: %v", st.Code())
			r.Logger.Infof("Error message: %s", st.Message())

		}
		return err
	} else {
		r.Logger.Info("Success")
	}
	return nil
}

func (r *GRPCReportAdapter) getHMACHeader(payload []*common.Metrics) (string, error) {
	binary, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return utils.HashData(binary, r.HMACKey)
}

func (r *GRPCReportAdapter) encrypt(payload []*common.Metrics) ([]byte, error) {
	var raw []byte
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	cipher, err := encryptData(raw, r.CryptoKey)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}