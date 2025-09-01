package connection

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// DBType constants
const (
	DBTypeSQLite  = "sqlite"
	DBTypeTurso   = "turso"
	MIGRATION_DIR = "migrations"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// DBConnector defines the interface for database operations.
type DBConnector interface {
	Connect() (*sql.DB, string, error)
	Migrate(db *sql.DB, dbType string) error
	GetDBInfo() (string, string)
}

// runGooseMigrations is a helper function to apply migrations.
func runGooseMigrations(db *sql.DB, dbDialect goose.Dialect) error {

	fsys, err := fs.Sub(embedMigrations, MIGRATION_DIR)
	if err != nil {
		return fmt.Errorf("failed to access embedded migrations: %w", err)
	}

	provider, err := goose.NewProvider(dbDialect, db, fsys)
	if err != nil {
		return fmt.Errorf("error creating goose provider: %w", err)
	}

	_, err = provider.Up(context.Background())
	if err != nil {
		if err.Error() == "no migrations found" {
			return fmt.Errorf("failed to apply migrations: %w.", err)
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migration(s) applied successfully.")
	return nil
}
