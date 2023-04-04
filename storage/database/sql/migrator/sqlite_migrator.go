package migrator

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sqlite_migrations
var SQLiteMigrations embed.FS

type sqlitemigrator struct {
	db *sql.DB
}

func NewSQLiteMigrator(db *sql.DB) *sqlitemigrator {
	return &sqlitemigrator{
		db: db,
	}
}

func (migrator *sqlitemigrator) Migrate() error {
	sourceInstance, err := httpfs.New(http.FS(SQLiteMigrations), "sqlite_migrations")
	if err != nil {
		return fmt.Errorf("invalid source instance, %w", err)
	}
	targetInstance, err := sqlite.WithInstance(migrator.db, new(sqlite.Config))
	if err != nil {
		return fmt.Errorf("invalid target sqlite instance, %w", err)
	}
	m, err := migrate.NewWithInstance(
		"httpfs", sourceInstance, "sqlite", targetInstance)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating: %v", err.Error())
	}
	return sourceInstance.Close()
}
