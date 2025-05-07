# git clone 

https://github.com/go-git/go-git/blob/main/_examples/clone/auth/ssh/private_key/main.go

CheckArgs("<url>", "<directory>", "<private_key_file>")

##

https://github.com/joerdav/xc

```sh
go install github.com/joerdav/xc/cmd/xc@latest

# vscode
code --install-extension https://marketplace.visualstudio.com/items?itemName=xc-vscode.xc-vscode


```


# ensure the right git / ssh config is setup.
```sh
code $HOME/.ssh/config

````


## windows

```sh

```

## unix 

```sh




# https://github.com/joeblew999/golang

rm -rf golang

# fails
#go run . https://github.com/joeblew999/golang.git $PWD/product $HOME/.ssh/joeblew999_github.com

# works without git
go run . git@github.com-joeblew999:joeblew999/golang.git $PWD/product $HOME/.ssh/joeblew999_github.com

# works with git.
git clone git@github.com-joeblew999:joeblew999/golang.git

```
