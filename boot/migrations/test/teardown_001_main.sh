#!/usr/bin/env bash

# Shell script to attempt to uninstall specified packages using Homebrew.
# This is intended for macOS and Linux systems where Homebrew is used.

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if Homebrew is installed
if ! command_exists brew; then
    echo "Error: Homebrew (brew) is not installed."
    echo "Please install Homebrew first to manage these packages."
    echo "Installation instructions can be found at: https://brew.sh/"
    exit 1
fi

echo "Homebrew found. Proceeding with uninstallation attempts..."

echo ""
echo "---------------------------------------------------------------------"
echo "Teardown script finished basic environment checks (e.g., Homebrew)."
echo "Package uninstallation and VS Code extensions are now primarily"
echo "managed by the Go application based on config.json."
echo "This script can be used for other post-teardown tasks or system"
echo "cleanup not covered by the JSON."
echo "---------------------------------------------------------------------"
# Add any other non-package, non-vscode-extension teardown tasks here.