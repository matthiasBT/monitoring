package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matthiasBT/monitoring/internal/adapters"
	"github.com/matthiasBT/monitoring/internal/collector"
	"github.com/matthiasBT/monitoring/internal/config"
)

const updateURL = "/update"

func main() {
	logger := adapters.SetupLogger()
	conf, err := config.InitAgentConfig()
	if err != nil {
		logger.Fatal(err)
	}
	done := make(chan bool)
	reporterInfra := collector.ReporterInfra{
		Logger:       logger,
		CurrSnapshot: nil,
		ReportTicker: time.NewTicker(time.Duration(conf.ReportInterval) * time.Second),
		Done:         done,
		ServerAddr:   conf.Addr,
		UpdateURL:    updateURL,
	}
	pollerInfra := collector.PollerInfra{
		Logger:       logger,
		PollCount:    0,
		CurrSnapshot: nil,
		PollTicker:   time.NewTicker(time.Duration(conf.PollInterval) * time.Second),
		Done:         done,
	}
	reporter := collector.Reporter{Infra: &reporterInfra}
	go reporter.Report()
	poller := collector.Poller{Infra: &pollerInfra}
	go poller.Poll()
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Stopping the agent")
}
