// Package main is the entry point for the server application in the monitoring system.
// It sets up and runs the GRPC server, handling configuration, logging, and graceful shutdowns.
// The package integrates various internal components.
package main

import (
	"log"
	"net"

	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/startup"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// main is the entry function of the application. It sets up and starts the GRPC server,
// including configuration, logging, storage. It also manages the application's lifecycle,
// handling initialization and graceful shutdown.
func main() {
	utils.PrintBuildFlags(buildVersion, buildDate, buildCommit)

	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}

	done := make(chan struct{}, 1)
	tickerChan := startup.SetupTicker(conf)
	retrier := startup.SetupRetrier(conf, logger)

	keeper := startup.SetupKeeper(conf, logger, retrier)
	if keeper != nil {
		defer keeper.Shutdown()
	}
	storage := adapters.NewMemStorage(done, tickerChan, logger, keeper)
	startup.PrepareStorage(conf, keeper, storage)
	key, err := conf.ReadPrivateKey()
	if err != nil {
		panic(err)
	}

	srv := startup.SetupGRPCServer(logger, storage, conf.HMACKey, conf.TrustedSubnet, key)

	go func() {
		logger.Infof("Launching the server at %s\n", conf.Addr)
		lis, err := net.Listen("tcp", conf.Addr)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := srv.Serve(lis); err != nil { // TODO: !errors.Is(err, http.ErrServerClosed) ?
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	startup.GracefulShutdownGRPC(srv, done, logger)
}
