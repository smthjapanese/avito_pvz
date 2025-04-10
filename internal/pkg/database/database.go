package database

import (
	"database/sql"
)

type Database struct {
	*sql.DB
}

func (db *Database) Close() error {
	return db.DB.Close()
}
