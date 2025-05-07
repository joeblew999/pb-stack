# git clone 

https://github.com/go-git/go-git/blob/main/_examples/clone/auth/ssh/private_key/main.go

CheckArgs("<url>", "<directory>", "<private_key_file>")

```sh

code $HOME/.ssh/config


# https://github.com/joeblew999/golang

rm -rf golang

# fails
go run . https://github.com/joeblew999/golang.git $PWD/product $HOME/.ssh/joeblew999_github.com

# works
go run . git@github.com-joeblew999:joeblew999/golang.git $PWD/product $HOME/.ssh/joeblew999_github.com

# works
git clone git@github.com-joeblew999:joeblew999/golang.git

```
