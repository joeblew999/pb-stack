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

# List of packages to attempt to uninstall.
# These are analogous to the winget packages:
# - Oven-sh.Bun (winget) -> bun (brew)
# - Git.Git (winget) -> git (brew)
# - GoLang.Go (winget) -> go (brew)
# - Microsoft.OpenSSH.Preview (winget) -> openssh (brew's version)
# - GnuWin32.Which (winget) -> which (brew's version, often gnu-which)
packages_to_uninstall=(
    "bun"
    "git"
    "go"
    "openssh"
    "which"
)

echo ""
echo "Attempting to uninstall the following Homebrew packages:"
for pkg_name in "${packages_to_uninstall[@]}"; do
    echo "  - ${pkg_name}"
done
echo "------------------------------------"

for pkg in "${packages_to_uninstall[@]}"; do
    echo ""
    echo "Attempting to uninstall '$pkg'..."
    if brew uninstall "$pkg"; then
        echo "Successfully uninstalled '$pkg'."
    else
        # brew uninstall exits with non-zero if not installed or other error
        echo "Could not uninstall '$pkg'. It might not be installed via Homebrew, or another error occurred."
    fi
done

echo ""
echo "------------------------------------"
echo "Uninstallation process finished."
echo "Note: System-provided versions of some tools (like git, ssh, or which) might still be present if they were not installed/managed by Homebrew."