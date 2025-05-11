# PB-Stack Bootstrapper

This application helps in setting up and tearing down development environments, including package installations and script executions. It can be run in GUI mode or CLI mode.

## Running the Application

### GUI Mode

To run in GUI mode (default):
```bash
go run ./main.go
```
Or using the Taskfile:
```bash
task gui
```

### CLI Mode

To run in CLI mode, you must use the `-cli` flag along with either `-setup` or `-teardown`.

**Available CLI Flags:**

*   `-cli`: (bool)
    Run in command-line mode instead of GUI. This flag is required for CLI operations.
*   `-setup`: (bool)
    Run setup scripts and operations. Must be used with `-cli`.
*   `-teardown`: (bool)
    Run teardown scripts and operations. Must be used with `-cli`.
    *Note: `-setup` and `-teardown` are mutually exclusive.*
*   `-package <name>`: (string)
    Specify a package name (e.g., Winget ID or Homebrew formula) for a targeted setup or teardown operation.
*   `-migrationSet <set_name>`: (string, default: "main")
    Specify which set of migrations to use from the `migrations` folder (e.g., "main", "test"). This affects which `config.json`, `extensions.txt`, and scripts are used.
*   `-logFile <path>`: (string)
    Path to a log file. If empty, logs are written to stderr only.
*   `-debug`: (bool)
    Enable verbose debug logging.

**CLI Examples:**

*   Run general setup using the "main" migration set:
    `go run ./main.go -cli -setup`
*   Run general setup using the "test" migration set and enable debug logging:
    `go run ./main.go -cli -setup -migrationSet test -debug`
*   Setup a specific package:
    `go run ./main.go -cli -setup -package Git.Git`
*   Teardown a specific package and log to a file:
    `go run ./main.go -cli -teardown -package Git.Git -logFile ./app.log`