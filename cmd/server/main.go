package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/matthiasBT/monitoring/internal/infra/compression"
	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/entities"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/adapters"
	"github.com/matthiasBT/monitoring/internal/server/usecases"
)

func setupServer(logger logging.ILogger, controller *usecases.BaseController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logging.Middleware(logger))
	r.Use(compression.MiddlewareReader)
	r.Use(compression.MiddlewareWriter)
	r.Mount("/", controller.Route())
	return r
}

func main() {
	logger := logging.SetupLogger()
	conf, err := server.InitConfig()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Infof("Server config: %v\n", *conf)

	storageEvents := make(chan struct{})
	storage := &adapters.MemStorage{
		Metrics: make(map[string]*entities.Metrics),
		Logger:  logger,
		Events:  storageEvents,
	}

	done := make(chan bool)
	var tickerChan <-chan time.Time
	if conf.StoresSync() {
		tickerChan = make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(*conf.StoreInterval) * time.Second)
		tickerChan = ticker.C
	}
	fileStorage := adapters.FileStorage{
		Logger:        logger,
		Storage:       storage,
		Path:          conf.FileStoragePath,
		Done:          done,
		Tick:          tickerChan,
		StorageEvents: storageEvents,
	}
	go fileStorage.Dump()
	controller := usecases.NewBaseController(logger, storage, conf.TemplatePath)

	r := setupServer(logger, controller)
	logger.Fatal(http.ListenAndServe(conf.Addr, r))

	// TODO: implement graceful shutdown
	//quitChannel := make(chan os.Signal, 1)
	//signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	//<-quitChannel
	//fmt.Println("Stopping the server")
}
