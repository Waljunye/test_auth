package pg

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

type migrationsConfig interface {
	DbUrl() string
	MigrationsPath() string
}

func MigrateUp(cfg migrationsConfig) (err error) {
	m, err := migrate.New(
		cfg.MigrationsPath(),
		cfg.DbUrl())
	if err != nil {
		return
	}
	version, dirty, err := m.Version()
	log.Info().Uint("version", version).Bool("dirty", dirty).Err(err).Msg("migration info")

	if err != nil {
		return
	}
	if err = m.Up(); err != nil {
		return
	}

	return
}

func MigrateDown(cfg migrationsConfig) (err error) {
	m, err := migrate.New(
		cfg.MigrationsPath(),
		cfg.DbUrl())
	if err != nil {
		return
	}
	if err = m.Down(); err != nil {
		return
	}

	return
}
