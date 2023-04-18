package sqldb

import (
	"database/sql"

	"github.com/block-wallet/campaigns-service/storage/database/config"
	"github.com/block-wallet/campaigns-service/utils/logger"

	"github.com/block-wallet/campaigns-service/storage/database/sql/migrator"
	"github.com/block-wallet/campaigns-service/storage/database/sql/postgresql"
)

func NewSQLDatabase(dbConfig *config.DBConfig) (*sql.DB, error) {
	var db *sql.DB
	var err error
	var sqlMigrator migrator.SQLMigrator
	switch dbConfig.DBType {
	case config.PostgresDBType:
		{
			db, err = postgresql.NewPosgtresDatabase(dbConfig.SQLConfig.Connection)
			if err != nil {
				return nil, err
			}
			sqlMigrator = migrator.NewPostgresMigrator(db)
		}
	}

	if !dbConfig.SQLConfig.SkipMigrations {
		logger.Sugar.Debug("running migrations...")
		err = sqlMigrator.Migrate()
		if err != nil {
			logger.Sugar.Errorf("error running migrations: %v", err.Error())
			return nil, err
		}
	}

	return db, nil
}
