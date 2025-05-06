# todo

## web / wasm ?

GUI works, but no server off course. Just add PC..

```sh
task go:dev:wasm:serve
```

---

TUI almost runs was WASM in a browser

```sh
task go:dev:wasm:serve
```

```sh

# github.com/charmbracelet/bubbletea
/Users/apple/workspace/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.5/tea.go:324:8: p.listenForResize undefined (type *Program has no field or method listenForResize)
/Users/apple/workspace/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.5/tea.go:410:8: undefined: suspendSupported
/Users/apple/workspace/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.5/tea.go:571:13: undefined: openInputTTY
/Users/apple/workspace/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.5/tea.go:580:13: undefined: openInputTTY
/Users/apple/workspace/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.5/tty.go:19:2: undefined: suspendProcess
/Users/apple/workspace/go/pkg/mod/github.com/charmbracelet/bubbletea@v1.3.5/tty.go:31:14: p.initInput undefined (type *Program has no field or method initInput)
exit status 1

```

---

web works of course

```sh
task go:dev
```

