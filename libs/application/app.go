package application

import (
	"auth/libs/listeners"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

type app struct {
	listeners map[int]listeners.PortListener
	workers   []listeners.BackgroundWorker
}

func New(listeners map[int]listeners.PortListener, workers []listeners.BackgroundWorker) *app {
	return &app{
		listeners: listeners,
		workers:   workers,
	}
}

func (app *app) Run(stopChan chan os.Signal) {
	for port, listener := range app.listeners {
		log.Info().Msg(fmt.Sprintf("running: %s. port: %v", listener.Info(), port))
		go func(listener listeners.PortListener, port int) {
			err := listener.Run(port)
			if err != nil {
				log.Error().Err(err)
				return
			}
		}(listener, port)
	}
	for _, worker := range app.workers {
		log.Info().Msg(fmt.Sprintf("running: %s", worker.Info()))
		go func(worker listeners.BackgroundWorker) {
			err := worker.Start()
			if err != nil {
				log.Error().Err(err)
				return
			}
		}(worker)
	}

	<-stopChan

	for _, listener := range app.listeners {
		err := listener.Stop()
		if err != nil {
			log.Error().Msg(fmt.Sprintf("stop listener: %v", err))
			return
		}
	}

	log.Info().Msg("all listeners stopped")

	for _, worker := range app.workers {
		err := worker.Stop()
		if err != nil {
			log.Error().Msg(fmt.Sprintf("stop worker: %v", err))
		}
	}
	log.Info().Msg("all workers stopped")
}
