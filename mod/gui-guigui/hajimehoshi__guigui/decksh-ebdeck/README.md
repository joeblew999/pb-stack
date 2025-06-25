# Decksh-Ebdeck Integration

This project integrates [decksh](https://github.com/ajstarks/decksh) (a domain-specific language for presentations) with [ebdeck](https://github.com/ajstarks/ebcanvas/tree/main/ebdeck) (an Ebitengine-based presentation renderer) and provides a GUI interface using [guigui](https://github.com/hajimehoshi/guigui).

## Overview

- **decksh**: A DSL for creating presentations that generates deck XML markup
- **ebdeck**: A 2D canvas API built on Ebitengine for rendering presentations
- **guigui**: A GUI framework for Go built with Ebitengine
- **Integration**: Bridge decksh scripts to ebdeck rendering with a modern GUI

## Features

- ✅ Process decksh (.dsh) files and generate deck XML
- ✅ Parse deck XML into structured presentation data
- ✅ Render presentations using Ebitengine-based canvas
- ✅ GUI presentation viewer with navigation controls (when dependencies allow)
- ✅ Keyboard navigation (arrow keys, page up/down, space, home, end)
- ✅ Export presentations to PNG files
- ✅ Support for basic decksh elements (text, shapes, colors)
- ✅ Modern build automation with Taskfile and Makefile
- ✅ Go workspace integration for clean module management
- ✅ Comprehensive testing and validation

## Quick Start

### Prerequisites

1. **Go 1.23+** - Make sure you have Go installed
2. **Task** (optional) - For modern build automation: [Install Task](https://taskfile.dev/installation/)

### Setup

#### Using Taskfile (Recommended)
```bash
cd decksh-ebdeck
task dev-setup    # Complete setup including decksh installation
task run-hello    # Run hello world example
```

#### Manual Setup (Alternative)
```bash
cd decksh-ebdeck
go install github.com/ajstarks/decksh/cmd/decksh@latest
go mod tidy
go build -o bin/decksh-ebdeck ./cmd/decksh-ebdeck
```

## Usage

### Using Taskfile (Recommended)

```bash
# Getting started
task                   # Show available tasks and project variables
task help              # Show detailed help and usage information
task dev-setup         # Complete development setup

# Quick commands
task build             # Build the application
task run-hello         # Run hello world example
task run-demo          # Run demo example
task render-hello      # Render hello world to PNG
task info              # Show project information

# Development workflow
task test              # Run tests
task validate          # Validate decksh example files
task fmt               # Format Go code

# Debug and troubleshooting
task debug-hello       # Show decksh XML output for hello.dsh
task deps-check        # Check if dependencies are installed
```



### Command Line Interface

```bash
# View presentation info
./bin/decksh-ebdeck -input examples/hello.dsh

# Launch GUI presentation viewer (when dependencies allow)
./bin/decksh-ebdeck -input examples/hello.dsh -gui

# Render to PNG file
./bin/decksh-ebdeck -input examples/hello.dsh -output presentation.png

# Custom window size
./bin/decksh-ebdeck -input examples/demo.dsh -gui -width 1400 -height 900

# Verbose output
./bin/decksh-ebdeck -input examples/demo.dsh -verbose
```

### GUI Controls

- **Arrow Keys**: Navigate between slides (left/right)
- **Page Up/Down**: Navigate between slides
- **Space**: Next slide
- **Home**: Go to first slide
- **End**: Go to last slide
- **ESC**: Exit presentation
- **Mouse**: Click Previous/Next buttons

## Project Structure

```
decksh-ebdeck/
├── cmd/decksh-ebdeck/          # Main CLI application
│   └── main.go
├── pkg/
│   ├── bridge/                 # decksh to ebdeck bridge logic
│   │   ├── decksh.go          # decksh processing
│   │   └── presentation.go    # presentation data structures
│   ├── renderer/              # ebdeck rendering wrapper
│   │   └── ebcanvas.go        # canvas rendering implementation
│   └── gui/                   # guigui-based presentation viewer
│       └── presentation.go    # GUI application
├── examples/                  # Sample decksh scripts
│   ├── hello.dsh             # Simple hello world presentation
│   └── demo.dsh              # Feature demonstration
├── go.mod                    # Go module (uses workspace)
├── Taskfile.yml              # Modern build automation
└── README.md                 # This file
```

## Supported Decksh Elements

Currently supported elements:

- **Text**: `text`, `ctext`, `etext` (left, center, right aligned)
- **Shapes**: `rect`, `circle`
- **Lines**: `line`
- **Colors**: Named colors, hex colors (#rrggbb), RGB colors (rgb(r,g,b))
- **Slides**: Background and foreground colors

### Planned Support

- Images (`image`, `cimage`)
- More shapes (`ellipse`, `polygon`, `arc`)
- Text formatting (fonts, styles)
- Gradients and advanced colors
- Animations and transitions

## Example Decksh Files

### Hello World (`examples/hello.dsh`)

```decksh
deck
    slide "white" "black"
        ctext "Hello, World!" 50 70 8
        ctext "Welcome to decksh-ebdeck integration" 50 50 4
        circle 50 15 5 "blue"
    eslide
edeck
```

### Demo Presentation (`examples/demo.dsh`)

See `examples/demo.dsh` for a more comprehensive example showing various decksh features.

## Development

### Adding New Element Support

1. Add parsing logic in `pkg/bridge/decksh.go`
2. Add rendering logic in `pkg/renderer/ebcanvas.go`
3. Test with example decksh files

### Dependencies

- [decksh](https://github.com/ajstarks/decksh) - DSL for presentations
- [ebcanvas](https://github.com/ajstarks/ebcanvas) - 2D canvas API
- [guigui](https://github.com/hajimehoshi/guigui) - GUI framework
- [Ebitengine](https://github.com/hajimehoshi/ebiten) - 2D game engine

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Submit a pull request

## License

This project follows the same license as the integrated components. See individual component licenses for details.

## Build Automation

This project uses **Taskfile** for modern, cross-platform build automation:

### Taskfile Features
- **Modern**: YAML-based configuration with better cross-platform support
- **Task Dependencies**: Automatic dependency resolution
- **Built-in Help**: `task help` shows all available tasks
- **Clean Syntax**: Easy to read and maintain

### Available Tasks

| Task | Description |
|------|-------------|
| `build` | Build the application |
| `deps` | Install dependencies including decksh |
| `run-hello` | Run hello world example |
| `run-demo` | Run demo example |
| `test` | Run tests |
| `validate` | Validate decksh example files |
| `render-hello` | Render hello world to PNG |
| `clean` | Clean build artifacts |
| `info` | Show project information |
| `dev-setup` | Complete development setup |

## Troubleshooting

### Common Issues

1. **"decksh command not found"**
   - Quick fix: `task deps`
   - Manual: `go install github.com/ajstarks/decksh/cmd/decksh@latest`
   - Make sure `$GOPATH/bin` is in your `$PATH`

2. **"Failed to parse deck XML"**
   - Check your decksh syntax
   - Debug: `task debug-hello` or `decksh < yourfile.dsh`
   - Validate: `task validate`

3. **"GUI not launching"**
   - GUI requires network connectivity for dependencies
   - Use PNG rendering instead: `task render-hello`
   - Try running with `-verbose` flag for more information

4. **Task not found**
   - Install Task: [taskfile.dev/installation](https://taskfile.dev/installation/)
   - Alternative: Use manual Go commands

## Examples

Try the included examples:

### Using Taskfile (Recommended)
```bash
# Quick start
task run-hello         # Run hello world example
task run-demo          # Run demo example
task render-hello      # Render hello world to PNG
task render-demo       # Render demo to PNG

# Debug and validate
task debug-hello       # Show decksh XML output
task validate          # Validate all examples
```

### Using Direct Commands
```bash
# Simple hello world
./bin/decksh-ebdeck -input examples/hello.dsh

# Feature demonstration
./bin/decksh-ebdeck -input examples/demo.dsh

# Render to PNG
./bin/decksh-ebdeck -input examples/hello.dsh -output hello.png
```


