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

# Uninstall VS Code extensions not listed in extensions.txt
if command_exists code && [ -f "extensions.txt" ]; then
    echo ""
    echo "Uninstalling VS Code extensions not listed in extensions.txt ..."

    # Create an array of extensions from extensions.txt
    readarray -t desired_extensions < extensions.txt

    # Get a list of currently installed extensions
    installed_extensions=$(code --list-extensions)

    # Iterate through installed extensions and uninstall those not in the desired list
    while IFS= read -r installed_extension; do
        found=false
        for desired_extension in "${desired_extensions[@]}"; do
            if [[ "$installed_extension" == "$desired_extension" ]]; then
                found=true
                break
            fi
        done
        if ! $found; then
            echo "Uninstalling extension: $installed_extension..."
            code --uninstall-extension "$installed_extension"
        fi
    done <<< "$installed_extensions"
    echo "VS Code extension uninstallation finished."
else
    echo ""
    if ! command_exists code; then
        echo "VS Code (code) not found. Skipping extension uninstallation."
    elif [ ! -f "extensions.txt" ]; then
        echo "extensions.txt not found. Skipping extension uninstallation."
    fi
fi