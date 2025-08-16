package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sriram15/progressor-todo-app/internal"
)

// GetAppDir returns the absolute path to the application's configuration and data directory.
// It creates the directory if it doesn't exist.
// The location is determined by the OS:
// - Windows: %APPDATA%\progressor
// - macOS: $HOME/Library/Application Support/progressor
// - Linux: $XDG_DATA_HOME/progressor or $HOME/.local/share/progressor
func GetAppDir() (string, error) {
	var appDir string

	switch runtime.GOOS {
	case "windows":
		appDir = os.Getenv("APPDATA")
		if appDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		appDir = filepath.Join(appDir, internal.APP_NAME)

	case "darwin": // macOS
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			return "", fmt.Errorf("HOME environment variable not set")
		}
		appDir = filepath.Join(homeDir, "Library", "Application Support", internal.APP_NAME)

	case "linux":
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome == "" {
			homeDir := os.Getenv("HOME")
			if homeDir == "" {
				return "", fmt.Errorf("neither XDG_DATA_HOME nor HOME environment variables are set")
			}
			xdgDataHome = filepath.Join(homeDir, ".local", "share")
		}
		appDir = filepath.Join(xdgDataHome, internal.APP_NAME)

	default: // Fallback for other OSes
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Ensure the app-specific directory exists
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create app directory %s: %w", appDir, err)
	}

	return appDir, nil
}
