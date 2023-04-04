package config

type SQLConfig struct {
	Connection     string
	SkipMigrations bool
}

func NewSQLConfig(connection string, skipMigrations bool) *SQLConfig {
	return &SQLConfig{
		Connection:     connection,
		SkipMigrations: skipMigrations,
	}
}
