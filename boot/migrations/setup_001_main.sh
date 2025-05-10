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
    "visual-studio-code" # Added VS Code Cask
)

echo ""
echo "Attempting to install the following Homebrew packages:"
for pkg_name in "${packages_to_install[@]}"; do
    echo "  - ${pkg_name}"
done
echo "------------------------------------"

for pkg in "${packages_to_install[@]}"; do
    if [[ "$pkg" == "go" ]] && command_exists go; then
        echo ""
        echo "Go already detected. Skipping installation."
        continue
    fi

    if [[ "$pkg" == "bun" ]] && command_exists bun; then
        echo ""
        echo "Bun already detected. Skipping installation."
        continue
    fi

    if [[ "$pkg" == "git" ]] && command_exists git; then
        echo ""
        echo "Git already detected. Skipping installation."
        continue
    fi

    if [[ "$pkg" == "openssh" ]] && command_exists ssh; then
        echo ""
        echo "OpenSSH already detected. Skipping installation."
        continue
    fi

    if [[ "$pkg" == "coreutils" ]] && command_exists which; then
        echo ""
        echo "coreutils (which) already detected. Skipping installation."
        continue
    fi

    # VS Code is a cask, and 'code' command might exist if installed manually or via other means
    if [[ "$pkg" == "visual-studio-code" ]] && command_exists code; then
        echo ""
        echo "Visual Studio Code (code command) already detected. Skipping installation via Homebrew."
        continue
    fi

    echo ""
    echo "Attempting to install '$pkg'..." && brew install "$pkg" && echo "Successfully installed '$pkg'." || echo "Could not install '$pkg'. It might not be available or another error occurred."
done

echo ""
echo "------------------------------------"
echo "Installation process finished."
echo "Note: Some packages (like git, ssh, or which) might have system-provided versions that take precedence over Homebrew-installed versions."

    # Install VS Code extensions if VS Code is installed and extensions.txt exists
    if command_exists code && [ -f "extensions.txt" ]; then
        echo ""
        echo "Installing VS Code extensions from extensions.txt ..."
        while read -r extension; do
            if [[ ! "$extension" =~ ^#.*$ && -n "$extension" ]]; then # Skip comments and empty lines
                code --list-extensions | grep -q "$extension" && echo "Extension $extension already installed, skipping." && continue
                echo "Installing extension: $extension..."
                code --install-extension "$extension"
            fi
        done < extensions.txt
        echo "VS Code extension installation finished."
    else
        echo ""
        if ! command_exists code; then
            echo "VS Code (code) not found. Skipping extension installation."
        elif [ ! -f "extensions.txt" ]; then
            echo "extensions.txt not found. Skipping extension installation."
        fi
    fi

echo ""
echo "Listing installed VS Code extensions:"
code --list-extensions