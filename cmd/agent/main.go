// Package main is the entry point for the monitoring agent application.
// It initializes and orchestrates various components like logging, configuration,
// data reporting, and polling. The application handles periodic data collection
// and reporting to a central server.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matthiasBT/monitoring/internal/agent/adapters"
	"github.com/matthiasBT/monitoring/internal/agent/entities"
	"github.com/matthiasBT/monitoring/internal/agent/usecases/poll"
	"github.com/matthiasBT/monitoring/internal/agent/usecases/report"
	"github.com/matthiasBT/monitoring/internal/infra/config/agent"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
)

// setupRetrier configures and returns a Retrier based on the provided agent configuration.
// It sets up retry attempts, intervals, and logging for handling network-related retries.
func setupRetrier(conf *agent.Config, logger logging.ILogger) utils.Retrier {
	return utils.Retrier{
		Attempts:         conf.RetryAttempts,
		IntervalFirst:    conf.RetryIntervalInitial,
		IntervalIncrease: conf.RetryIntervalBackoff,
		Logger:           logger,
	}
}

// main is the entry function of the application. It sets up logging, configuration,
// data reporting, and polling mechanisms. It orchestrates the agent's lifecycle,
// including handling graceful shutdowns.
func main() {
	logger := logging.SetupLogger()
	conf, err := agent.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}
	done := make(<-chan bool)
	dataExchange := entities.SnapshotWrapper{CurrSnapshot: nil}
	retrier := setupRetrier(conf, logger)
	reporter := report.Reporter{
		Logger: logger,
		Data:   &dataExchange,
		Ticker: time.NewTicker(time.Duration(conf.ReportInterval) * time.Second),
		Done:   done,
		SendAdapter: adapters.NewHTTPReportAdapter(
			logger,
			conf.Addr,
			conf.UpdateURL,
			retrier,
			[]byte(conf.HMACKey),
			conf.RateLimit,
		),
	}
	poller := poll.Poller{
		Logger:    logger,
		PollCount: 0,
		Data:      &dataExchange,
		Ticker:    time.NewTicker(time.Duration(conf.PollInterval) * time.Second),
		Done:      done,
	}
	go reporter.Report()
	go poller.Poll()
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Stopping the agent")
}
