# binrun

upstream: https://github.com/PerplexedSphex/binrun

fork: https://github.com/joeblew999/binrun

- NATS running in process

- DS using to update GUI.


Good for SVG Updating. 

For example, https://data-star.dev/examples/progress_bar, uses SVG to update GUI.

APi Generation: https://ogen.dev/docs/spec/extensions/#streaming-json-encoding makes it SSE for DS ?

## tractor

https://github.com/tractordev/toolkit-go/blob/main/go.mod

```sh
go run tractor.dev/toolkit-go/desktop/_example@latest

go install tractor.dev/toolkit-go/desktop/_example@latest

```

## task


The Task team  have a TUI under development.

https://github.com/go-task/task/issues/2077#issuecomment-2938743684

https://github.com/go-task/task/tree/tui

https://github.com/go-task/task/tree/tui/internal/tui


Will impact other attempts:

- https://github.com/titpetric/task-ui/issues/15
- https://github.com/aleksandersh/task-tui/issues/3
- https://github.com/HeyMegabyte/task-dash

```sh

go install github.com/go-task/task/v3/cmd/task@tui

task --tui

``` 

## examples

Task file runs it all.