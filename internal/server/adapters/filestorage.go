package adapters

import (
	"fmt"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type FileStorage struct {
	Logger        logging.ILogger
	Storage       entities.Storage
	Done          <-chan bool
	Tick          <-chan time.Time
	StorageEvents <-chan struct{}
}

func (fs *FileStorage) Dump() {
	for {
		select {
		case <-fs.Done:
			fs.Logger.Infoln("Stopping the Dump job")
			return
		case tick := <-fs.Tick:
			fs.Logger.Infof("Dump job is ticking at %v\n", tick)
			Dummy()
		case <-fs.StorageEvents:
			fs.Logger.Infoln("Received storage event")
			Dummy()
		}
	}
}

func Dummy() {
	fmt.Println("Hello from Dummy!")
}
