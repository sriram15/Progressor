package connection

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/sriram15/progressor-todo-app/internal"
	"github.com/sriram15/progressor-todo-app/internal/utils"
	_ "github.com/tursodatabase/go-libsql" // Driver for libSQL
)

// SQLiteConnector handles SQLite database connections.
type SQLiteConnector struct {
	profileName string
}

// NewSQLiteConnector creates a new SQLiteConnector for a specific profile.
// If profileName is empty or "default", it will use the original database name.
func NewSQLiteConnector(profileName string) *SQLiteConnector {
	return &SQLiteConnector{profileName: profileName}
}

func (c *SQLiteConnector) getDBPath() (string, error) {
	appDir, err := utils.GetAppDir()
	if err != nil {
		return "", fmt.Errorf("failed to get app directory: %w", err)
	}
	
	dbName := internal.DATABASE_NAME
	if c.profileName != "" && strings.ToLower(c.profileName) != "default" {
		dbName = fmt.Sprintf("progressor-%s.db", c.profileName)
	}
	
	return filepath.Join(appDir, dbName), nil
}

// Connect establishes a connection to the SQLite database.
func (c *SQLiteConnector) Connect() (*sql.DB, string, error) {
	dbPath, err := c.getDBPath()
	if err != nil {
		return nil, "", err
	}

	// The connection string for local libSQL/SQLite is just the path.
	// We prepend "file:" for local paths if not already a URL.
	dsn := dbPath
	if !strings.HasPrefix(dsn, "file:") && !strings.HasPrefix(dsn, "http:") && !strings.HasPrefix(dsn, "https:") {
		// Forcing file protocol for local paths with go-libsql
		// On Windows, filepath.ToSlash is important for converting backslashes
		dsn = "file:" + filepath.ToSlash(dbPath)
	}

	fmt.Printf("Attempting to connect to SQLite DB at: %s\n", dsn)

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, "", fmt.Errorf("failed to open SQLite database at %s: %w", dsn, err)
	}

	if err = db.Ping(); err != nil {
		return nil, "", fmt.Errorf("failed to ping SQLite database at %s: %w", dsn, err)
	}

	fmt.Println("Successfully connected to SQLite database.")
	return db, DBTypeSQLite, nil
}

// Migrate applies database migrations for SQLite.
func (sc *SQLiteConnector) Migrate(db *sql.DB, dbType string) error {
	fmt.Printf("Applying migrations to SQLite DB with dialect: %s\n", "sqlite3")
	return runGooseMigrations(db, goose.Dialect("sqlite3"))
}

func (sc *SQLiteConnector) GetDBInfo() (string, string) {
	dbPath, err := sc.getDBPath()
	if err != nil {
		return "LOCAL", "Unable to determine SQLite path"
	}
	return "LOCAL", dbPath
}