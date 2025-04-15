# test


only has rc file.
no task file.

## bugs

The variable extrapolation is not immutable:

https://github.com/go-task/task/issues/2180

```sh
task go

- bin
GO_BIN_NAME: go
GO_BIN_WHICH: /opt/homebrew/bin/go
GO_BIN_VERSION: go version go1.24.2 darwin/arm64

- env

- var
GO_VAR_SRC_PATH: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-rc/cmd/02
GO_VAR_BIN_NAME: 02
GO_VAR_BIN_PATH: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-rc/cmd/02/.bin/02



apple@192-168-1-5 02 % task go

- bin
GO_BIN_NAME: go
GO_BIN_WHICH: /opt/homebrew/bin/go
GO_BIN_VERSION: go version go1.24.2 darwin/arm64

- env

- var
GO_VAR_SRC_PATH:
GO_VAR_BIN_NAME:
GO_VAR_BIN_PATH: /
``` 
