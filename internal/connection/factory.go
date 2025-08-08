package connection

import (
	"fmt"
	"os"
	"sync"
)

// DATABASE_NAME is the default name for the SQLite database file.
const DATABASE_NAME = "progressor.db"

var (
	connector DBConnector
	manager   *DBManager
	once      sync.Once
)

// InitDB determines the database type from environment variables,
// creates the appropriate connector, connects, applies migrations,
// and initializes the DBManager.
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

		manager = NewDBManager(connectedDB)
		fmt.Println("Database ready.")
	})
	return dbErr
}

// GetManager returns the singleton DBManager instance.
// It will panic if InitDB has not been called successfully.
func GetManager() *DBManager {
	if manager == nil {
		panic("DBManager not initialized. Call connection.InitDB() first.")
	}
	return manager
}

func GetDBInfo() (string, string) {
	if connector == nil {
		return "Unknown", "Connector not initialized"
	}
	return connector.GetDBInfo()
}
