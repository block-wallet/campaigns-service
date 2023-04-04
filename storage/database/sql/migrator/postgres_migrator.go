package migrator

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations
var Migrations embed.FS

type postgremigrator struct {
	db *sql.DB
}

func NewPostgresMigrator(db *sql.DB) *postgremigrator {
	return &postgremigrator{
		db: db,
	}
}

func (migrator *postgremigrator) Migrate() error {
	sourceInstance, err := httpfs.New(http.FS(Migrations), "migrations")
	if err != nil {
		return fmt.Errorf("invalid source instance, %w", err)
	}
	targetInstance, err := postgres.WithInstance(migrator.db, new(postgres.Config))
	if err != nil {
		return fmt.Errorf("invalid target postgres instance, %w", err)
	}
	m, err := migrate.NewWithInstance(
		"httpfs", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %w", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error migrating: %v", err.Error())
	}
	return sourceInstance.Close()
}
