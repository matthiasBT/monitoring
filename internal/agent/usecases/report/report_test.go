package report

import (
	"crypto/rand"
	"encoding/json"
	"testing"

	"github.com/matthiasBT/monitoring/internal/agent/adapters"
	"github.com/matthiasBT/monitoring/internal/agent/entities"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/stretchr/testify/assert"
)

func TestReporter_report(t *testing.T) {
	hmacKey := make([]byte, 32)
	if _, err := rand.Read(hmacKey); err != nil {
		t.Fatalf("Failed to create HMAC key: %v", err)
	}
	reportAdapter := adapters.NewHTTPReportAdapter(
		logging.SetupLogger(),
		"0.0.0.0:8000",
		"/updates/",
		utils.Retrier{
			Attempts:         1,
			IntervalFirst:    0,
			IntervalIncrease: 0,
			Logger:           logging.SetupLogger(),
		},
		hmacKey,
		nil,
		0,
	)
	reportAdapter.Jobs = make(chan []byte, 1)
	data := entities.SnapshotWrapper{
		CurrSnapshot: &entities.Snapshot{
			Gauges:   map[string]float64{"GFoo": 11.54, "GBar": 22.11},
			Counters: map[string]int64{"CFoo": -1, "CBar": 0},
		},
	}
	r := &Reporter{
		Logger:      logging.SetupLogger(),
		Data:        &data,
		SendAdapter: reportAdapter,
	}
	r.report()
	job := <-reportAdapter.Jobs
	var result []*common.Metrics
	if err := json.Unmarshal(job, &result); err != nil {
		t.Errorf("Failed to unmarshal report adapter result: %v", err)
	}
	assert.Equal(t, len(result), 4)

	assert.Equal(t, result[0].ID, "GFoo")
	assert.Equal(t, result[0].MType, common.TypeGauge)
	assert.Equal(t, *result[0].Value, 11.54)
	assert.Equal(t, result[0].Delta == nil, true)

	assert.Equal(t, result[1].ID, "GBar")
	assert.Equal(t, result[1].MType, common.TypeGauge)
	assert.Equal(t, *result[1].Value, 22.11)
	assert.Equal(t, result[1].Delta == nil, true)

	assert.Equal(t, result[2].ID, "CFoo")
	assert.Equal(t, result[2].MType, common.TypeCounter)
	assert.Equal(t, result[2].Value == nil, true)
	assert.Equal(t, *result[2].Delta, int64(-1))

	assert.Equal(t, result[3].ID, "CBar")
	assert.Equal(t, result[3].MType, common.TypeCounter)
	assert.Equal(t, result[3].Value == nil, true)
	assert.Equal(t, *result[3].Delta, int64(0))
}
