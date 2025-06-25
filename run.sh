#!/bin/bash

# GUI Robot - Build and Run Script
# Alternative to Taskfile for when Task has issues

set -e

APP_NAME="gui-robot"
BIN_DIR="bin"
CMD_DIR="./cmd/gui-robot"

function show_help() {
    echo "GUI Robot - AI-Powered GUI Automation"
    echo "======================================"
    echo ""
    echo "Usage: ./run.sh [command]"
    echo ""
    echo "Commands:"
    echo "  build            - Build the application"
    echo "  test             - Run tests"
    echo "  clean            - Clean build artifacts"
    echo "  demo             - Run automation demo"
    echo "  screen-info      - Get screen information"
    echo "  interactive      - Start interactive mode"
    echo "  help             - Show this help"
    echo ""
    echo "Examples:"
    echo "  ./run.sh build"
    echo "  ./run.sh demo"
    echo "  ./run.sh screen-info"
}

function build() {
    echo "Building GUI Robot..."
    mkdir -p $BIN_DIR
    go build -o $BIN_DIR/$APP_NAME $CMD_DIR
    echo "✓ Build complete: $BIN_DIR/$APP_NAME"
}

function test() {
    echo "Running tests..."
    go test -v ./...
}

function clean() {
    echo "Cleaning build artifacts..."
    rm -rf $BIN_DIR
    rm -f *.png
    echo "✓ Clean complete"
}

function demo() {
    build
    echo "Running GUI Robot demo..."
    ./$BIN_DIR/$APP_NAME -command get_screen_info
    echo "✓ Screen info retrieved"
    ./$BIN_DIR/$APP_NAME -command move_mouse -params '{"x":100,"y":100}'
    echo "✓ Mouse moved"
    ./$BIN_DIR/$APP_NAME -command type -params '{"text":"Hello from GUI Robot!"}'
    echo "✓ Text typed"
    ./$BIN_DIR/$APP_NAME -command wait -params '{"duration":1}'
    echo "✓ Demo completed!"
}

function screen_info() {
    build
    ./$BIN_DIR/$APP_NAME -command get_screen_info
}

function interactive() {
    build
    ./$BIN_DIR/$APP_NAME -interactive
}

# Main command handling
case "${1:-help}" in
    build)
        build
        ;;
    test)
        test
        ;;
    clean)
        clean
        ;;
    demo)
        demo
        ;;
    screen-info)
        screen_info
        ;;
    interactive)
        interactive
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
