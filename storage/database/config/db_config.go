package config

type DBType string

const (
	PostgresDBType DBType = "PostgreSQL"
)

type DBConfig struct {
	DBType         DBType
	SQLConfig      *SQLConfig
	MigrationsPath *string
}

func NewDBConfig(DBType DBType, sqlConfig *SQLConfig) *DBConfig {
	return &DBConfig{
		DBType:    DBType,
		SQLConfig: sqlConfig,
	}
}
