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


### WASM

https://github.com/magodo/go-wasmww can run golang as wasm workers, to allow the GUI and the workers to be decoupled.

"
At its basic, it abstracts the exec.Cmd structure, to allow the main thread (the parent process) to create a web worker (the child process), by specifying the WASM URL, together with any arguments or environment variables, if any.
"








