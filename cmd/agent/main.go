package main

import (
	"fmt"
	"github.com/matthiasBT/monitoring/internal/collector"
	"github.com/matthiasBT/monitoring/internal/config"
	"time"
)

const patternUpdate = "/update"

func main() {
	conf := config.InitAgentConfig()
	pollCnt := 0
	var wrapper collector.SnapshotWrapper
	go collector.Report(&wrapper, conf.ReportInterval, conf.ServerAddr, patternUpdate)
	for {
		pollCnt += 1
		fmt.Printf("Starting iteration %v\n", pollCnt)
		wrapper.CurrSnapshot = collector.Collect(pollCnt)
		fmt.Printf("Finished iteration %v\n", pollCnt)
		time.Sleep(conf.PollInterval)
	}
}
