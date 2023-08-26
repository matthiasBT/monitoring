package poll

import (
	"sort"
	"testing"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/stretchr/testify/assert"
)

func TestPoll(t *testing.T) {
	c := PollerInfra{
		Logger:       logging.SetupLogger(),
		PollCount:    12,
		CurrSnapshot: nil,
		PollTicker:   nil,
		Done:         nil,
	}
	poller := Poller{Infra: &c}
	poller.currentSnapshot()
	assert.NotContains(t, c.CurrSnapshot.Gauges, "PollCount")
	assert.Equalf(t, map[string]int64{"PollCount": 12}, c.CurrSnapshot.Counters, "Counters don't match")
	gauges := make([]string, 0, len(c.CurrSnapshot.Gauges))
	for key := range c.CurrSnapshot.Gauges {
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
