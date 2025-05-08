package connection

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

// DBType constants
const (
	DBTypeSQLite = "sqlite"
	DBTypeTurso  = "turso"
)

// MIGRATION_DIRECTORY is relative to the project root when the binary runs.
const MIGRATION_DIRECTORY = "./internal/database/migrations"

// DBConnector defines the interface for database operations.
type DBConnector interface {
	Connect() (*sql.DB, string, error)
	Migrate(db *sql.DB, dbType string) error
	GetDBInfo() (string, string)
}

// runGooseMigrations is a helper function to apply migrations.
func runGooseMigrations(db *sql.DB, dbDialect goose.Dialect) error {
	fmt.Printf("Looking for migrations in: %s (relative to CWD)\n", MIGRATION_DIRECTORY)

	migrationFS := os.DirFS(MIGRATION_DIRECTORY)

	provider, err := goose.NewProvider(dbDialect, db, migrationFS)
	if err != nil {
		return fmt.Errorf("error creating goose provider: %w", err)
	}

	_, err = provider.Up(context.Background())
	if err != nil {
		if err.Error() == "no migrations found" {
			return fmt.Errorf("failed to apply migrations: %w. Ensure migration files exist in %s", err, MIGRATION_DIRECTORY)
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migration(s) applied successfully.")
	return nil
}
