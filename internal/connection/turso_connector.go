package connection

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/sriram15/progressor-todo-app/internal"
	libsql "github.com/tursodatabase/go-libsql"
)

// TursoConnector implements the DBConnector interface for Turso.
type TursoConnector struct {
	dbURL         string
	authToken     string
	loggableURL   string
	encryptionKey string
}

// NewTursoConnector creates a new TursoConnector.
func NewTursoConnector(dbURL, authToken, encryptionKey string) DBConnector {
	return &TursoConnector{dbURL: dbURL, authToken: authToken, encryptionKey: encryptionKey}
}

// Connect establishes a connection to the Turso database in embedded replica mode only.
func (tc *TursoConnector) Connect() (*sql.DB, string, error) {

	if tc.dbURL == "" || tc.authToken == "" || tc.encryptionKey == "" {
		return nil, DBTypeTurso, fmt.Errorf("Turso DB URL and Auth Token and EncryptionKey must be provided")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, DBTypeTurso, fmt.Errorf("could not get user home directory: %w", err)
	}

	appDir := filepath.Join(homeDir, fmt.Sprintf(".%s", internal.APP_NAME))
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, DBTypeTurso, fmt.Errorf("could not create app directory at %s: %w", appDir, err)
	}

	// Create a unique replica name based on the database URL to avoid conflicts
	sanitizedURL := strings.ReplaceAll(tc.dbURL, "://", "_")
	sanitizedURL = strings.ReplaceAll(sanitizedURL, ".", "_")
	sanitizedURL = strings.ReplaceAll(sanitizedURL, "/", "_")
	dbReplicaName := fmt.Sprintf("replica-%s.db", sanitizedURL)
	dbPath := filepath.Join(appDir, dbReplicaName)

	fmt.Printf("Creating local Turso replica at: %s\n", dbPath)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, tc.dbURL,
		libsql.WithAuthToken(tc.authToken),
		libsql.WithEncryption(tc.encryptionKey),
	)
	if err != nil {
		return nil, DBTypeTurso, fmt.Errorf("failed to create Turso connector: %w", err)
	}

	tursoDb := sql.OpenDB(connector)

	if err := tursoDb.Ping(); err != nil {
		tursoDb.Close()
		return nil, DBTypeTurso, fmt.Errorf("failed to ping Turso database: %w", err)
	}

	tc.loggableURL = tc.dbURL // Log the actual remote URL
	return tursoDb, DBTypeTurso, nil
}

// Migrate applies database migrations for Turso.
func (tc *TursoConnector) Migrate(db *sql.DB, dbType string) error {
	fmt.Printf("Applying migrations to Turso DB with dialect: %s\n", "sqlite3")
	return runGooseMigrations(db, goose.Dialect("sqlite3"))
}

func (tc *TursoConnector) GetDBInfo() (string, string) {
	return "TURSO", tc.loggableURL
}
