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
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

    # Check if installation was successful
    if command_exists brew; then
        echo "Homebrew installed successfully."
        # Add Homebrew to the PATH if it's not already there.
        if [[ ":$PATH:" != *":/opt/homebrew/bin:"* ]]; then
            echo 'Adding Homebrew to PATH in ~/.zshrc and ~/.bashrc...'
            echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zshrc
            echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.bashrc
            source ~/.zshrc
            source ~/.bashrc
        fi
    else
        echo "Error: Failed to install Homebrew."
        exit 1
    fi
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
        echo "Installing VS Code extensions from extensions.txt..."
        while read -r extension; do
            if [[ ! "$extension" =~ ^#.*$ && -n "$extension" ]]; then # Skip comments and empty lines
                echo "Installing extension: $extension"
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