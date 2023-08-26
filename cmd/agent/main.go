package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	"github.com/matthiasBT/monitoring/internal/agent/usecases/poll"
	"github.com/matthiasBT/monitoring/internal/agent/usecases/report"
	"github.com/matthiasBT/monitoring/internal/infra/config/agent"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

const updateURL = "/update"

func main() {
	logger := logging.SetupLogger()
	conf, err := agent.InitAgentConfig()
	if err != nil {
		logger.Fatal(err)
	}
	done := make(chan bool)
	dataExchange := entities.SnapshotWrapper{CurrSnapshot: nil}
	reporterInfra := report.ReporterInfra{
		Logger:       logger,
		Data:         &dataExchange,
		ReportTicker: time.NewTicker(time.Duration(conf.ReportInterval) * time.Second),
		Done:         done,
		ServerAddr:   conf.Addr,
		UpdateURL:    updateURL,
	}
	pollerInfra := poll.PollerInfra{
		Logger:     logger,
		PollCount:  0,
		Data:       &dataExchange,
		PollTicker: time.NewTicker(time.Duration(conf.PollInterval) * time.Second),
		Done:       done,
	}
	reporter := report.Reporter{Infra: &reporterInfra}
	go reporter.Report()
	poller := poll.Poller{Infra: &pollerInfra}
	go poller.Poll()
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Stopping the agent")
}
