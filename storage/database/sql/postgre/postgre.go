package postgre

import (
	"database/sql"

	"github.com/block-wallet/campaigns-service/utils/logger"
)

func NewPosgtreDatabase(connection string) (*sql.DB, error) {
	logger.Sugar.Debugf("Connecting to: %v", connection)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	logger.Sugar.Debugf("Successfully opened database")
	return db, nil
}
