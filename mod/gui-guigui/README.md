# gui-native

## guigui

https://github.com/hajimehoshi/guigui

A Native GUI system for web, desktop and mobile that uses golang.

status: working with task file really well for native and wasm.

## ebitenui

https://github.com/ebitenui/ebitenui

## Usage ideas

### Task gui

Livekit seems to use it at https://github.com/livekit/livekit-cli/blob/main/pkg/bootstrap/bootstrap.go

The task files that is ues are here https://github.com/livekit-examples/index?tab=readme-ov-file



Get all tasks:

```sh
task --json --list-all
```

Run a task:

```sh
task git:status:src --verbose

task: [git:status:src] echo ''

task: [git:status:src] cd /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/mod/gui-guigui/guigui && git status
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean
task: [git:status:src] echo ''

```








