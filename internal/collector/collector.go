package collector

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"

	"github.com/matthiasBT/monitoring/internal/interfaces"
)

type Snapshot struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

type SnapshotWrapper struct {
	CurrSnapshot *Snapshot
}

type Context struct {
	Logger       interfaces.ILogger
	PollCount    int64
	CurrSnapshot *Snapshot
	PollTicker   *time.Ticker
	ReportTicker *time.Ticker
	Done         chan bool
	ServerAddr   string
	UpdateURL    string
}

type Reporter struct {
	Ctx *Context
}

type Poller struct {
	Ctx *Context
}

func (p *Poller) Poll() {
	for {
		select {
		case <-p.Ctx.Done:
			p.Ctx.Logger.Infoln("Stopping the Poll job")
			return
		case tick := <-p.Ctx.PollTicker.C:
			p.Ctx.PollCount += 1
			p.Ctx.Logger.Infof("Poll job #%v is ticking at %v\n", p.Ctx.PollCount, tick)
			p.currentSnapshot()
		}
	}
}

func (p *Poller) currentSnapshot() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	p.Ctx.CurrSnapshot = &Snapshot{
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
			"PollCount": p.Ctx.PollCount,
		},
	}
	p.Ctx.Logger.Infoln("Created another metrics snapshot")
}

func (r *Reporter) Report() {
	for {
		select {
		case <-r.Ctx.Done:
			r.Ctx.Logger.Infoln("Stopping the Report job")
			return
		case tick := <-r.Ctx.ReportTicker.C:
			r.Ctx.Logger.Infof("Report job is ticking at %v\n", tick)
			r.report()
		}
	}
}

func (r *Reporter) report() {
	if r.Ctx.CurrSnapshot == nil {
		r.Ctx.Logger.Infoln("Data for report is not ready yet")
		return
	}
	// saving the address of the current snapshot, so it doesn't get overwritten
	snapshot := r.Ctx.CurrSnapshot
	r.Ctx.Logger.Infof("Reporting snapshot, memory address: %v\n", &snapshot)
	for name, val := range snapshot.Gauges {
		path := buildGaugePath(r.Ctx.UpdateURL, name, val)
		r.reportMetric(path)
	}
	for name, val := range snapshot.Counters {
		path := buildCounterPath(r.Ctx.UpdateURL, name, val)
		r.reportMetric(path)
	}
	r.Ctx.Logger.Infoln("All metrics have been reported")
}

func (r *Reporter) reportMetric(path string) {
	u := url.URL{
		Scheme: "http",
		Host:   r.Ctx.ServerAddr,
		Path:   path,
	}
	resp, err := http.Post(u.String(), "text/plain", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		// trying to submit everything we can, hence no aborting the iteration when encountering an error
		r.Ctx.Logger.Infof("Failed to report a metric. POST %v: %v\n", path, err.Error())
		return
	} else {
		r.Ctx.Logger.Infof("Success: POST %v\n", path)
	}
}

func buildGaugePath(patternUpdate string, name string, val float64) string {
	return fmt.Sprintf("%s/%s/%s/%v", patternUpdate, "gauge", name, val)
}

func buildCounterPath(patternUpdate string, name string, val int64) string {
	return fmt.Sprintf("%s/%s/%s/%v", patternUpdate, "counter", name, val)
}
