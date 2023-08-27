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
)

func main() {
	logger := logging.SetupLogger()
	conf, err := agent.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}
	done := make(<-chan bool)
	dataExchange := entities.SnapshotWrapper{CurrSnapshot: nil}
	reporter := report.Reporter{
		Logger: logger,
		Data:   &dataExchange,
		Ticker: time.NewTicker(time.Duration(conf.ReportInterval) * time.Second),
		Done:   done,
		SendAdapter: &adapters.HTTPReportAdapter{
			Logger:     logger,
			ServerAddr: conf.Addr,
			UpdateURL:  conf.UpdateURL,
		},
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
