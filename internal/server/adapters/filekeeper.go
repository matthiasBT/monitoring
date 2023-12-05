// Package adapters provides functionality for managing file-based storage,
// including reading from and writing to a file. It handles JSON serialization
// of metrics data and ensures thread safety with a mutex.
package adapters

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/infra/utils"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// FileKeeper is a struct that manages file operations and holds a logger,
// the path to the file storage, a retrier for handling retry logic, and a mutex for
// synchronizing operations.
type FileKeeper struct {
	Lock    *sync.Mutex     // Mutex for synchronization
	Path    string          // Path to the file storage
	Logger  logging.ILogger // Logger for logging activities
	Retrier utils.Retrier   // Retrier for retry logic
}

// NewFileKeeper creates and returns a new FileKeeper instance with the provided configuration,
// logger, and retrier. It initializes a mutex for thread safety.
func NewFileKeeper(conf *server.Config, logger logging.ILogger, retrier utils.Retrier) entities.Keeper {
	return &FileKeeper{
		Logger:  logger,
		Path:    conf.FileStoragePath,
		Retrier: retrier,
		Lock:    &sync.Mutex{},
	}
}

// Flush writes a snapshot of storage data (metrics) to a file in JSON format,
// ensuring thread safety and error handling.
func (fs *FileKeeper) Flush(ctx context.Context, storageSnapshot []*common.Metrics) error {
	fs.Logger.Infoln("Starting saving the storage data")

	fs.Lock.Lock()
	defer fs.Lock.Unlock()

	file, err := os.OpenFile(fs.Path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fs.Logger.Errorf("Failed to open storage file: %s\n", err.Error())
		return err
	}
	defer file.Close()

	for _, metrics := range storageSnapshot {
		var err error
		var body []byte
		body, err = json.Marshal(metrics)
		if err != nil {
			fs.Logger.Errorf("Failed to marshal a metric: %s, %s\n", metrics.ID, err.Error())
			return err
		}
		if _, err = file.Write(body); err != nil {
			fs.Logger.Errorf("Failed to write a metric to the file %s\n", err.Error())
			return err
		}
		if _, err = file.WriteString("\n"); err != nil {
			fs.Logger.Errorf("Failed to write a newline to the file %s\n", err.Error())
			return err
		}
	}
	fs.Logger.Infoln("Saving complete")
	return nil
}

// Restore reads and returns all metrics from the file storage, parsing them from JSON.
func (fs *FileKeeper) Restore() []*common.Metrics {
	fs.Logger.Infoln("Restoring the storage data")
	var result []*common.Metrics

	file, err := os.OpenFile(fs.Path, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		fs.Logger.Errorf("Failed to open storage file: %v\n", err.Error())
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		metrics := common.Metrics{}
		err = json.Unmarshal([]byte(scanner.Text()), &metrics)
		if err != nil {
			fs.Logger.Errorf("Failed to unmarshal data from file: %v\n", err.Error())
			panic(err)
		}
		result = append(result, &metrics)
	}
	fs.Logger.Infoln("Success")
	return result
}

// Ping is a no-op for the FileKeeper, as it does not require a live connection.
func (fs *FileKeeper) Ping(context.Context) error {
	return nil
}

// Shutdown logs that no shutdown action is needed for the FileKeeper.
func (fs *FileKeeper) Shutdown() {
	fs.Logger.Infoln("No shutdown action needed")
}
