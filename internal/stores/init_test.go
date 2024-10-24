package stores

import (
	"database/sql"
	"flag"
	"os"
	"testing"

	"auth/cmd/config"
)

var (
	usersStore *UsersStore
	dbConn     *sql.DB
)

func TestMain(m *testing.M) {
	cfgFile := flag.String("cfg-file", "", "define .env file")
	flag.Parse()

	cfg, err := config.LoadFromEnv(cfgFile)
	if err != nil {
		panic(err)
	}

	dbConn, err = sql.Open("postgres", cfg.DbUrl())
	if err != nil {
		panic(err)
	}
	usersStore = NewUsersStore(dbConn)

	os.Exit(m.Run())
}
