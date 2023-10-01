package poll

import (
	"fmt"
	"sort"
	"testing"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	poller := Poller{
		Logger:    logging.SetupLogger(),
		PollCount: 12,
		Data:      &entities.SnapshotWrapper{CurrSnapshot: nil},
		Ticker:    nil,
		Done:      nil,
	}
	poller.currentSnapshot()
	assert.NotContains(t, poller.Data.CurrSnapshot.Gauges, "PollCount")
	assert.Equalf(t, map[string]int64{"PollCount": 12}, poller.Data.CurrSnapshot.Counters, "Counters don't match")
	gauges := make([]string, 0, len(poller.Data.CurrSnapshot.Gauges))
	for key := range poller.Data.CurrSnapshot.Gauges {
		gauges = append(gauges, key)
	}
	sort.Strings(gauges)
	expectedGauges := []string{
		"Alloc",
		"BuckHashSys",
		"FreeMemory",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"RandomValue",
		"StackInuse",
		"StackSys",
		"Sys",
		"TotalAlloc",
		"TotalMemory",
	}
	if cpuCount, err := cpu.Counts(true); err != nil {
		t.Fatalf("Failed to get the number of CPUs: %v", err)
	} else {
		for i := 1; i <= cpuCount; i++ {
			name := fmt.Sprintf("CPUutilization%d", i)
			expectedGauges = append(expectedGauges, name)
		}
		sort.Strings(expectedGauges)
	}
	assert.EqualValues(t, expectedGauges, gauges)
}
