package connection

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	// Turso driver (libsql)
	"github.com/pressly/goose/v3"
	_ "github.com/tursodatabase/go-libsql"
)

// TursoConnector implements the DBConnector interface for Turso.
type TursoConnector struct {
	loggableURL string
}

// NewTursoConnector creates a new TursoConnector.
func NewTursoConnector() DBConnector {
	return &TursoConnector{}
}

// Connect establishes a connection to the Turso database.
func (tc *TursoConnector) Connect() (*sql.DB, string, error) {
	tursoURL := os.Getenv("TURSO_DB_URL")
	if tursoURL == "" {
		return nil, DBTypeTurso, fmt.Errorf("TURSO_DB_URL is not set in the environment or .env file")
	}

	db, err := sql.Open("libsql", tursoURL)
	if err != nil {
		return nil, DBTypeTurso, fmt.Errorf("sql: failed to open Turso DB: %w", err)
	}

	loggableURL := tursoURL
	if strings.Contains(loggableURL, "authToken=") {
		loggableURL = strings.Split(loggableURL, "authToken=")[0] + "authToken=****"
	}
	tc.loggableURL = loggableURL

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, DBTypeTurso, fmt.Errorf("sql: failed to ping TursoDB. Check URL (%s). Original error: %w", loggableURL, err)
	}
	return db, DBTypeTurso, nil
}

// Migrate applies database migrations for Turso.
func (tc *TursoConnector) Migrate(db *sql.DB, dbType string) error {
	fmt.Printf("Applying migrations to Turso DB with dialect: %s\n", "sqlite3")
	return runGooseMigrations(db, goose.Dialect("sqlite3"))
}

func (tc *TursoConnector) GetDBInfo() (string, string) {
	return "TURSO", tc.loggableURL
}
