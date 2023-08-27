package adapters

import (
	"encoding/json"
	"os"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type FileStorage struct {
	Logger        logging.ILogger
	Storage       entities.Storage
	Path          string
	Done          <-chan bool
	Tick          <-chan time.Time
	StorageEvents <-chan struct{}
}

func (fs *FileStorage) Dump() {
	for {
		select {
		case <-fs.Done:
			fs.Logger.Infoln("Stopping the Dump job")
			fs.save() // TODO: graceful shutdown
		case tick := <-fs.Tick:
			fs.Logger.Infof("Dump job is ticking at %v\n", tick)
			fs.save()
		case <-fs.StorageEvents:
			fs.Logger.Infoln("Received storage event")
			fs.save()
		}
	}
}

func (fs *FileStorage) save() {
	// TODO: add mutex, must be the same as in the storage!
	fs.Logger.Infoln("Starting saving the storage data")
	data, err := fs.Storage.GetAll()
	if err != nil {
		fs.Logger.Errorf("Failed to receive data from storage: %s\n", err.Error())
		return
	}

	file, err := os.OpenFile(fs.Path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fs.Logger.Errorf("Failed to open storage file: %s\n", err.Error())
		return
	}
	defer file.Close()

	for _, metrics := range data {
		body, err := json.Marshal(metrics)
		if err != nil {
			fs.Logger.Errorf("Failed to marshal a metric: %s, %s\n", metrics.ID, err.Error())
			return
		}
		if _, err := file.Write(body); err != nil {
			fs.Logger.Errorf("Failed to write a metric to the file %s\n", err.Error())
			return
		}
		if _, err = file.WriteString("\n"); err != nil {
			fs.Logger.Errorf("Failed to write a newline to the file %s\n", err.Error())
			return
		}
	}
	fs.Logger.Infoln("Saving complete")
}
