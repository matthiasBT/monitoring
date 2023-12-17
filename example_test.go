package monitoring_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
	"github.com/sirupsen/logrus"
)

func Example() {
	logger := logging.SetupLogger()
	logger.SetLevel(logrus.FatalLevel) // to avoid printing unnecessary logs
	storage := adapters.NewMemStorage(nil, nil, logger, nil)
	controller := usecases.NewBaseController(logger, storage, "/")

	updateCounter(100500, controller)
	updateGauge(5.5, controller)
	getCounter(controller)
	getGauge(controller)
	updateCounter(11, controller)
	updateGauge(1.5, controller)
	getCounter(controller)
	getGauge(controller)

	// Output:
	// {id:FooBar, type:counter, delta:100500}
	// {id:BarFoo, type:gauge, value:5.5}
	// {id:FooBar, type:counter, delta:100511}
	// {id:BarFoo, type:gauge, value:1.5}
}

func updateCounter(value int64, controller *usecases.BaseController) {
	w := httptest.NewRecorder()
	newCounter := common.Metrics{
		ID:    "FooBar",
		MType: "counter",
		Delta: &value,
		Value: nil,
	}
	body, _ := json.Marshal(newCounter)
	addCounterReq := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(body))
	addCounterReq.Header.Set("Content-Type", "application/json")
	controller.UpdateMetric(w, addCounterReq)
}

func updateGauge(value float64, controller *usecases.BaseController) {
	w := httptest.NewRecorder()
	newGauge := common.Metrics{
		ID:    "BarFoo",
		MType: "gauge",
		Delta: nil,
		Value: &value,
	}
	body, _ := json.Marshal(newGauge)
	addGaugeReq := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewReader(body))
	addGaugeReq.Header.Set("Content-Type", "application/json")
	controller.UpdateMetric(w, addGaugeReq)
}

func getCounter(controller *usecases.BaseController) {
	w := httptest.NewRecorder()
	newCounter := common.Metrics{
		ID:    "FooBar",
		MType: "counter",
		Delta: nil,
		Value: nil,
	}
	body, _ := json.Marshal(newCounter)
	getCounterReq := httptest.NewRequest(http.MethodGet, "/value/", bytes.NewReader(body))
	getCounterReq.Header.Set("Content-Type", "application/json")
	controller.GetMetric(w, getCounterReq)
	printSorted(w.Body.Bytes())
}

func getGauge(controller *usecases.BaseController) {
	w := httptest.NewRecorder()
	newGauge := common.Metrics{
		ID:    "BarFoo",
		MType: "gauge",
		Delta: nil,
		Value: nil,
	}
	body, _ := json.Marshal(newGauge)
	getGaugeReq := httptest.NewRequest(http.MethodGet, "/value/", bytes.NewReader(body))
	getGaugeReq.Header.Set("Content-Type", "application/json")
	controller.GetMetric(w, getGaugeReq)
	printSorted(w.Body.Bytes())
}

func printSorted(body []byte) {
	var result common.Metrics
	json.Unmarshal(body, &result)
	if result.MType == common.TypeGauge {
		fmt.Printf("{%s:%v, %s:%v, %s:%v}\n", "id", result.ID, "type", result.MType, "value", *result.Value)
	} else {
		fmt.Printf("{%s:%v, %s:%v, %s:%v}\n", "id", result.ID, "type", result.MType, "delta", *result.Delta)
	}
}
