# git

We need a golang git server and cli.

We also need ssh cli and server


## https://github.com/go-git/cli  

cli

- go install github.com/go-git/cli/cmd/gogit@latest

server

- go install github.com/go-git/cli/cmd/gogit-http-server@latest

## https://github.com/charmbracelet/soft-serve

works

# "github.com/charmbracelet/ssh" & "github.com/charmbracelet/wish"

https://github.com/charmbracelet/huh/releases/tag/v0.7.0

example: https://github.com/charmbracelet/huh/blob/main/examples/ssh-form/main.go

```sh
go install github.com/charmbracelet/huh/examples/ssh-form@latest
ssh-form
```
