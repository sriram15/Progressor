package connection

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/sriram15/progressor-todo-app/internal/database"
)

// DATABASE_NAME is the default name for the SQLite database file.
const DATABASE_NAME = "progressor.db"

var connector DBConnector

// OpenDB determines the database type from environment variables,
// creates the appropriate connector, connects, and applies migrations.
func OpenDB() (*sql.DB, error) {
	dbType := os.Getenv("DB_TYPE")

	switch dbType {
	case DBTypeSQLite, "sqlite3", "": // Default to SQLite
		fmt.Println("Using SQLite database connector.")
		connector = NewSQLiteConnector()
	case DBTypeTurso:
		fmt.Println("Using Turso database connector.")
		connector = NewTursoConnector()
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s. Supported types are 'sqlite' or 'turso'", dbType)
	}

	db, actualDBType, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database (%s): %w", dbType, err)
	}

	fmt.Printf("Successfully established database connection (Type: %s).\n", actualDBType)

	if err := connector.Migrate(db, actualDBType); err != nil {
		// It's important to close the DB if migration fails to prevent leaks.
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations for %s: %w", actualDBType, err)
	}

	fmt.Println("Database ready.")
	return db, nil
}

func GetDBQuery() (*database.Queries, error) {

	db, err := OpenDB()
	if err != nil {
		return nil, err
	}

	queries := database.New(db)
	return queries, nil
}

func GetDBInfo() (string, string) {
	return connector.GetDBInfo()
}
