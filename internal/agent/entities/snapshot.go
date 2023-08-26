package entities

type Snapshot struct {
	Gauges   map[string]float64
	Counters map[string]int64
}
