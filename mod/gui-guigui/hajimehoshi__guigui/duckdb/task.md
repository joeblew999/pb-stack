# Taskfile System

Built with [Task](https://github.com/go-task/task/) ([v3.44.0](https://github.com/go-task/task/releases/tag/v3.44.0))

## Overview

This taskfile system is designed to provide a consistent and efficient development environment for both human developers and AI assistants working in VS Code. It implements common development workflows and utilities while keeping everything contained within the project directory.

## Core Features

### 1. Local Shell Completions
- Generates completion scripts for bash, zsh, and fish shells
- All completion files are stored in the local `./completions` directory
- No global system pollution
- Easy to enable/disable per project

```bash
# Generate completion scripts
task task:completion:create

# Get instructions for enabling completions
task task:completion:enable
```

### 2. Process Management
- List running processes: `task task:process:list`
- Show all process details: `task task:process:list:all`
- Search for specific processes: `task task:process:search -- <name>`
- Kill processes by name or PID: `task task:process:kill -- <name|pid>`

### 3. File Management
- Basic file listing: `task task:file:list`
- More features planned

### 4. Binary Management (hatch-server example)
- Local binary installation in `.bin` directory
- Cross-platform support (handles .exe on Windows)
- Version information and health checks
- Process control (start/stop)

## Integration with VS Code

The taskfile system is designed to work seamlessly with VS Code, providing tools that both human developers and AI assistants (like GitHub Copilot) can use effectively. While some features overlap with Model Context Protocol (MCP) systems, having these tools directly available in VS Code improves the development workflow.

## Project Structure

```
.
├── completions/        # Shell completion scripts
├── .bin/              # Local binary installations
├── Taskfile.yml       # Main task file
├── task.taskfile.yml  # Core utilities (process, file management)
└── hatch.taskfile.yml # Binary-specific tasks
```

## Common Tasks

### Core Utilities
```bash
# List available tasks
task --list-all

# Show process information
task task:process:list
task task:process:list:all
task task:process:search -- <name>

# File operations
task task:file:list
```

### Binary Management
```bash
# Install/update hatch-server
task hatch:dep

# Check binary location and info
task hatch:dep:which

# Start/stop server
task hatch:serve
task hatch:run:kill
```

## Future Plans

- Expand file management capabilities
- Develop a proper Taskfile MCP implementation
- Add more developer utilities
- Improve cross-platform support

