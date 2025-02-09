package internal

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
	_ "github.com/tursodatabase/go-libsql"
)

const MIGRATION_PATH = "database/migrations"
const DATABASE_NAME = "progressor.db"

//go:embed database/migrations/*.sql
var embedMigrations embed.FS

// getDatabasePath returns the path to the SQLite database file.
func GetDatabasePath(prependUrl string) (string, error) {

	// Check if DATBASE_PATH env variable is set
	if databasePathEnv := os.Getenv("DATABASE_PATH"); databasePathEnv != "" {
		databasePathEnv = filepath.Join(databasePathEnv, DATABASE_NAME)
		return databasePathEnv, nil
	}

	var appDir string
	// Determine the appropriate directory based on OS
	switch runtime.GOOS {
	case "windows":
		appDir = os.Getenv("APPDATA")
	case "darwin": // macOS
		appDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		appDir = filepath.Join(os.Getenv("HOME"), ".config")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	// Create a directory for your app
	appDir = filepath.Join(prependUrl, appDir, APP_NAME)
	if err := os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", err
	}

	// Define the database file path
	dbPath := filepath.Join(appDir, DATABASE_NAME)
	return dbPath, nil
}

func OpenDB() (*sql.DB, error) {
	dbPath, err := GetDatabasePath("file:")
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %v", err)
	}

	// db, err := sql.Open("sqlite3", dbPath)
	fmt.Println(dbPath)
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		return nil, fmt.Errorf("sql: failed to open DB: %v", err)
	}

	goose.SetBaseFS(embedMigrations)
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, os.DirFS("./internal/database/migrations"))
	if err != nil {
		return nil, fmt.Errorf("error creating goose provider: ", err)
	}

	results, err := provider.Up(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to create migration: %v", err)
	}

	for _, r := range results {
		fmt.Printf("OK   %s (%s)\n", r.Source.Path, r.Duration)
	}

	fmt.Println("Migration created successfully.")
	return db, nil
}

func CloseDb(db *sql.DB) {
	if err := db.Close(); err != nil {
		fmt.Println("Failed to close DB:", err)
	}
}
