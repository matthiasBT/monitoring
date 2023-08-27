package poll

import (
	"math/rand"
	"runtime"
	"time"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

type PollerInfra struct {
}

type Poller struct {
	Logger    logging.ILogger
	PollCount int64
	Data      *entities.SnapshotWrapper
	Ticker    *time.Ticker
	Done      <-chan bool
}

func (p *Poller) Poll() {
	for {
		select {
		case <-p.Done:
			p.Logger.Infoln("Stopping the Poll job")
			return
		case tick := <-p.Ticker.C:
			p.PollCount += 1
			p.Logger.Infof("Poll job #%v is ticking at %v\n", p.PollCount, tick)
			p.currentSnapshot()
		}
	}
}

func (p *Poller) currentSnapshot() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	snapshot := &entities.Snapshot{
		Gauges: map[string]float64{
			"Alloc":         float64(rtm.Alloc),
			"BuckHashSys":   float64(rtm.BuckHashSys),
			"Frees":         float64(rtm.Frees),
			"GCCPUFraction": rtm.GCCPUFraction,
			"GCSys":         float64(rtm.GCSys),
			"HeapAlloc":     float64(rtm.HeapAlloc),
			"HeapIdle":      float64(rtm.HeapIdle),
			"HeapInuse":     float64(rtm.HeapInuse),
			"HeapObjects":   float64(rtm.HeapObjects),
			"HeapReleased":  float64(rtm.HeapReleased),
			"HeapSys":       float64(rtm.HeapSys),
			"LastGC":        float64(rtm.LastGC),
			"Lookups":       float64(rtm.Lookups),
			"MCacheInuse":   float64(rtm.MCacheInuse),
			"MCacheSys":     float64(rtm.MCacheSys),
			"MSpanInuse":    float64(rtm.MSpanInuse),
			"MSpanSys":      float64(rtm.MSpanSys),
			"Mallocs":       float64(rtm.Mallocs),
			"NextGC":        float64(rtm.NextGC),
			"NumForcedGC":   float64(rtm.NumForcedGC),
			"NumGC":         float64(rtm.NumGC),
			"OtherSys":      float64(rtm.OtherSys),
			"PauseTotalNs":  float64(rtm.PauseTotalNs),
			"StackInuse":    float64(rtm.StackInuse),
			"StackSys":      float64(rtm.StackSys),
			"Sys":           float64(rtm.Sys),
			"TotalAlloc":    float64(rtm.TotalAlloc),
			"RandomValue":   rand.Float64(),
		},
		Counters: map[string]int64{
			"PollCount": p.PollCount,
		},
	}
	p.Data.CurrSnapshot = snapshot
	p.Logger.Infoln("Created another metrics snapshot")
}
