package adapters

import (
	"context"
	"os"
	"sync"
	"testing"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
)

func TestFileKeeper_Flush(t *testing.T) {
	tests := []struct {
		name            string
		storageSnapshot []*common.Metrics
		wantErr         bool
	}{
		{
			name: "flush_snapshot",
			storageSnapshot: []*common.Metrics{
				{
					ID:    "BarFoo1",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(44.1),
				},
				{
					ID:    "BarFoo2",
					MType: common.TypeGauge,
					Delta: nil,
					Value: ptrfloat64(55.5),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.CreateTemp("", "example")
			if err != nil {
				t.Fatal("Failed to create a temporary file")
			}
			defer os.Remove(file.Name())
			defer file.Close()
			fs := &FileKeeper{
				Logger: logging.SetupLogger(),
				Path:   file.Name(),
				Lock:   &sync.Mutex{},
			}
			if err := fs.Flush(context.Background(), tt.storageSnapshot); err != nil {
				t.Errorf("Flush() error = %v", err)
			}
			restoredState := fs.Restore()
			if len(restoredState) != len(tt.storageSnapshot) ||
				!compare(restoredState[0], tt.storageSnapshot[0]) ||
				!compare(restoredState[1], tt.storageSnapshot[1]) {
				t.Errorf("State after Restore() is not equal")
			}
		})
	}
}
