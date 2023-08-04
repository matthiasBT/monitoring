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
	go collector.Report(&wrapper, time.Duration(conf.ReportInterval)*time.Second, conf.ServerAddr, patternUpdate)
	for {
		pollCnt += 1
		fmt.Printf("Starting iteration %v\n", pollCnt)
		wrapper.CurrSnapshot = collector.Collect(pollCnt)
		fmt.Printf("Finished iteration %v\n", pollCnt)
		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
	}
}
