## Process Compose

https://github.com/F1bonacc1/process-compose

Process Compose is a simple and flexible scheduler and orchestrator to manage non-containerized applications.

For Darwin, Linux and Windows.

**Why?** Because sometimes you just don't want to deal with docker files, volume definitions, networks and docker registries.
Since it's written in Go, Process Compose is a single binary file and has no other dependencies.

1. Create a Config as process-compose.yaml:

```yaml
version: "0.5"

processes:
  hello:
    command: echo 'Hello World'
  pc:
    command: echo 'From Process Compose'
    depends_on:
      hello:
        condition: process_completed
```

2. Start it by running `process-compose` from your terminal.

Check the [Documentation](https://f1bonacc1.github.io/process-compose/launcher/) for more advanced use cases.

#### Features:

- Processes execution (in parallel or/and serially)
- Processes dependencies and startup order
- Process recovery policies
- Manual process [re]start
- Processes arguments `bash` or `zsh` style (or define your own shell)
- Per process and global environment variables
- Per process or global (single file) logs
- Health checks (liveness and readiness)
- Terminal User Interface (TUI) or CLI modes
- Forking (services or daemons) processes
- REST API (OpenAPI a.k.a Swagger)
- Logs caching
- Functions as both server and client
- Configurable shortcuts
- Merge Configuration Files
- Namespaces
- Run Multiple Replicas of a Process
- Run a Foreground Process 
- Themes Support
