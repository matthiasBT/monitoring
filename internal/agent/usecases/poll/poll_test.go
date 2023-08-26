package poll

import (
	"sort"
	"testing"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/stretchr/testify/assert"
)

func TestPoll(t *testing.T) {
	c := PollerInfra{
		Logger:     logging.SetupLogger(),
		PollCount:  12,
		Data:       &entities.SnapshotWrapper{CurrSnapshot: nil},
		PollTicker: nil,
		Done:       nil,
	}
	poller := Poller{Infra: &c}
	poller.currentSnapshot()
	snap := c.Data.CurrSnapshot
	assert.NotContains(t, snap.Gauges, "PollCount")
	assert.Equalf(t, map[string]int64{"PollCount": 12}, snap.Counters, "Counters don't match")
	gauges := make([]string, 0, len(snap.Gauges))
	for key := range snap.Gauges {
		gauges = append(gauges, key)
	}
	sort.Strings(gauges)

	expectedGauges := []string{
		"Alloc",
		"BuckHashSys",
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
	}
	assert.EqualValues(t, expectedGauges, gauges)
}
