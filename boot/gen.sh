#!/bin/bash
#use the vscode cli to generate the vscode extensions into the extensions.txt file, so that the boot files can use it.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
code --list-extensions > "$SCRIPT_DIR/extensions.txt"
echo "VS Code extensions list saved to $SCRIPT_DIR/extensions.txt"