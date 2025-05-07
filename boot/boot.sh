#!/usr/bin/env bash

# Shell script to attempt to install specified packages using Homebrew.
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

echo "Homebrew found. Proceeding with installation attempts..."

# List of packages to attempt to uninstall.
# These are analogous to the winget packages:
# - Oven-sh.Bun (winget) -> bun (brew)
# - Git.Git (winget) -> git (brew)
# - GoLang.Go (winget) -> go (brew)
# - Microsoft.OpenSSH.Preview (winget) -> openssh (brew's version)
# - GnuWin32.Which (winget) -> which (brew's version, often coreutils)
packages_to_install=(
    "bun"
    "git"
    "go"
    "openssh"
    "coreutils"
)

echo ""
echo "Attempting to install the following Homebrew packages:"
for pkg_name in "${packages_to_install[@]}"; do
    echo "  - ${pkg_name}"
done
echo "------------------------------------"

for pkg in "${packages_to_install[@]}"; do
    echo ""
    echo "Attempting to install '$pkg'..."
    if brew install "$pkg"; then
        echo "Successfully installed '$pkg'."
    else
        # brew uninstall exits with non-zero if not installed or other error
        echo "Could not uninstall '$pkg'. It might not be installed via Homebrew, or another error occurred."
    fi
done

echo ""
echo "------------------------------------"
echo "Installation process finished."
echo "Note: Some packages (like git, ssh, or which) might have system-provided versions that take precedence over Homebrew-installed versions."