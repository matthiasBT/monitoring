package adapters

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"time"

	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type FileStorage struct {
	Logger        logging.ILogger
	Storage       entities.Storage
	Path          string
	Done          <-chan struct{}
	Tick          <-chan time.Time
	StorageEvents <-chan struct{}
	Lock          *sync.Mutex
	StoreSync     bool
	inited        bool
}

func (fs *FileStorage) Dump() {
	for {
		select {
		case <-fs.Done:
			fs.Logger.Infoln("Stopping the Dump job")
			fs.save()
			return
		case tick := <-fs.Tick:
			if !fs.StoreSync {
				fs.Logger.Infof("Dump job is ticking at %v\n", tick)
				fs.save()
			}
		case <-fs.StorageEvents:
			if fs.StoreSync {
				fs.Logger.Infoln("Received storage event")
				fs.save()
			}
		}
	}
}

func (fs *FileStorage) save() {
	fs.Logger.Infoln("Starting saving the storage data")

	fs.Lock.Lock()
	defer fs.Lock.Unlock()

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
		fs.Logger.Infof("Inited: %v. Dumped: %v\n", fs.inited, string(body))
	}
	fs.Logger.Infoln("Saving complete")
}

func (fs *FileStorage) InitStorage() map[string]*common.Metrics {
	fs.inited = true
	fs.Logger.Infoln("Starting restoring the storage data")
	var result = make(map[string]*common.Metrics)

	file, err := os.OpenFile(fs.Path, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		fs.Logger.Errorf("Can't init storage: %v\n", err.Error())
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
		data, _ := json.Marshal(metrics)
		fs.Logger.Infof("Loaded: %v\n", string(data))
		result[metrics.ID] = &metrics
	}
	fs.Logger.Infoln("Success")
	return result
}
