#!/bin/bash
#use the vscode cli to generate the vscode extensions into the extensions.txt file, so that the boot files can use it.

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
EXTENSIONS_FILE="$SCRIPT_DIR/extensions.txt"

# Get the current list of extensions
current_extensions=$(code --list-extensions)

# Save the current list to a temporary file
echo "$current_extensions" > "$EXTENSIONS_FILE.tmp"

# Compare with the existing extensions.txt file, if it exists
if [ -f "$EXTENSIONS_FILE" ]; then
    new_extensions=$(comm -23 <(sort "$EXTENSIONS_FILE.tmp") <(sort "$EXTENSIONS_FILE"))
    deleted_extensions=$(comm -13 <(sort "$EXTENSIONS_FILE.tmp") <(sort "$EXTENSIONS_FILE"))

    echo "New extensions:"
    echo "$new_extensions"
    echo "Deleted extensions:"
    echo "$deleted_extensions"
fi

# Overwrite the extensions.txt file with the current list
mv "$EXTENSIONS_FILE.tmp" "$EXTENSIONS_FILE"
echo "VS Code extensions list saved to $EXTENSIONS_FILE"