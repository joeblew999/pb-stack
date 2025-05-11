#!/usr/bin/env bash

# boot.sh
# Shell script to attempt to install specified packages using Homebrew.
# This is intended for macOS and Linux systems where Homebrew is used.

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Homebrew is installed
if command_exists brew; then
    echo "Homebrew found. Proceeding with installation attempts..."
else
    echo "Homebrew (brew) not found. Attempting to install..."
    # Attempt to install Homebrew
    # The Homebrew installer will guide the user on adding brew to their PATH for future sessions.
    if /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"; then
        echo "Homebrew installation script finished."
        echo "IMPORTANT: Please follow any instructions output by the Homebrew installer regarding PATH configuration for your shell to ensure 'brew' is available in new terminal sessions."
    else
        echo "Error: Homebrew installation script failed."
        exit 1
    fi
fi

# Ensure Homebrew environment is configured for the *current script execution*.
# This is crucial if Homebrew was just installed or if the user's shell environment isn't fully set up yet.
echo ""
echo "Configuring Homebrew environment for this script session..."
BREW_EXECUTABLE=$(command -v brew) # Try to find brew in PATH first (might be set by user or previous run)

if [ -z "$BREW_EXECUTABLE" ] || ! [ -x "$BREW_EXECUTABLE" ]; then
    # If not found via command -v, check standard locations
    if [ -x "/opt/homebrew/bin/brew" ]; then # Apple Silicon
        BREW_EXECUTABLE="/opt/homebrew/bin/brew"
    elif [ -x "/usr/local/bin/brew" ]; then # Intel Macs or Linuxbrew default
        BREW_EXECUTABLE="/usr/local/bin/brew"
    fi
fi

if [ -x "$BREW_EXECUTABLE" ]; then
    eval "$($BREW_EXECUTABLE shellenv)"
    echo "Homebrew environment configured for this script."
else
    echo "Warning: Could not find 'brew' executable to configure environment for this script. Subsequent brew commands might fail."
    # Depending on strictness, you might want to exit 1 here if brew is essential.
fi

echo ""
echo "---------------------------------------------------------------------"
echo "Setup script finished basic environment checks (e.g., Homebrew)."
echo "Package installations and VS Code extensions are now primarily"
echo "managed by the Go application based on config.json."
echo "This script can be used for other pre-setup tasks or system"
echo "configurations not covered by the JSON."
echo "---------------------------------------------------------------------"
# Add any other non-package, non-vscode-extension setup tasks here.
# For example, creating directories, setting specific environment variables for the session, etc.