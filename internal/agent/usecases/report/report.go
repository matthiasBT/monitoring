// Package report contains the functionalities for reporting system metrics.
// It includes structures and methods for periodically gathering and sending
// metrics to a specified reporting service.

package report

import (
	"time"

	"github.com/matthiasBT/monitoring/internal/agent/entities"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

// Reporter is responsible for periodically reporting system metrics. It retrieves
// metrics from a SnapshotWrapper, formats them, and sends them using a provided
// reporting adapter. Reporting occurs at intervals defined by a Ticker and can
// be stopped via a Done channel.
type Reporter struct {
	// Logger is used to log informational and error messages during the reporting process.
	Logger logging.ILogger

	// Data holds the current snapshot of the system metrics to be reported.
	Data *entities.SnapshotWrapper

	// Ticker controls the intervals at which the metrics are reported.
	Ticker *time.Ticker

	// Done is a channel used to signal when the reporting should be stopped.
	Done <-chan bool

	// SendAdapter is an interface that defines the method of sending the metrics.
	SendAdapter entities.IReporter
}

// Report continuously reports system metrics at intervals defined by the Reporter's Ticker.
// It logs each reporting event and calls the internal report method. Reporting can
// be stopped by sending a signal on the Reporter's Done channel.
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
	var batch = make([]*common.Metrics, 0, len(snapshot.Gauges)+len(snapshot.Gauges))
	for name, val := range snapshot.Gauges {
		metric := common.Metrics{
			ID:    name,
			MType: common.TypeGauge,
			Delta: nil,
			Value: &val,
		}
		batch = append(batch, &metric)
	}
	for name, val := range snapshot.Counters {
		metric := common.Metrics{
			ID:    name,
			MType: common.TypeCounter,
			Delta: &val,
			Value: nil,
		}
		batch = append(batch, &metric)
	}
	if len(batch) == 0 {
		r.Logger.Infoln("Nothing to report yet")
		return
	}
	r.Logger.Infoln("All metrics have been prepared for report")
	if err := r.SendAdapter.ReportBatch(batch); err != nil {
		r.Logger.Errorf("Failed to report a batch. Error: %v\n", err.Error())
	}
}

func (r *Reporter) send(metric common.Metrics) {
	if err := r.SendAdapter.Report(&metric); err != nil {
		r.Logger.Errorf("Failed to report a metric: %v. Error: %v\n", metric, err)
	}
}
