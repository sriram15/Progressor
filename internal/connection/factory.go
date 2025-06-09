package connection

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

// DATABASE_NAME is the default name for the SQLite database file.
const DATABASE_NAME = "progressor.db"

var (
	connector DBConnector
	db        *sql.DB
	once      sync.Once
	dbMutex   sync.Mutex
)

// InitDB determines the database type from environment variables,
// creates the appropriate connector, connects, and applies migrations.
func InitDB() error {
	var dbErr error
	once.Do(func() {
		dbType := os.Getenv("DB_TYPE")

		switch dbType {
		case DBTypeSQLite, "sqlite3", "": // Default to SQLite
			fmt.Println("Using SQLite database connector.")
			connector = NewSQLiteConnector()
		case DBTypeTurso:
			fmt.Println("Using Turso database connector.")
			connector = NewTursoConnector()
		default:
			dbErr = fmt.Errorf("unsupported DB_TYPE: %s. Supported types are 'sqlite' or 'turso'", dbType)
			return
		}

		connectedDB, actualDBType, err := connector.Connect()
		if err != nil {
			dbErr = fmt.Errorf("failed to connect to database (%s): %w", dbType, err)
			return
		}

		fmt.Printf("Successfully established database connection (Type: %s).\n", actualDBType)

		if err := connector.Migrate(connectedDB, actualDBType); err != nil {
			// It's important to close the DB if migration fails to prevent leaks.
			connectedDB.Close()
			dbErr = fmt.Errorf("failed to apply migrations for %s: %w", actualDBType, err)
			return
		}

		db = connectedDB
		fmt.Println("Database ready.")
	})
	return dbErr
}

// GetDB checks if the current db connection is alive, and if not, reopens it.
func GetDB() (*sql.DB, func()) {
	dbMutex.Lock()
	return db, func() {
		dbMutex.Unlock()
	}
}

func GetDBInfo() (string, string) {
	return connector.GetDBInfo()
}
