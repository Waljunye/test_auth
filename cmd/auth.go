package main

import (
	"auth/internal/contracts/events"
	"auth/internal/services"
	"auth/internal/stores"
	"auth/libs"
	"database/sql"
	"flag"
	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"strings"

	"auth/cmd/config"
	"auth/internal/contracts/open_api"
	v1 "auth/internal/contracts/open_api/v1"
	"auth/libs/application"
	"auth/libs/listeners"
	"auth/libs/pg"

	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

}

func main() {
	cfgFile := flag.String("cfg-file", "", "define .env file")
	flag.Parse()

	cfg, err := config.LoadFromEnv(cfgFile)
	if err != nil {
		panic(err)
	}

	dbConn, err := sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}

	switch cfg.AppMode() {
	case "app":
		err = run(cfg, dbConn)
		if err != nil {
			panic(err)
		}
		return
	case "migrateUp":
		log.Info().Msg("migrate up")

		err = pg.MigrateUp(cfg)
		if err != nil {
			log.Warn().Err(err).Msg("migrate up failed")
		}
	case "migrateDown":
		err = pg.MigrateDown(cfg)
		if err != nil {
			panic(err)
		}
	}
}

func run(cfg *config.Config, db *sql.DB) (err error) {
	app, err := build(cfg, db)
	if err != nil {
		return
	}

	log.Info().Msg("build finished, starting app")
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	app.Run(stopChan)
	return
}

func build(cfg *config.Config, db *sql.DB) (result app, err error) {
	usersStore := stores.NewUsersStore(db)

	authService := services.NewAuthService(usersStore, cfg)

	baseContract := v1.NewBaseContract()
	authContract := v1.NewAuthContract(authService)

	openApiListener := open_api.New(nil, baseContract, authContract)

	kafkaConfig := &sarama.Config{}

	kafkaDriver := libs.NewKafkaDriver(strings.Split(cfg.KafkaServers(), ","), cfg.KafkaConsumerGroup(), kafkaConfig)

	eventsListener := events.NewBarEventContract(kafkaDriver)
	lsnrs := map[int]listeners.PortListener{
		cfg.PublicApiPort(): openApiListener,
	}
	wrkrs := []listeners.BackgroundWorker{
		eventsListener,
	}

	result = application.New(lsnrs, wrkrs)

	return
}
