package collector

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type Snapshot struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

type SnapshotWrapper struct {
	CurrSnapshot *Snapshot
}

type Context struct {
	PollCount    int64
	CurrSnapshot *Snapshot
	PollTicker   *time.Ticker
	ReportTicker *time.Ticker
	Done         chan bool
	ServerAddr   string
	UpdateURL    string
}

func Poll(c *Context) {
	fmt.Println(">>")
	for {
		select {
		case <-c.Done:
			fmt.Println("Stopping the Poll job")
			return
		case tick := <-c.PollTicker.C:
			c.PollCount += 1
			fmt.Printf("Poll job #%v is ticking at %v\n", c.PollCount, tick)
			currentSnapshot(c)
		}
	}
}

func currentSnapshot(c *Context) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	c.CurrSnapshot = &Snapshot{
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
			"PollCount": c.PollCount,
		},
	}
	fmt.Println("Created another metrics snapshot")
}

func Report(c *Context) {
	for {
		select {
		case <-c.Done:
			fmt.Println("Stopping the Report job")
			return
		case tick := <-c.ReportTicker.C:
			fmt.Printf("Report job is ticking at %v\n", tick)
			report(c)
		}
	}
}

func report(c *Context) {
	if c.CurrSnapshot == nil {
		fmt.Println("Data for report is not ready yet")
		return
	}
	// saving the address of the current snapshot, so it doesn't get overwritten
	snapshot := c.CurrSnapshot
	fmt.Printf("Reporting snapshot, memory address: %v\n", &snapshot)
	for name, val := range snapshot.Gauges {
		path := buildGaugePath(c.UpdateURL, name, val)
		reportMetric(c.ServerAddr, path)
	}
	for name, val := range snapshot.Counters {
		path := buildCounterPath(c.UpdateURL, name, val)
		reportMetric(c.ServerAddr, path)
	}
	fmt.Println("All metrics have been reported")
}

func buildGaugePath(patternUpdate string, name string, val float64) string {
	return fmt.Sprintf("%s/%s/%s/%v", patternUpdate, "gauge", name, val)
}

func buildCounterPath(patternUpdate string, name string, val int64) string {
	return fmt.Sprintf("%s/%s/%s/%v", patternUpdate, "counter", name, val)
}

func reportMetric(addr string, path string) {
	u := url.URL{
		Scheme: "http",
		Host:   addr,
		Path:   path,
	}
	resp, err := http.Post(u.String(), "text/plain", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		// trying to submit everything we can, hence no aborting the iteration when encountering an error
		fmt.Printf("Failed to report a metric. POST %v: %v\n", path, err.Error())
		return
	} else {
		fmt.Printf("Success: POST %v\n", path)
	}
}
