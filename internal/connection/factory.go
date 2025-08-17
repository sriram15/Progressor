package connection

import (
	"fmt"

	"github.com/sriram15/progressor-todo-app/internal/profile"
)

// DATABASE_NAME is the default name for the SQLite database file.
const DATABASE_NAME = "progressor.db"

// InitDB is deprecated. Use NewManagerForProfile instead.
// This function is part of the old singleton pattern and will be removed.
func InitDB() error {
	return fmt.Errorf("InitDB is deprecated and should not be called")
}

// GetManager is deprecated. The DBManager should be retrieved from the active AppSession.
// This function is part of the old singleton pattern and will be removed.
func GetManager() *DBManager {
	panic("GetManager is deprecated. The DBManager should be retrieved from the active AppSession.")
}

// GetDBInfo is deprecated.
// This function is part of the old singleton pattern and will be removed.
func GetDBInfo() (string, string) {
	return "Unknown", "GetDBInfo is deprecated."
}

// NewManagerForProfile creates a new DBManager for a given profile.
// It handles the creation of the correct connector (SQLite or Turso),
// fetches the token for Turso, connects, and runs migrations.
func NewManagerForProfile(p *profile.Profile) (*DBManager, error) {
	var connector DBConnector

	fmt.Print("Creating database manager for profile: %s\n", p.DBType)

	switch p.DBType {
	case profile.DBTypeSQLite, "": // Default to SQLite
		fmt.Printf("Creating SQLite connector for profile: %s\n", p.Name)
		connector = NewSQLiteConnector(p.Name)
	case profile.DBTypeTurso:
		fmt.Printf("Creating Turso connector for profile: %s\n", p.Name)
		// The AuthTokenKey from the profile is the key for the keychain
		token, err := profile.GetToken(p.AuthTokenKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get token for profile %s: %w", p.Name, err)
		}
		connector = NewTursoConnector(p.DBUrl, token)
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s", p.DBType)
	}

	connectedDB, actualDBType, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database for profile %s: %w", p.Name, err)
	}

	if err := connector.Migrate(connectedDB, actualDBType); err != nil {
		connectedDB.Close()
		return nil, fmt.Errorf("failed to apply migrations for profile %s: %w", p.Name, err)
	}

	manager := NewDBManager(connectedDB)
	fmt.Println("Database manager created successfully.")
	return manager, nil
}
