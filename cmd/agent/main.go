package main

import (
	"fmt"
	"github.com/matthiasBT/monitoring/internal/collector"
	"time"
)

const pollInterval = 2 * time.Second
const reportInterval = 10 * time.Second
const serverAddr = ":8080"
const patternUpdate = "/update"

func main() {
	pollCnt := 0
	var wrapper collector.SnapshotWrapper
	go collector.Report(&wrapper, reportInterval, serverAddr, patternUpdate)
	for {
		pollCnt += 1
		fmt.Printf("Starting iteration %v\n", pollCnt)
		wrapper.CurrSnapshot = collector.Collect(pollCnt)
		fmt.Printf("Finished iteration %v\n", pollCnt)
		time.Sleep(pollInterval)
	}
}
