package collector

import (
	"sort"
	"testing"

	"github.com/matthiasBT/monitoring/internal/adapters"
	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	c := PollerInfra{
		Logger:       adapters.SetupLogger(),
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

func Test_buildCounterPath(t *testing.T) {
	type args struct {
		patternUpdate string
		name          string
		val           int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "metric 1",
			args: args{
				patternUpdate: "/update",
				name:          "FooBar",
				val:           23,
			},
			want: "/update/counter/FooBar/23",
		},
		{
			name: "metric 2",
			args: args{
				patternUpdate: "/send",
				name:          "BarFoo",
				val:           1,
			},
			want: "/send/counter/BarFoo/1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildCounterPath(tt.args.patternUpdate, tt.args.name, tt.args.val); got != tt.want {
				t.Errorf("buildCounterPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildGaugePath(t *testing.T) {
	type args struct {
		patternUpdate string
		name          string
		val           float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "metric 1",
			args: args{
				patternUpdate: "/update",
				name:          "FooBar",
				val:           23.4,
			},
			want: "/update/gauge/FooBar/23.4",
		},
		{
			name: "metric 2",
			args: args{
				patternUpdate: "/send",
				name:          "BarFoo",
				val:           0.01,
			},
			want: "/send/gauge/BarFoo/0.01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildGaugePath(tt.args.patternUpdate, tt.args.name, tt.args.val); got != tt.want {
				t.Errorf("buildGaugePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
