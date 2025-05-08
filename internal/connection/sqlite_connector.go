package connection

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/sriram15/progressor-todo-app/internal"
	_ "github.com/tursodatabase/go-libsql" // Driver for libSQL
)

// SQLiteConnector handles SQLite database connections.
type SQLiteConnector struct {
	// No specific fields needed for now, as path is determined at connection time.
}

// NewSQLiteConnector creates a new SQLiteConnector.
func NewSQLiteConnector() *SQLiteConnector {
	return &SQLiteConnector{}
}

// Connect establishes a connection to the SQLite database.
func (c *SQLiteConnector) Connect() (*sql.DB, string, error) {
	dbPath, err := getLocalOSPath()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get local SQLite path: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, "", fmt.Errorf("failed to create database directory %s: %w", dir, err)
	}

	// The connection string for local libSQL/SQLite is just the path.
	// For URLs (like Turso), it might be different.
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

// GetLocalOSPath returns the path to the SQLite database file.
// It prioritizes LOCAL_DATABASE_PATH env var, then XDG directories, then current dir.
func getLocalOSPath() (string, error) {
	// Check if LOCAL_DATABASE_PATH env variable is set
	if databasePathEnv := os.Getenv("LOCAL_DATABASE_PATH"); databasePathEnv != "" {
		// If LOCAL_DATABASE_PATH is a directory, append DATABASE_NAME. Otherwise, use as is.
		info, err := os.Stat(databasePathEnv)
		if err == nil && info.IsDir() {
			return filepath.Join(databasePathEnv, DATABASE_NAME), nil
		}
		// If it's a file (or error checking Stat failed, assume it's a full path)
		return databasePathEnv, nil
	}

	// Try XDG data directory (Linux, macOS)
	var appDir string
	switch runtime.GOOS {
	case "windows":
		appDir = os.Getenv("APPDATA")
		if appDir == "" {
			fmt.Println("APPDATA not set, falling back to current directory for SQLite DB.")
			appDir = "." // Fallback
		} else {
			appDir = filepath.Join(appDir, internal.APP_NAME) // Store in AppData/YourAppName
		}
	case "darwin": // macOS
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			return "", fmt.Errorf("HOME not set, cannot determine default SQLite path")
		}
		appDir = filepath.Join(homeDir, "Library", "Application Support", internal.APP_NAME)
	case "linux":
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome == "" {
			homeDir := os.Getenv("HOME")
			if homeDir == "" {
				return "", fmt.Errorf("XDG_DATA_HOME and HOME not set, cannot determine default SQLite path")
			}
			xdgDataHome = filepath.Join(homeDir, ".local", "share")
		}
		appDir = filepath.Join(xdgDataHome, internal.APP_NAME)
	default: // Fallback for other OSes
		fmt.Printf("Unsupported OS (%s) for XDG path, falling back to current directory for SQLite DB.\n", runtime.GOOS)
		appDir = "." // Fallback to current directory
	}

	// Ensure the app-specific directory exists
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory %s: %w", appDir, err)
	}
	return filepath.Join(appDir, DATABASE_NAME), nil
}

func (sc *SQLiteConnector) GetDBInfo() (string, string) {
	dbPath, err := getLocalOSPath()
	if err != nil {
		return "LOCAL", "Unable to determine SQLite path"
	}
	return "LOCAL", dbPath
}
