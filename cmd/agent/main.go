package main

import (
	"fmt"
	"github.com/matthiasBT/monitoring/internal/metrics"
	"time"
)

const pollInterval = 2 * time.Second
const reportInterval = 10 * time.Second
const serverAddr = ":8080"
const patternUpdate = "/update"

// todo: handle errors better?
func main() {
	pollCnt := 0
	var wrapper metrics.SnapshotWrapper
	go metrics.Report(&wrapper, reportInterval, serverAddr, patternUpdate)
	for {
		pollCnt += 1
		fmt.Printf("Starting iteration %v\n", pollCnt)
		wrapper.CurrSnapshot = metrics.Collect(pollCnt)
		fmt.Printf("Finished iteration %v\n", pollCnt)
		time.Sleep(pollInterval)
	}
}
