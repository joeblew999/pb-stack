# carapace-bin

A cross-platform command-line completion tool that generates shell completions for numerous CLI applications. Works across Bash, Fish, Zsh, and other shells.

Part of a Model Context Protocol (MCP) implementation for integrating smart command completion into AI-assisted workflows.

## Tasks

Run tasks using: `task carapace:<taskname>`

- `carapace:dep` - Install dependencies and build the tool
- `carapace:clean` - Clean build artifacts
- `carapace:test` - Run tests

## Usage

```bash
carapace [command]
```

Provides smart command-line completions for supported commands.

## Links

- Documentation: https://carapace.sh

- Spec repo: https://github.com/carapace-sh/carapace-spec

- Code repos:
    - https://github.com/carapace-sh

    - https://github.com/carapace-sh/carapace
        - see

    - https://github.com/carapace-sh/carapace-bin
        - entry points:
            - https://github.com/carapace-sh/carapace-bin/blob/master/cmd/carapace-fmt/main.go
            - https://github.com/carapace-sh/carapace-bin/blob/master/cmd/carapace-generate/gen.go
            - https://github.com/carapace-sh/carapace-bin/blob/master/cmd/carapace-lint/main.go
            - https://github.com/carapace-sh/carapace-bin/blob/master/cmd/carapace-parse/main.go
            - https://github.com/carapace-sh/carapace-bin/blob/master/cmd/carapace-shim/main.go
            - https://github.com/carapace-sh/carapace-bin/blob/master/cmd/carapace/main.go

        - used by: https://github.com/carapace-sh/carapace-bin/network/dependents
            - https://github.com/carapace-sh/freckles

    - https://github.com/carapace-sh/freckles
        - https://github.com/carapace-sh/freckles/blob/master/cmd/freckles/main.go
        - used by: https://github.com/carapace-sh/freckles/network/dependents

    
    - https://github.com/carapace-sh/carapace-bridge
        - https://github.com/carapace-sh/carapace-bridge/blob/master/cmd/carapace-bridge/main.go
        - used by: https://github.com/carapace-sh/carapace-bridge/network/dependents


    - https://github.com/carapace-sh/carapace-selfupdate
        - https://github.com/carapace-sh/carapace-selfupdate/blob/master/cmd/carapace-selfupdate/main.go
        - used by: https://github.com/carapace-sh/carapace-selfupdate/network/dependents

    - https://github.com/carapace-sh/carapace-shlex
        - https://github.com/carapace-sh/carapace-shlex/blob/master/cmd/carapace-shlex/main.go
        - used by: https://github.com/carapace-sh/carapace-shlex/network/dependents
            - seems to be used by all of carapace ?




