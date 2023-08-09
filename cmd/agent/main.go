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
	conf := config.InitAgentConfig(logger)
	ctx := collector.Context{
		Logger:       logger,
		PollCount:    0,
		CurrSnapshot: nil,
		PollTicker:   time.NewTicker(time.Duration(conf.PollInterval) * time.Second),
		ReportTicker: time.NewTicker(time.Duration(conf.ReportInterval) * time.Second),
		Done:         make(chan bool),
		ServerAddr:   conf.Addr,
		UpdateURL:    updateURL,
	}
	reporter := collector.Reporter{Ctx: &ctx}
	go reporter.Report()
	poller := collector.Poller{Ctx: &ctx}
	go poller.Poll()
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Stopping the agent")
}
