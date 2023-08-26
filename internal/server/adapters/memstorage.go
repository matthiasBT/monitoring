package adapters

import (
	"fmt"
	"strconv"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type MemStorage struct {
	MetricsGauge   map[string]float64
	MetricsCounter map[string]int64
	Logger         logging.ILogger
}

func (storage *MemStorage) Add(update entities.MetricUpdate) {
	storage.Logger.Infof("Updating metrics with %+v\n", update)
	switch update.Type {
	case entities.TypeGauge:
		storage.Logger.Infof("Old metric value: %f\n", storage.MetricsGauge[update.Name])
		val, _ := strconv.ParseFloat(update.Value, 64)
		storage.MetricsGauge[update.Name] = val
		storage.Logger.Infof("New metric value: %f\n", storage.MetricsGauge[update.Name])
	case entities.TypeCounter:
		storage.Logger.Infof("Old metric value: %d\n", storage.MetricsCounter[update.Name])
		val, _ := strconv.ParseInt(update.Value, 10, 64)
		storage.MetricsCounter[update.Name] += val
		storage.Logger.Infof("New metric value: %d\n", storage.MetricsCounter[update.Name])
	}
}

func (storage *MemStorage) Get(mType string, name string) (string, error) {
	storage.Logger.Infof("Getting metric of type %s named %s\n", mType, name)
	switch mType {
	case entities.TypeGauge:
		if val, ok := storage.MetricsGauge[name]; ok {
			res := strconv.FormatFloat(val, 'f', -1, 64)
			return res, nil
		}
		storage.Logger.Infoln("No such Gauge metric")
		return "", entities.ErrUnknownMetricName
	case entities.TypeCounter:
		if val, ok := storage.MetricsCounter[name]; ok {
			res := fmt.Sprintf("%d", val)
			return res, nil
		}
		storage.Logger.Infoln("No such Counter metric")
		return "", entities.ErrUnknownMetricName
	default:
		storage.Logger.Infoln("Invalid metric type")
		return "", entities.ErrInvalidMetricType
	}
}

func (storage *MemStorage) GetAll() map[string]string {
	res := make(map[string]string, len(storage.MetricsGauge)+len(storage.MetricsCounter))
	for name := range storage.MetricsGauge {
		valStr, _ := storage.Get(entities.TypeGauge, name)
		res[name] = valStr
	}
	for name := range storage.MetricsCounter {
		valStr, _ := storage.Get(entities.TypeCounter, name)
		res[name] = valStr
	}
	return res
}
