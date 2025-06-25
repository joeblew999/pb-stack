# xtask

**A unified Go binary that embeds Task + all cross-platform development tools**

## 🎯 Problem Statement

Task is not fully cross-platform - Taskfiles often contain Unix-specific commands like `2>/dev/null`, `pkill`, `command -v` that break on Windows. This creates a "works on my machine" problem and prevents true cross-platform development workflows.

## 🚀 Solution

A single Go binary that embeds Task + all necessary cross-platform tools, creating a unified, all-Go development command that works identically on Windows, macOS, and Linux.

## Conventions

.bin for binaries

.wasm for wasm binaries. 

.data for data, like ./.data/nats

## 🔧 Core Architecture

### Hybrid CLI/Daemon with Embedded NATS JetStream

**YES! xtask operates as both CLI tool AND server with embedded NATS JetStream**

#### CLI Mode (instant execution)
- Binary detection (`xtask which go`)
- File downloads (`xtask got https://...`)
- Cross-platform silent execution (`xtask silent cmd`)
- Cross-platform process management (`xtask kill-port 8080`)

#### Server Mode (persistent service)
- Start with embedded NATS JetStream (`xtask --server`)
- HTTP/gRPC API server (:8080)
- Web interface (/web)
- NATS coordination (:4222)
- Cluster connectivity
- Real-time coordination

### Smart CLI Client
- CLI automatically uses best execution method
- Supports local execution, server-specific execution, and cluster coordination

### API-Driven Everything
- RESTful HTTP API endpoints
- WebSocket for real-time updates
- Web interface with visual dashboard, NATS monitoring, and real-time logs

## 📦 Embedded Tools

### Core Tools
- **go-task/task** - Task runner (core functionality)
- **go-which** - Cross-platform binary detection
- **got** - Cross-platform file downloads
- **task-tui** - Interactive Task interface
- **process-compose** - Process orchestration
- **go-git** - Pure Go git operations

### System Libraries
- **gopsutil** - Cross-platform system operations
- **archives** - Archive/compression handling
- **afero** - Filesystem abstraction

### Custom Utilities
- **Silent execution** - Cross-platform `2>/dev/null`
- **Port management** - Cross-platform `pkill`
- **File operations** - Cross-platform `cp`, `tree`
- **Network utilities** - Port checking, health checks

## 🎯 Cross-Platform Utilities

### Shell Operation Replacements
- Silent execution (2>/dev/null equivalent)
- DevNull (>/dev/null 2>&1 equivalent)
- Cross-platform port killing
- Port waiting with timeout
- Parallel command execution
- Command retry with count

### File Operations
- Cross-platform tree view
- Cross-platform file copy
- Archive extraction and creation

### Network Operations
- Enhanced file downloads
- HTTP health checks
- Port scanning

## 🔧 Taskfile Integration

### Base Configuration
- Unified tool access through Taskfile variables

## 🌟 Server Mode Features

### Embedded NATS JetStream
```bash
# Start server with embedded NATS
xtask --server
# ✅ NATS server running on :4222
# ✅ JetStream enabled for persistence
# ✅ Cluster-ready for global coordination
# ✅ No external dependencies
```

### HTTP/gRPC API Server
```bash
# RESTful API
curl http://localhost:8080/api/v1/tools/which/go
curl -X POST http://localhost:8080/api/v1/tools/got \
  -d '{"url": "https://example.com/file.zip", "output": "/tmp/file.zip"}'

# WebSocket for real-time updates
ws://localhost:8080/ws

# Web interface
http://localhost:8080/web
```

### Smart CLI Routing
```bash
# CLI automatically chooses best execution method:
xtask which go
# 1. Local server (fastest)
# 2. Remote cluster (coordinated)
# 3. Local execution (fallback)

# Environment variables control behavior:
export XTASK_SERVER_URL="http://build-server:8080"
export XTASK_CLUSTER_URL="nats://global.cluster:4222"
```

### Real-Time Coordination
- **Live task execution** - See commands running across nodes
- **Resource monitoring** - CPU, memory, disk usage
- **Log streaming** - Real-time output from all nodes
- **Health monitoring** - Node status and availability

## 🌟 Benefits

### Cross-Platform Consistency
- ✅ Same commands work on Windows, macOS, Linux
- ✅ No more platform-specific shell scripts
- ✅ No more "works on my machine" issues

### All-Go Toolchain
- ✅ No Python dependencies
- ✅ No OS-specific utilities required
- ✅ Single binary distribution
- ✅ Perfect for containers and CI/CD

### Hybrid Architecture Benefits
- ✅ **CLI simplicity** - Works without server for basic tasks
- ✅ **Server power** - Advanced coordination when needed
- ✅ **API integration** - Perfect for automation and tooling
- ✅ **Real-time collaboration** - Team coordination via NATS

### Professional Development Platform
- ✅ Enterprise-ready
- ✅ Reproducible builds
- ✅ Version-controlled toolchain
- ✅ AI-friendly automation
- ✅ **Global coordination** - Distributed team support
- ✅ **Embedded infrastructure** - No external dependencies

## 🚀 Implementation Phases

### Phase 1: Core Integration
- Embed Task + go-which + got
- Basic cross-platform utilities (silent, devnull)
- Multi-binary support via symlinks

### Phase 2: Enhanced Features
- Add process-compose integration
- File operations (copy, archive, tree)
- Network utilities (port checking, health checks)

### Phase 3: Advanced Capabilities
- Interactive TUI components
- Git operations integration
- Template and code generation
- Advanced process management

### Phase 4: Ecosystem Integration
- Full Taskfile integration
- NATS JetStream coordination
- AI agent compatibility
- Community adoption

**xtask completes the vision of a sophisticated, cross-platform development ecosystem that actually works everywhere!** 🚀






