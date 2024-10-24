package pg

import "database/sql"

type Db struct {
	db *sql.DB
}

func (db *Db) test() {
	db.db.Begin()
}
