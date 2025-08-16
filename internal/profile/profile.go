package profile

// DBType defines the type of database connection.
// Use constants DBTypeSQLite or DBTypeTurso.
type DBType string

const (
	DBTypeSQLite DBType = "sqlite"
	DBTypeTurso  DBType = "turso"
)

// Config represents the root of the profiles.json file,
// containing a list of all user profiles.
type Config struct {
	Profiles []Profile `json:"profiles"`
}

// Profile defines the configuration for a single user profile.
// It holds the necessary information to connect to a specific database.
type Profile struct {
	// ID is the unique identifier for the profile (e.g., a UUID).
	ID string `json:"id"`

	// Name is the user-friendly name for the profile (e.g., "Work", "Personal").
	Name string `json:"name"`

	// DBType specifies the database backend, either "sqlite" or "turso".
	DBType DBType `json:"dbType"`

	// DBPath is the absolute path to the database file. Used only when DBType is "sqlite".
	DBPath string `json:"dbPath,omitempty"`

	// DBUrl is the remote database URL. Used only when DBType is "turso".
	DBUrl string `json:"dbUrl,omitempty"`

	// AuthTokenKey is the key used to look up the Turso auth token in the OS keychain.
	// This is NOT the token itself. Used only when DBType is "turso".
	AuthTokenKey string `json:"authTokenKey,omitempty"`
}
