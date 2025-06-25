#!/bin/bash
# GUI Robot Launcher
APP_DIR="$(dirname "$0")/GUI Robot.app"
if [ -d "$APP_DIR" ]; then
    open "$APP_DIR"
else
    echo "GUI Robot.app not found!"
    exit 1
fi
