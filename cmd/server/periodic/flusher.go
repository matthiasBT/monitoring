package periodic

import (
	"context"
	"time"

	"github.com/matthiasBT/monitoring/internal/infra/config/server"
	"github.com/matthiasBT/monitoring/internal/infra/logging"
	"github.com/matthiasBT/monitoring/internal/server/entities"
)

// todo: get rid of

type Flusher struct {
	Storage entities.Storage
	Keeper  entities.Keeper
	Tick    <-chan time.Time
	IsSync  bool
	Done    <-chan struct{}
	Logger  logging.ILogger
}

func NewFlusher(
	conf *server.Config,
	logger logging.ILogger,
	storage entities.Storage,
	keeper entities.Keeper,
	done chan struct{},
) Flusher {
	var tickerChan <-chan time.Time
	if conf.FlushesSync() {
		tickerChan = make(chan time.Time) // will never be used
	} else {
		ticker := time.NewTicker(time.Duration(*conf.StoreInterval) * time.Second)
		tickerChan = ticker.C
	}
	return Flusher{
		Logger:  logger,
		Storage: storage,
		Done:    done,
		Tick:    tickerChan,
		IsSync:  conf.FlushesSync(),
		Keeper:  keeper,
	}
}

func (pf *Flusher) Flush(ctx context.Context) {
	pf.Logger.Infoln("Launching the periodic Flush job")
	for {
		select {
		case <-pf.Done:
			pf.Logger.Infoln("Stopping the periodic Flush job")
			pf.flush(ctx, true)
			return
		case tick := <-pf.Tick:
			if !pf.IsSync { // the "else" is unreachable here, just a matter of precaution
				pf.Logger.Infof("The periodic Flush job is ticking at %v\n", tick)
				pf.flush(ctx, false)
			}
		}
	}
}

func (pf *Flusher) flush(ctx context.Context, mustSucceed bool) {
	data, err := pf.Storage.Snapshot(ctx)
	if err != nil {
		pf.Logger.Errorf("Failed to get data from storage: %s\n", err.Error())
		if mustSucceed {
			panic(err)
		}
		return
	}
	if err := pf.Keeper.Flush(ctx, data); err != nil {
		pf.Logger.Errorf("Failed to flush data: %s\n", err.Error())
		if mustSucceed {
			panic(err)
		}
		return
	}
}
