package adapters

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	common "github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

type FileKeeper struct {
	Logger    logging.ILogger
	Storage   entities.Storage
	Path      string
	Done      <-chan struct{}
	Tick      <-chan time.Time
	Lock      *sync.Mutex
	StoreSync bool
}

func NewFileKeeper(
	conf *server.Config,
	logger logging.ILogger,
	storage entities.Storage,
	done chan struct{},
) entities.Keeper {
	var tickerChan <-chan time.Time
	if conf.StoresSync() {
		tickerChan = make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(*conf.StoreInterval) * time.Second)
		tickerChan = ticker.C
	}
	return &FileKeeper{
		Logger:    logger,
		Storage:   storage,
		Path:      conf.FileStoragePath,
		Done:      done,
		Tick:      tickerChan,
		Lock:      &sync.Mutex{},
		StoreSync: conf.StoresSync(),
	}
}

func (fs *FileKeeper) FlushPeriodic() {
	for {
		select {
		case <-fs.Done:
			fs.Logger.Infoln("Stopping the Flush job")
			fs.Flush()
			return
		case tick := <-fs.Tick:
			if !fs.StoreSync { // the "else" is unreachable here, just a matter of precaution
				fs.Logger.Infof("Flush job is ticking at %v\n", tick)
				fs.Flush()
			}
		}
	}
}

func (fs *FileKeeper) Flush() {
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
	}
	fs.Logger.Infoln("Saving complete")
}

func (fs *FileKeeper) Restore() map[string]*common.Metrics {
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
		result[metrics.ID] = &metrics
	}
	fs.Logger.Infoln("Success")
	return result
}
