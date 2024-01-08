package startup

import (
	"context"
	"crypto/rsa"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/compression"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/secure"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/entities"
	"github.com/matthiasBT/monitoring/internal/server/servergrpc"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
	pb "github.com/matthiasBT/monitoring/proto"
	"google.golang.org/grpc"
)

// SetupHTTPServer configures and returns a new HTTP router with middleware and routes.
// It includes logging, compression, optional HMAC checking, and controller routes.
func SetupHTTPServer(
	logger logging.ILogger, controller *usecases.BaseController, hmacKey, subnet string, key *rsa.PrivateKey,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logging.Middleware(logger))
	r.Use(compression.MiddlewareReader, compression.MiddlewareWriter)
	if hmacKey != "" {
		r.Use(secure.MiddlewareHashReader(hmacKey), secure.MiddlewareHashWriter(hmacKey))
	}
	if key != nil {
		r.Use(secure.MiddlewareCryptoReader(key))
	}
	if subnet != "" {
		r.Use(secure.MiddlewareIPFilter(logger, subnet))
	}
	r.Mount("/", controller.Route())
	return r
}

func SetupGRPCServer(
	logger logging.ILogger, storage entities.Storage, hmacKey, subnet string, key *rsa.PrivateKey,
) *grpc.Server {
	srv := servergrpc.NewServer(logger, storage)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(servergrpc.LoggingInterceptor),
	)
	pb.RegisterMonitoringServer(s, srv)
	return s
}

// SetupKeeper initializes and returns the appropriate Keeper (database or file) based on configuration.
// It configures the storage mechanism for the server, handling data persistence.
func SetupKeeper(conf *server.Config, logger logging.ILogger, retrier utils.Retrier) entities.Keeper {
	if conf.Flushes() {
		if conf.DatabaseDSN != "" {
			db := adapters.OpenDB(conf.DatabaseDSN)
			return adapters.NewDBKeeper(db, logger, retrier)
		} else {
			return adapters.NewFileKeeper(conf, logger, retrier)
		}
	}
	return nil
}

// SetupRetrier configures and returns a Retrier based on the provided server configuration.
// It sets up retry attempts, intervals, and logging for handling network-related retries.
func SetupRetrier(conf *server.Config, logger logging.ILogger) utils.Retrier {
	return utils.Retrier{
		Attempts:         conf.RetryAttempts,
		IntervalFirst:    conf.RetryIntervalInitial,
		IntervalIncrease: conf.RetryIntervalBackoff,
		Logger:           logger,
	}
}

// PrepareStorage loads the initial state of the storage and launches periodic storage data flushing if necessary
func PrepareStorage(conf *server.Config, keeper entities.Keeper, storage entities.Storage) {
	if conf.Flushes() {
		if conf.Restore {
			state := keeper.Restore()
			storage.Init(state)
		}
		if !conf.FlushesSync() {
			go storage.FlushPeriodic(context.Background())
		}
	}
}

// SetupTicker creates and returns a ticker channel based on the configuration.
// It's used for periodic operations like data flushing.
func SetupTicker(conf *server.Config) <-chan time.Time {
	if conf.FlushesSync() {
		return make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(conf.StoreInterval) * time.Second)
		return ticker.C
	}
}

// GracefulShutdownHTTP handles the graceful shutdown of the HTTP server.
// It listens for system signals and shuts down the server after processing ongoing requests.
func GracefulShutdownHTTP(srv *http.Server, done chan struct{}, logger logging.ILogger) {
	gracefulShutdown(done, logger)
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err.Error())
	}
}

// GracefulShutdownGRPC handles the graceful shutdown of the HTTP server.
// It listens for system signals and shuts down the server after processing ongoing requests.
func GracefulShutdownGRPC(srv *grpc.Server, done chan struct{}, logger logging.ILogger) {
	gracefulShutdown(done, logger)
	srv.GracefulStop()
}

func gracefulShutdown(done chan struct{}, logger logging.ILogger) {
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-quitChannel
	logger.Infof("Received signal: %v\n", sig)
	done <- struct{}{}
	time.Sleep(5 * time.Second)
}
