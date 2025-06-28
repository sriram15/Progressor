#!/bin/bash

# Get the absolute path of the directory where the script is located
# This makes the script runnable from anywhere.
SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)

# --- Configuration ---
# Set the project directory to where the script is.
PROJECT_DIR="$SCRIPT_DIR"

# --- Environment Variables ---
# Load Turso credentials from a .env file in the project root if it exists.
# This is a best practice to avoid hardcoding secrets.
if [ -f "$PROJECT_DIR/.env" ]; then
  export $(grep -v '^#' "$PROJECT_DIR/.env" | xargs)
fi

# Ensure Turso variables are set, otherwise exit.
if [ -z "$TURSO_DB_PATH" ] || [ -z "$TURSO_AUTH_TOKEN" ]; then
    echo "Error: TURSO_DB_PATH and TURSO_AUTH_TOKEN must be set in .env file."
    # You can use a system notification here for better visibility on startup
    # On macOS: osascript -e 'display notification "Turso credentials not set" with title "App Startup Error"'
    # On Linux: notify-send "App Startup Error" "Turso credentials not set"
    exit 1
fi

# --- Execution ---
echo "Starting Wails dev server for project at $PROJECT_DIR..."

# Navigate to the project directory
cd "$PROJECT_DIR"

# Launch the Wails development server.
# The 'wails3 dev' command automatically handles building the Go source.
# We add 'nohup' and '&' to ensure it runs in the background and doesn't
# die if the parent shell that started it closes.
nohup wails3 dev &

echo "App started in background."