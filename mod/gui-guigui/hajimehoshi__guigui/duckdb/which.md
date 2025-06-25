# go-which

A cross-platform Go implementation of the Unix `which` command. Works on Linux, macOS, and Windows.

Part of a Model Context Protocol (MCP) implementation for integrating command discovery into AI-assisted workflows.

## Tasks

Run tasks using: `task which:<taskname>`

- `which:dep` - Install dependencies and build the tool
- `which:clean` - Clean build artifacts
- `which:which` - Find binary path
- `which:test` - Run tests

## Usage

```bash
go-which [command]
```

Finds the location of a command in your PATH.
