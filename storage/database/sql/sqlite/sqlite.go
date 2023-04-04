package sqlite

import (
	"database/sql"

	"github.com/block-wallet/campaigns-service/utils/logger"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLiteDatabase(connection string) (*sql.DB, error) {
	logger.Sugar.Debugf("Connecting to: %v", connection)
	db, err := sql.Open("sqlite3", connection)
	if err != nil {
		return nil, err
	}
	logger.Sugar.Debugf("Successfully opened database")
	return db, nil
}

func (db *SQLite) Ping() error {
	return db.db.Ping()
}
