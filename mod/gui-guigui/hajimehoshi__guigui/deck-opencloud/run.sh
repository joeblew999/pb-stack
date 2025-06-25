#!/bin/bash

# OpenCloud - Task Runner Script
# This script ensures we use the correct Taskfile

set -e

TASKFILE="./Taskfile.yml"

if [ ! -f "$TASKFILE" ]; then
    echo "❌ Taskfile.yml not found in current directory"
    exit 1
fi

# If no arguments provided, show default task
if [ $# -eq 0 ]; then
    task --taskfile="$TASKFILE"
else
    # Pass all arguments to task
    task --taskfile="$TASKFILE" "$@"
fi
