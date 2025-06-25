#!/bin/bash

# DataStar WASM Rod Testing Script
# This script helps run Rod tests with proper setup

set -e

echo "🤖 DataStar WASM Rod Testing"
echo "============================"

# Check if applications are running
check_service() {
    local url=$1
    local name=$2
    
    if curl -s "$url" > /dev/null 2>&1; then
        echo "✅ $name is running at $url"
        return 0
    else
        echo "❌ $name is not running at $url"
        return 1
    fi
}

# Function to start services if needed
start_services() {
    echo "🚀 Starting DataStar services..."
    
    # Start server in background
    echo "📡 Starting server..."
    cd .. && task server &
    SERVER_PID=$!
    
    # Start WASM service worker
    echo "🌐 Starting WASM service worker..."
    task wasm &
    WASM_PID=$!
    
    # Start Todo WASM
    echo "📝 Starting Todo WASM..."
    task wasm-todo &
    TODO_PID=$!
    
    # Wait for services to start
    echo "⏳ Waiting for services to start..."
    sleep 5
    
    # Store PIDs for cleanup
    echo "$SERVER_PID" > /tmp/ds-wasm-server.pid
    echo "$WASM_PID" > /tmp/ds-wasm-wasm.pid
    echo "$TODO_PID" > /tmp/ds-wasm-todo.pid
}

# Function to stop services
stop_services() {
    echo "🛑 Stopping services..."
    
    if [ -f /tmp/ds-wasm-server.pid ]; then
        kill $(cat /tmp/ds-wasm-server.pid) 2>/dev/null || true
        rm -f /tmp/ds-wasm-server.pid
    fi
    
    if [ -f /tmp/ds-wasm-wasm.pid ]; then
        kill $(cat /tmp/ds-wasm-wasm.pid) 2>/dev/null || true
        rm -f /tmp/ds-wasm-wasm.pid
    fi
    
    if [ -f /tmp/ds-wasm-todo.pid ]; then
        kill $(cat /tmp/ds-wasm-todo.pid) 2>/dev/null || true
        rm -f /tmp/ds-wasm-todo.pid
    fi
    
    # Also use task kill as backup
    cd .. && task kill 2>/dev/null || true
}

# Cleanup on exit
trap stop_services EXIT

# Check if services are already running
echo "🔍 Checking if services are running..."

SERVER_RUNNING=false
WASM_RUNNING=false
TODO_RUNNING=false

if check_service "http://localhost:8081/health" "Server"; then
    SERVER_RUNNING=true
fi

if check_service "http://localhost:8082" "WASM Service Worker"; then
    WASM_RUNNING=true
fi

if check_service "http://localhost:8083" "Todo WASM"; then
    TODO_RUNNING=true
fi

# Start services if not running
if [ "$SERVER_RUNNING" = false ] || [ "$WASM_RUNNING" = false ] || [ "$TODO_RUNNING" = false ]; then
    echo "⚠️  Some services are not running. Starting them..."
    start_services
    
    # Re-check services
    echo "🔍 Re-checking services..."
    sleep 3
    
    check_service "http://localhost:8081/health" "Server" || {
        echo "❌ Failed to start server"
        exit 1
    }
    
    check_service "http://localhost:8082" "WASM Service Worker" || {
        echo "❌ Failed to start WASM service worker"
        exit 1
    }
    
    check_service "http://localhost:8083" "Todo WASM" || {
        echo "❌ Failed to start Todo WASM"
        exit 1
    }
else
    echo "✅ All services are already running"
fi

echo ""
echo "🧪 Running Rod tests..."

# Determine test mode
HEADLESS=${HEADLESS:-true}
if [ "$1" = "--visible" ] || [ "$1" = "-v" ]; then
    HEADLESS=false
    echo "👁️  Running tests with visible browser"
else
    echo "🕶️  Running tests in headless mode"
fi

# Create screenshots directory
mkdir -p screenshots

# Run the tests
if [ "$HEADLESS" = true ]; then
    go test -v ./tests -timeout 60s
else
    go test -v ./tests -timeout 60s -args -headless=false
fi

echo ""
echo "✅ Rod tests completed successfully!"
echo "📸 Screenshots saved in: tests/screenshots/"
echo "📋 Test logs available above"
