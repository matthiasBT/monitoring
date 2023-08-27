package report

import (
	"time"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

type Reporter struct {
	Logger      logging.ILogger
	Data        *entities.SnapshotWrapper
	Ticker      *time.Ticker
	Done        <-chan bool
	SendAdapter entities.IReporter
}

func (r *Reporter) Report() {
	for {
		select {
		case <-r.Done:
			r.Logger.Infoln("Stopping the Report job")
			return
		case tick := <-r.Ticker.C:
			r.Logger.Infof("Report job is ticking at %v\n", tick)
			r.report()
		}
	}
}

func (r *Reporter) report() {
	if r.Data.CurrSnapshot == nil {
		r.Logger.Infoln("Data for report is not ready yet")
		return
	}
	// saving the address of the current snapshot, so it doesn't get overwritten
	snapshot := r.Data.CurrSnapshot
	r.Logger.Infof("Reporting snapshot, memory address: %v\n", &snapshot)
	for name, val := range snapshot.Gauges {
		metric := common.Metrics{
			ID:    name,
			MType: common.TypeGauge,
			Delta: nil,
			Value: &val,
		}
		r.send(metric)
	}
	for name, val := range snapshot.Counters {
		metric := common.Metrics{
			ID:    name,
			MType: common.TypeCounter,
			Delta: &val,
			Value: nil,
		}
		r.send(metric)
	}
	r.Logger.Infoln("All metrics have been reported")
}

func (r *Reporter) send(metric common.Metrics) {
	if err := r.SendAdapter.Report(&metric); err != nil {
		r.Logger.Errorf("Failed to report a metric: %v. Error: %v\n", metric, err)
	}
}
