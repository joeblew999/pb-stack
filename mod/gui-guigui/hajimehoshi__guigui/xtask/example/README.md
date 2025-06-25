# xtask Example Project

This example demonstrates the basic capabilities of xtask - a unified Go binary that embeds Task + all cross-platform development tools.

## 🚀 Quick Start

```bash
# Show available tasks
task

# Run the full demonstration
task demo

# Or try individual features
task build-xtask    # Build xtask binary
task which          # Test binary detection
task download       # Test file downloads
task health-check   # Test HTTP health checks
task tree           # Show directory tree
task build          # Build example app
task run            # Run example app
```

## 🌟 Server Mode

xtask can run as a server with embedded NATS JetStream for team coordination:

```bash
# Start server
task server-start

# Check status
task server-status

# Test distributed features
task distributed-test

# Stop server
task server-stop
```

## 🔧 xtask Features Demonstrated

### CLI Tools (Cross-Platform)
- **`xtask which`** - Binary detection (replaces Unix `which`)
- **`xtask got`** - File downloads (replaces `curl`/`wget`)
- **`xtask silent`** - Silent execution (replaces `2>/dev/null`)
- **`xtask kill-port`** - Port management (replaces `pkill`)
- **`xtask wait-for-port`** - Port waiting
- **`xtask tree`** - Directory trees (replaces Unix `tree`)
- **`xtask health-check`** - HTTP health checks

### Server Features
- **HTTP API** - RESTful endpoints for all tools
- **WebSocket** - Real-time command streaming
- **NATS JetStream** - Embedded message streaming
- **Web Interface** - Browser-based dashboard
- **Cluster Coordination** - Multi-node synchronization

## 🌍 Cross-Platform Benefits

This example works identically on:
- ✅ **Windows** - No PowerShell/CMD differences
- ✅ **macOS** - No BSD vs GNU tool differences  
- ✅ **Linux** - No distribution-specific issues

All through a single Go binary with no external dependencies!

## 📁 Project Structure

```
example/
├── Taskfile.yml     # Task definitions using xtask
├── main.go          # Simple Go application
├── go.mod           # Go module definition
└── README.md        # This file
```

## 🎯 What This Demonstrates

1. **Unified Toolchain** - One binary for all development tools
2. **Cross-Platform** - Same commands work everywhere
3. **Task Integration** - Seamless integration with Task runner
4. **Server Mode** - Optional server mode for team coordination
5. **API-Driven** - Everything available via HTTP API
6. **Real-Time** - WebSocket support for live updates

## 🔗 Related

- [xtask README](../README.md) - Full xtask documentation
- [NATS Integration](../NATS.md) - Distributed coordination details
- [Hosting Guide](../HOSTING.md) - Global deployment strategy
