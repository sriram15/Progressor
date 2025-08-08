package connection

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/sriram15/progressor-todo-app/internal"
	libsql "github.com/tursodatabase/go-libsql"
)

// TursoConnector implements the DBConnector interface for Turso.
type TursoConnector struct {
	loggableURL string
}

// NewTursoConnector creates a new TursoConnector.
func NewTursoConnector() DBConnector {
	return &TursoConnector{}
}

// Connect establishes a connection to the Turso database in embedded replica mode only.
func (tc *TursoConnector) Connect() (*sql.DB, string, error) {
	tursoDbPath := os.Getenv("TURSO_DB_PATH")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")
	encryptionKey := os.Getenv("TURSO_ENCRYPTION_KEY")
	dbReplicaName := "progressor-replica.db"

	if tursoDbPath == "" || authToken == "" {
		return nil, DBTypeTurso, fmt.Errorf("TURSO_DB_PATH, TURSO_AUTH_TOKEN must be set in the environment or .env file for embedded replica mode")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error creating using the home directory for replica:", err)
		return nil, DBTypeTurso, err
	}

	dbDir := filepath.Join(homeDir, fmt.Sprintf(".%s", internal.APP_NAME))
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Errorf("Could not create database directory at %s: %v", dbDir, err)
		return nil, DBTypeTurso, err
	}

	dbPath := filepath.Join(dbDir, dbReplicaName)

	fmt.Printf("Created a local temp replica at: %s\n", dbPath)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, tursoDbPath,
		libsql.WithAuthToken(authToken),
		libsql.WithEncryption(encryptionKey),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		return nil, DBTypeTurso, err
	}

	tursoDb := sql.OpenDB(connector)

	if err := tursoDb.Ping(); err != nil {
		return nil, DBTypeTurso, err
	}
	// defer tursoDb.Close()

	tc.loggableURL = tursoDbPath
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
