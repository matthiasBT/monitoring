package report

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

// TODO: add adapters layer

type ReporterInfra struct {
	Logger       logging.ILogger
	CurrSnapshot *entities.Snapshot
	ReportTicker *time.Ticker
	Done         chan bool
	ServerAddr   string
	UpdateURL    string
}

type Reporter struct {
	Infra *ReporterInfra
}

func (r *Reporter) Report() {
	for {
		select {
		case <-r.Infra.Done:
			r.Infra.Logger.Infoln("Stopping the Report job")
			return
		case tick := <-r.Infra.ReportTicker.C:
			r.Infra.Logger.Infof("Report job is ticking at %v\n", tick)
			r.report()
		}
	}
}

func (r *Reporter) report() {
	if r.Infra.CurrSnapshot == nil {
		r.Infra.Logger.Infoln("Data for report is not ready yet")
		return
	}
	// saving the address of the current snapshot, so it doesn't get overwritten
	snapshot := r.Infra.CurrSnapshot
	r.Infra.Logger.Infof("Reporting snapshot, memory address: %v\n", &snapshot)
	for name, val := range snapshot.Gauges {
		path := buildGaugePath(r.Infra.UpdateURL, name, val)
		r.reportMetric(path)
	}
	for name, val := range snapshot.Counters {
		path := buildCounterPath(r.Infra.UpdateURL, name, val)
		r.reportMetric(path)
	}
	r.Infra.Logger.Infoln("All metrics have been reported")
}

func (r *Reporter) reportMetric(path string) {
	u := url.URL{
		Scheme: "http",
		Host:   r.Infra.ServerAddr,
		Path:   path,
	}
	resp, err := http.Post(u.String(), "text/plain", nil)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		// trying to submit everything we can, hence no aborting the iteration when encountering an error
		r.Infra.Logger.Infof("Failed to report a metric. POST %v: %v\n", path, err.Error())
		return
	} else {
		r.Infra.Logger.Infof("Success: POST %v\n", path)
	}
}

func buildGaugePath(patternUpdate string, name string, val float64) string {
	return fmt.Sprintf("%s/%s/%s/%v", patternUpdate, "gauge", name, val)
}

func buildCounterPath(patternUpdate string, name string, val int64) string {
	return fmt.Sprintf("%s/%s/%s/%v", patternUpdate, "counter", name, val)
}
