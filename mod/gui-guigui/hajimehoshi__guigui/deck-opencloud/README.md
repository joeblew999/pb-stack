# deck-opencloud


This project runs [`opencloud`](https://github.com/opencloud-eu/opencloud) (a Golang collaboration server) 

This gives use a File System that is collaborative, and searchable, and indexed.

## Integration plan

With [`Guigui`](https://github.com/hajimehoshi/guigui) (a Golang GUI framework) with decksh (a domain-specific language for presentations).

With Datastar ( https://github.com/datastar-app/datastar ), we can build a Web GUI for the OpenCloud server that allows users to view, search, and collaborate on documents. 

## 🚀 Quickstart

Get up and running with the clean **Build → Run → Test** workflow:

```bash
# Full development workflow (recommended)
task pc:dev
```

This will:
- 📦 **Build** all binaries (`opencloud-server`, `opencloud-gui`, `opencloud-cli`)
- 📊 **Create** and index sample documents
- 🚀 **Run** the collaboration server with health checks
- 🧪 **Test** all endpoints to verify everything works

**🌐 Open your browser to:** http://localhost:8080

### Alternative Workflows

```bash
# Minimal (just build + run server)
task pc:dev-minimal

# GUI development (build + run server + GUI)
task pc:dev-gui

# Stop all processes
task pc:dev-stop
```

### Individual Tasks

```bash
# See all available commands
task

# Build only
task build

# Test running server
task test-server
```

## Goals

- Run an Open cloud server in cmd/opencloud-server/main.go
- Run a Guigui client in cmd/opencloud-gui/main.go 
- Run a CLI in cmd/opencloud-cli/main.go

- Index documents in data/docs and data/code


## Binaries

This project produces three distinct binaries, each serving a specific purpose:

*   **`opencloud-gui`**: Runs the graphical user interface.
*   **`opencloud-server`**: Runs the OpenCloud collaboration server.
*   **`opencloud-cli`**: Provides command-line utilities for indexing and searching.

## Tools

This project uses consistent naming: `Taskfile.yml` and `processFile.yaml` for clear parallel structure.

Its important for there to be no duplicity in what the Taskfile does and what the processFile does.

The way to approach this is that the Taskfile builds everything, and then the processFile is started to run the binaries, and then the task file is used to call the running opencloud-server to do tests or other things.


### Taskfile

This project uses [Task](https://taskfile.dev/) for task automation.

The main `Taskfile.yml` is in the root of the workspace, with specialized tasks organized into separate included taskfiles.

**Structure:**
- `Taskfile.yml` - Main tasks (build, test, dependencies, etc.)
- `Taskfile.process.yml` - Process orchestration tasks (dev workflows, process management)
- `Taskfile.opencloud.yml` - External OpenCloud binary tasks (download, run, version)

**Process Task Patterns:**
- `pc` - Default command (shows status and available commands)
- `pc:dev*` - Development workflows (build + run + test combinations)
- `pc:run:<CMD>` - Direct process-compose CLI commands (up, down, logs, ps, etc.)
- `pc:dep` - Dependency management (install process-compose)

**Usage:**
- Main tasks: `task build`, `task test`, `task server`
- Process workflows: `task pc:dev`, `task pc:dev-minimal`, `task pc:dev-gui`
- Process CLI commands: `task pc:run:up`, `task pc:run:down`, `task pc:run:logs`, `task pc:run:ps`
- OpenCloud tasks: `task oc:server`, `task oc:version`, `task oc:check`

**Default Commands:**
- `task pc` - Show process-compose status and available commands
- `task oc` - Show OpenCloud binary status and available commands

The taskfile works on all platforms, such as Linux, macOS, and Windows.

The Taskfile is simple and clean with no special characters.

Remember that ":" inside echo statements breaks the taskfile.

### Processfile

This project uses [process-compose](https://github.com/F1bonacc1/process-compose) for **running binaries** with proper orchestration.

https://f1bonacc1.github.io/process-compose/

https://f1bonacc1.github.io/process-compose/cli/process-compose/ shows all the way to control it.



**Responsibilities:**
- Start/stop binaries with dependencies
- Health checks and readiness probes
- Process lifecycle management
- Cross-platform process orchestration

The `Processfile.yaml` contains the process definitions and works cross-platform.

### Usage

Use the Task runner to see available commands and build the project. The project is set up to be developed within a mono repo structure.

## Search

`opencloud`'s search functionality is a key part of this integration.

### Current Implementation

- **Indexing:** By default, `opencloud` uses [Apache Tika](https://docs.opencloud.eu/docs/dev/server/Services/search/Search-info/) for indexing various file types.

- **Search Service:** The search service is built using [bleve](https://github.com/opencloud-eu/opencloud/tree/main/services/search).

- **Query Language:** It uses [KQL (Keyword Query Language)](https://github.com/opencloud-eu/opencloud/tree/main/pkg/kql).

### Integration Goals

The main goal is to enhance the search and indexing capabilities.

0. **using the existing opencloud pks to run a standard open cloud server:** The opencloud-server cmd just uses the existing opencloud pkgs ( https://github.com/opencloud-eu/opencloud/tree/main/opencloud ) to start a server, with our own pkgs to add our customizations.

1.  **Adding a native Golang indexing solution:** This will provide an alternative to the Java-based Apache Tika, simplifying the stack.

2.  **Adding Markdown support:** The new indexer should support `.md` files.

3.  **Adding Decksh support:** Once Markdown support is working, we will add support for `Decksh`.
