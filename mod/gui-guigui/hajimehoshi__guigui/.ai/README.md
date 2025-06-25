# decksh-ebdeck Integration

This project integrates decksh (a domain-specific language for presentations) with ebdeck (an Ebitengine-based presentation renderer).

## URLS

- https://github.com/hajimehoshi/guigui
- https://github.com/ajstarks/decksh
- https://github.com/ajstarks/ebcanvas

## Overview

- **decksh**: DSL for creating presentations, generates deck XML markup
- **ebdeck**: 2D canvas API built on Ebitengine for rendering presentations
- **Integration**: Bridge decksh scripts to ebdeck rendering with guigui GUI

## Architecture

1. **decksh processor** - Parse .dsh files and generate deck XML
2. **ebdeck renderer** - Render presentations using Ebitengine
3. **guigui wrapper** - Provide GUI interface for presentation viewing
4. **CLI tool** - Command-line interface for batch processing

## Components

- `decksh-ebdeck/cmd/decksh-ebdeck/` - Main CLI application
- `decksh-ebdeck/pkg/bridge/` - decksh to ebdeck bridge logic
- `decksh-ebdeck/pkg/renderer/` - ebdeck rendering wrapper
- `decksh-ebdeck/pkg/gui/` - guigui-based presentation viewer
- `decksh-ebdeck/examples/` - Sample decksh scripts

## Key Features Implemented

✅ **decksh Processing**: Calls external `decksh` command to convert .dsh files to deck XML
✅ **XML Parsing**: Parses deck XML into structured Go data types
✅ **Ebitengine Rendering**: Renders presentations using Ebitengine's 2D graphics
✅ **guigui GUI**: Modern GUI with slide navigation, keyboard controls
✅ **CLI Interface**: Command-line tool with multiple output options
✅ **Color Support**: Named colors, hex colors, RGB colors
✅ **Element Support**: Text, rectangles, circles, lines
✅ **Export**: Render presentations to PNG files
✅ **Go Workspace**: Uses go.work for clean module management

## Usage Examples

```bash
# Setup (one time)
cd decksh-ebdeck
./setup.sh

# View presentation in GUI
./bin/decksh-ebdeck -input examples/hello.dsh -gui

# Render to PNG
./bin/decksh-ebdeck -input examples/demo.dsh -output demo.png

# Using Makefile
make run-hello    # Launch hello world
make run-demo     # Launch demo
make deps         # Install dependencies
```

## Project Structure

```
decksh-ebdeck/
├── cmd/decksh-ebdeck/main.go          # CLI application
├── pkg/
│   ├── bridge/                        # decksh ↔ ebdeck bridge
│   │   ├── decksh.go                 # Process .dsh files
│   │   ├── presentation.go           # Data structures
│   │   └── bridge_test.go            # Tests
│   ├── renderer/ebcanvas.go          # Ebitengine rendering
│   └── gui/presentation.go           # guigui GUI viewer
├── examples/
│   ├── hello.dsh                     # Simple example
│   └── demo.dsh                      # Feature demo
├── go.mod                            # Go module (uses workspace)
├── Makefile                          # Build automation
├── setup.sh                          # Setup script
└── README.md                         # Documentation
```



