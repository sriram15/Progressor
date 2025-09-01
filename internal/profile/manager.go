package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sriram15/progressor-todo-app/internal"
	"github.com/sriram15/progressor-todo-app/internal/utils"
)

// ProfileManager handles reading from and writing to the profiles.json configuration file.
type ProfileManager struct {
	configPath string
}

// NewManager creates a new ProfileManager.
// It determines the correct application config path and ensures the directory exists.
func NewManager() (*ProfileManager, error) {
	appDir, err := utils.GetAppDir()
	if err != nil {
		return nil, err
	}

	log.Printf("Using application directory: %s", appDir)
	return &ProfileManager{
		configPath: filepath.Join(appDir, "profiles.json"),
	}, err
}

// LoadConfig reads the profiles.json file from the config path.
// If the file doesn't exist, it returns a new, empty Config struct.
func (m *ProfileManager) LoadConfig() (*Config, error) {
	bytes, err := os.ReadFile(m.configPath)

	if err != nil {
		// If the file simply doesn't exist, it's not an error.
		// This is expected on first launch.
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Profiles: []Profile{}}, nil
		}
		return nil, err // Otherwise, it's a real read error.
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig serializes the given Config struct to JSON and writes it to the config path.
func (m *ProfileManager) SaveConfig(config *Config) error {
	bytes, err := json.MarshalIndent(config, "", "  ") // Use indentation for readability
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, bytes, 0644) // 0644 is standard file permissions
}

// GetProfiles loads the configuration and returns all profiles.
func (m *ProfileManager) GetProfiles() ([]Profile, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load profiles config: %w", err)
	}
	return config.Profiles, nil
}

// GetProfile retrieves a profile by its ID.
func (m *ProfileManager) GetProfile(profileID string) (*Profile, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load profiles config: %w", err)
	}

	for _, p := range config.Profiles {
		if p.ID == profileID {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("profile with ID %s not found", profileID)
}

// CreateProfile creates a new profile, saves it to the config, and stores the Turso token if applicable.
func (m *ProfileManager) CreateProfile(newProfile Profile, tursoToken, encryptionKeyPath string) (*Profile, error) {
	config, err := m.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load profiles config: %w", err)
	}

	// Generate a unique ID for the new profile
	newProfile.ID = uuid.New().String()

	// Store Turso token if applicable
	if newProfile.DBType == DBTypeTurso {
		if tursoToken == "" {
			return nil, errors.New("Turso token is required for Turso database type")
		}

		if encryptionKeyPath == "" {
			return nil, errors.New("Encryption key path is required for Turso database type")
		}

		// Use the profile ID as the key for the token in the keyring
		newProfile.AuthTokenKey = newProfile.ID
		err := StoreToken(newProfile.AuthTokenKey, tursoToken)
		if err != nil {
			return nil, fmt.Errorf("failed to store Turso token: %w", err)
		}

		err = StoreToken(newProfile.EncryptionKeyPath, encryptionKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to store Encryption key path: %w", err)
		}
	}

	config.Profiles = append(config.Profiles, newProfile)

	if err := m.SaveConfig(config); err != nil {
		return nil, fmt.Errorf("failed to save profiles config: %w", err)
	}

	return &newProfile, nil
}

// DeleteProfile removes a profile by its ID.
func (m *ProfileManager) DeleteProfile(profileID string) error {
	config, err := m.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load profiles config: %w", err)
	}

	found := false
	for i, p := range config.Profiles {
		if p.ID == profileID {
			// Remove the profile from the slice
			config.Profiles = append(config.Profiles[:i], config.Profiles[i+1:]...)
			found = true
			// If it's a Turso profile, delete its token from the keyring
			if p.DBType == DBTypeTurso && p.AuthTokenKey != "" {
				if err := DeleteToken(p.AuthTokenKey); err != nil {
					log.Printf("Warning: Failed to delete Turso token for profile %s: %v", p.Name, err)
				}
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("profile with ID %s not found", profileID)
	}

	if err := m.SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save profiles config: %w", err)
	}
	return nil
}

// GetProfileDBPath returns the database file path for a given SQLite profile.
func (m *ProfileManager) GetProfileDBPath(profileID string) (string, error) {
	p, err := m.GetProfile(profileID)
	if err != nil {
		return "", err
	}
	if p.DBType != DBTypeSQLite {
		return "", fmt.Errorf("profile %s is not a SQLite database", p.Name)
	}
	// This logic is duplicated from sqlite_connector.go.
	// A better approach would be to centralize DB path generation.
	// For now, we'll replicate it.
	appDir, err := utils.GetAppDir()
	if err != nil {
		return "", fmt.Errorf("failed to get app directory: %w", err)
	}
	dbName := internal.DATABASE_NAME                          // Use the constant from connection package
	if p.Name != "" && strings.ToLower(p.Name) != "default" { // Assuming profile name is used for db name
		dbName = fmt.Sprintf("progressor-%s.db", p.Name)
	}
	return filepath.Join(appDir, dbName), nil
}

