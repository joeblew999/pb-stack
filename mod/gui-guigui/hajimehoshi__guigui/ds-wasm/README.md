# ds-wasm

This project demonstrates DataStar running as WASM in the browser with **multiple examples**:

1. **Hello World** - Basic DataStar patterns (server vs WASM comparison)
2. **Todo App** - Advanced DataStar features (forms, lists, persistence)

## The example we are replicating

**Official DataStar Repository**: https://github.com/starfederation/datastar
**Go Module**: `github.com/starfederation/datastar v1.0.0-beta.11`
**JavaScript CDN**: `https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-beta.11/bundles/datastar.js`
**Reference Example**: https://github.com/starfederation/datastar/tree/develop/examples/go/hello-world

This project replicates the official DataStar hello-world example using the official StarFederation DataStar module, demonstrating both server and WASM modes of the same application. 

We try to get the backend golang and the golang compiled to wasm to run in the browser to be as isomorphic as possible. 

Daisy GUI is https://github.com/CoreyCole/datastarui. Might use this later.

## WASM Web Workers approach

https://github.com/magodo/go-wasmww looks interesting. It looks like it is designed so that can run other golang wasm that is Late Bound and Late Loaded ? DataStar signals can help here too when we late load wasm ? 

Looks like an example Golang wasm code designed to be run inside go-wasmww needs to be written to use STD IO ? 

https://github.com/magodo/go-wasmww/blob/main/examples/shared/worker/main.go is an example worker.

https://github.com/nlepage/go-wasm-http-server helps with running perhaps ? . Its golang module is github.com/nlepage/go-wasm-http-server/v2




## Testing approach

Rod ( https://github.com/go-rod/rod ) seems to be the way Datastar can be tested in the actual browser ? 


## 🎯 Dual Mode DataStar

Run the same DataStar application in two different ways:

1. **Server Mode**: Traditional Go web server (http://localhost:8081)
2. **WASM Mode**: Go compiled to WebAssembly running in browser (http://localhost:8082)

## 🚀 Quickstart

```bash
# Build and run hello-world dual modes
task dev

# Or run individual examples
task server        # Server mode (port 8081)
task wasm           # Hello-world WASM (port 8082)
task wasm-todo      # Todo WASM example (port 8083)
```

**🌐 Open your browser to:**
- **Server Mode**: http://localhost:8081 (Hello World with Go HTTP server)
- **WASM Service Worker**: http://localhost:8082 (Service Worker with 3 URLs + SSE)
- **WASM Todo**: http://localhost:8083 (Todo App with Go WASM)

## 🔧 Individual Commands

```bash
# Build
task build-server     # Build server binary
task build-wasm       # Build hello-world WASM
task build-wasm-todo  # Build todo WASM
task build            # Build all binaries

# Run
task server          # Run server mode (port 8081)
task wasm            # Run WASM service worker (port 8082)
task wasm-todo       # Run todo WASM (port 8083)
task dev             # Run dual-mode with process orchestration
task dev-simple      # Run dual-mode (simple approach)

# Setup
task tools-install   # Install cross-platform tools (go-which, got)
task caddy-install   # Install Caddy web server
task tools-check     # Check if tools are installed
task caddy-check     # Check if Caddy is installed

# Process Management
task kill            # Kill all project processes (cross-platform)
task kill-server     # Kill ds-server processes
task kill-wasm       # Kill caddy processes on WASM port

# Test
task test-server     # Test server endpoints
task test-wasm       # Test WASM in browser (manual)
task test-rod        # Test with Rod browser automation (visible)
task test-rod-headless # Test with Rod (headless mode)
task test-all        # Run all automated tests
```

## 🌟 Technologies

- **DataStar**: Hypermedia framework (https://github.com/starfederation/datastar)
- **WebAssembly**: Go compiled to WASM for browser execution
- **Service Workers**: For browser integration (planned)
- **File System Access API**: For file operations (planned via https://github.com/tractordev/toolkit-go)

## 📁 Project Structure

```
ds-wasm/
├── Taskfile.yml           # Build & test automation
├── Processfile.yml        # Process orchestration
├── cmd/
│   ├── server/main.go     # Server mode implementation
│   ├── wasm/main.go       # Hello-world WASM implementation
│   ├── wasm-todo/main.go  # Todo WASM implementation
│   └── pkill/main.go      # Cross-platform process killer
├── web/
│   ├── index.html         # Hello-world WASM frontend
│   ├── wasm/              # Hello-world WASM files
│   └── todo/
│       ├── index.html     # Todo WASM frontend
│       └── wasm/          # Todo WASM files
└── bin/                   # Generated binaries
```

## 🎮 Features Demonstrated

### ✅ Hello World WASM Service Worker
- **Service Worker Pattern** - WASM acts as server inside browser
- **3 URL Routes** - `/`, `/hello`, `/status` endpoints
- **SSE Communication** - Server-Sent Events simulation
- **Real-time Updates** - Live connection status and message counts
- **DataStar Integration** - Full reactive UI with state management
- **Build constraints** (`//go:build js && wasm`)

### ✅ Todo WASM Example
- **Form handling** (add new todos)
- **List management** (display, filter todos)
- **State persistence** (localStorage)
- **Interactive UI** (toggle, delete todos)
- **Filtering** (all/active/completed)
- **Statistics** (count active/completed)
- **Pure WASM** (no server required)
- **Build constraints** (`//go:build js && wasm`)

### ✅ Rod Browser Testing
- **Real browser automation** - Tests run in actual browsers (Chrome/Firefox)
- **DataStar validation** - Verifies reactive UI updates work correctly
- **WASM testing** - Ensures WASM modules load and function properly
- **Cross-application testing** - Tests server vs WASM consistency
- **Screenshot capture** - Visual debugging and test documentation
- **Headless/visible modes** - Run tests with or without browser UI

### 🚧 Planned Features
- **Service Worker integration**
- **File System Access API**
- **Offline functionality**
- **Progressive Web App features**

## 🔍 Example Comparison

The same DataStar hello-world example demonstrates:

| Feature | Server Mode | WASM Mode |
|---------|-------------|-----------|
| **Rendering** | Server-side | Client-side |
| **State Management** | DataStar + Server | DataStar + WASM |
| **Network** | HTTP requests | Function calls |
| **Performance** | Network latency | Near-native speed |
| **Offline** | ❌ Requires server | ✅ Runs offline |

## 🧪 Testing

Open both modes side-by-side to compare:
1. **Increment/Decrement**: Pure DataStar reactivity
2. **Server Interaction**: Different backend handling
3. **Time Updates**: Different data sources
4. **Performance**: Compare response times

## WASM Web Workers Approach

### Why Web Workers + WASM is Perfect for DataStar

**WASM Web Workers** provide the ideal architecture for DataStar applications because they solve the fundamental challenges of reactive web applications:

#### **🔄 Non-Blocking Reactivity**
- **DataStar's reactive updates** run in background thread, keeping UI responsive
- **Heavy state computations** don't freeze the user interface
- **Real-time data processing** happens without blocking user interactions

#### **⚡ True Parallel Processing**
- **Multiple DataStar stores** can be managed simultaneously
- **Background data synchronization** while user continues working
- **Concurrent SSE/WebSocket handling** for real-time features

#### **🧠 Complex State Management**
- **Go's powerful type system** handles complex DataStar state logic
- **Business rules and validation** run efficiently in compiled WASM
- **Data transformations** leverage Go's excellent JSON/data handling

#### **📡 Background Communication**
- **SSE connections** managed in worker without blocking UI
- **API calls and data fetching** happen in background
- **WebSocket management** for real-time DataStar updates

### Architecture Benefits for DataStar

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Main Thread   │◄──►│   Web Worker     │◄──►│  WASM Module    │
│                 │    │                  │    │                 │
│ • DataStar UI   │    │ • State Manager  │    │ • Go Logic      │
│ • DOM Updates   │    │ • SSE Handler    │    │ • Business Rules│
│ • User Events   │    │ • Data Sync      │    │ • Validation    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

#### **🎯 Perfect DataStar Use Cases**
1. **Reactive Forms**: Complex validation logic in Go/WASM
2. **Real-time Dashboards**: Background data processing and updates
3. **Collaborative Apps**: Multi-user state synchronization
4. **Offline-First**: Background sync with conflict resolution
5. **Data-Heavy Apps**: Efficient processing of large datasets

#### **🛡️ Production Advantages**
- **Isolation**: Worker crashes don't affect main DataStar UI
- **Security**: WASM provides sandboxed execution environment
- **Performance**: Compiled Go code runs faster than JavaScript
- **Debugging**: Separate worker debugging in DevTools
- **Testing**: Rod can test worker communication patterns

### File System Integration

For file operations, the architecture uses the **File System Access API** via [toolkit-go](https://github.com/tractordev/toolkit-go), enabling:
- **Direct file access** from DataStar applications
- **Drag-and-drop file handling** with reactive UI updates
- **Local file persistence** for offline-first applications

## Cloudflare Hosting

https://github.com/syumai/workers allows us to have our Server running on Cloudflare Workers.

I am not sure yet how SSE will work.  https://github.com/syumai/workers/issues/164 has some info related to how to do it.



## Tooling

This project uses consistent naming: `Taskfile.yml` and `Processfile.yml` for clear parallel structure.

No absolute paths are used in the taskfile or the processfile.

All tooling is golang based and cross platform. No external dependencies like Python, curl, or wget are required. When we install any golang tool, the binary will have .exe on Windows and we account for this in the taskfile.

**Cross-Platform Binary Handling**: The Taskfile automatically handles `.exe` extensions on Windows using the `BIN_EXT` variable. For direct process-compose usage on Windows, manually add `.exe` to binary paths in `Processfile.yml`.

### Build & Development Tools

**Pure Go Toolchain**: All tools are Go-based for maximum cross-platform compatibility. No Python, curl, wget, or other external dependencies required.

**Go 1.24.4** (latest) - WASM support built-in with `wasm_exec.js` at `$(go env GOROOT)/lib/wasm/wasm_exec.js`
- WASM files use proper build constraints: `//go:build js && wasm`
- Ensures WASM code only compiles for the `js/wasm` target
- Prevents accidental compilation for other platforms

**Taskfile** ([go-task/task](https://github.com/go-task/task)) - Build & test automation using `Taskfile.yml`
- Cross-platform task automation
- Build server binary and WASM module
- Run tests against running services

**process-compose** ([F1bonacc1/process-compose](https://github.com/F1bonacc1/process-compose)) - Process orchestration using `Processfile.yml`
- Run binaries with dependencies and health checks
- Manage server and WASM dev server
- CLI docs: https://f1bonacc1.github.io/process-compose/cli/process-compose/

**Caddy** ([caddyserver.com](https://caddyserver.com/)) - WASM development server
- Modern web server with automatic HTTPS and file browsing
- Excellent WASM support with proper MIME types
- Cross-platform installation via `go install github.com/caddyserver/caddy/v2/cmd/caddy@latest`
- Zero configuration needed for static file serving
- Required for WASM development (no fallback needed)

**go-which** ([hairyhenderson/go-which](https://pkg.go.dev/github.com/hairyhenderson/go-which)) - Cross-platform `which` command
- Reliable binary detection across all platforms
- Used for detecting Caddy and other Go tools
- Install: `go install github.com/hairyhenderson/go-which/cmd/which@latest`

**gopsutil v4** ([shirou/gopsutil](https://github.com/shirou/gopsutil)) - Cross-platform process management
- Latest version with improved performance and Windows support
- Reliable process management across all platforms including Windows
- Used for building cross-platform `pkill` utility
- Handles process killing, listing, and monitoring
- Built into project as `bin/pkill` utility

**got** ([melbahja/got](https://github.com/melbahja/got)) - Cross-platform HTTP downloader
- Fast concurrent downloader, faster than curl and wget
- Cross-platform alternative to curl/wget
- Used for testing HTTP endpoints in tasks
- Install: `go install github.com/melbahja/got/cmd/got@latest`

**Rod** ([go-rod/rod](https://github.com/go-rod/rod)) - Browser automation for testing
- High-level driver based on DevTools Protocol
- Perfect for testing DataStar reactive applications
- Real browser testing (Chrome, Firefox, etc.)
- Screenshot capture and visual debugging
- Headless and visible browser modes
- Cross-platform browser automation

### Optional Tools

**mkcert** ([FiloSottile/mkcert](https://github.com/FiloSottile/mkcert)) - HTTPS locally for advanced web features

### Project Structure

```
.gitignore          # Ignores bin/, web/wasm/*.wasm, web/wasm/wasm_exec.js
Taskfile.yml        # Build & test automation (handles .exe on Windows)
Processfile.yml     # Process orchestration
go.mod              # Go dependencies (Official DataStar v1.0.0-beta.11, gopsutil v4)
bin/                # Built binaries (ds-server, pkill + .exe on Windows)
web/wasm/           # WASM build artifacts
cmd/pkill/          # Cross-platform process killer using gopsutil v4
```




