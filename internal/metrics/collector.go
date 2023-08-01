package metrics

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

func Collect(pollCnt int) *Snapshot {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	result := &Snapshot{
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
			"PollCount": int64(pollCnt),
		},
	}
	fmt.Println("Created another metrics snapshot")
	return result
}

func Report(wrapper *SnapshotWrapper, interval time.Duration, addr string, patternUpdate string) {
	for {
		if wrapper.CurrSnapshot == nil {
			fmt.Println("Data is not ready yet")
			time.Sleep(interval)
			continue
		}
		fmt.Printf("Reporting snapshot, memory address: %v", wrapper.CurrSnapshot)
		for name, val := range wrapper.CurrSnapshot.Gauges {
			path := buildGaugePath(patternUpdate, name, val)
			reportMetric(addr, path)
		}
		for name, val := range wrapper.CurrSnapshot.Counters {
			path := buildCounterPath(patternUpdate, name, val)
			reportMetric(addr, path)
		}
		time.Sleep(interval)
	}
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
	resp, _ := http.Post(u.String(), "text/plain", nil)
	defer resp.Body.Close()
}
