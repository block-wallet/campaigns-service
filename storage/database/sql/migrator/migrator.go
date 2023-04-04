package migrator

type SQLMigrator interface {
	Migrate() error
}
