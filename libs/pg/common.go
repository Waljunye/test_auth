package pg

import (
	"database/sql"
	"github.com/rs/zerolog/log"
)

func CloseStmt(stmt *sql.Stmt) {
	if stmt != nil {
		err := stmt.Close()
		if err != nil {
			log.Error().Err(err).Msg("close statement")
		}
	}
}

func CloseConnections(stmt *sql.Stmt, rows *sql.Rows) {
	if stmt != nil {
		err := stmt.Close()
		if err != nil {
			log.Error().Err(err).Msg("close statement")
		}
	}

	if rows != nil {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("close statement")
		}
	}
}
