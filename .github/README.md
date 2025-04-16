# Github

https://github.com/joeblew999/pb-stack

Task files for cross platform development and runtimes.

https://taskfile.dev
https://github.com/go-task/task


status ! WIP !

## Documentation

See [Doc](../doc/README.md) folder for Project Info.

## Task

Assuming you have golang installed, bootstrap task onto to your laptop...

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

## Env

Copy the .env-template to .env to suit your own git and github credentials.

```sh
    cp .env-template .env
```

```sh

# .env

# Each Repo MUST have this.

### task

BASE_TASK_VERSION_ENV=v3.42.1
# https://github.com/go-task/task/tree/recursive-config-search
#BASE_TASK_VERSION_ENV=recursive-config-search



### git

GIT_ORG_NAME=xxx

# below settings are not needed because my files conform,
# to the conventions expected at git.taskfile.yml

#GIT_SIGNING_USER_NAME=xxx
GIT_SIGNING_USER_EMAIL=xxx@gmail.com
#GIT_SIGNING_KEY_PRIV={{.HOME}}/.ssh/xxx_github.com
#GIT_SIGNING_KEY={{.HOME}}/.ssh/xxx_github.com.pub

```

## Task

Then Task takes over and does the next level of Bootstrap.

TASK files are used:

1. Locally for dev.

2. In Github Actions for CI and CD, along with TOFU files.

3. In Production for Upgrades, along with TOFU files.


```sh

task 

task base

```


## Dev-time

Pick a Project or Module

```sh
cd mod && task 
```

or

```sh
cd proj && task 
```

then to build ...

```sh
task
```

then to run ...

```sh
task go:run
```

then to package for distribution.


```sh
task base-bin-pack
``` 

then to push the binary for usage by others.


```sh
task base-bin-push
```


## Run-time

Pick the Project local or remote run time 

```sh
cd proj/example/local
```
or
```sh 
cd proj/example/remote
```
then,  to update your local binary.

```sh
task base-bin-pull
```

then to run your local binary.

```sh
task base-bin-run
```


## helpers for Testing on Desktops.

To install golang and git on Windows easily.

```sh 
winget install Golang.go
winget install Git.git
```

To easily control winget using golang using https://github.com/mbarbita/go-winget 

```sh
go install github.com/mbarbita/go-winget@latest 
```

To move ssh config & keys around between laptops for testing.

https://github.com/psanford/wormhole-william

```sh
go install github.com/psanford/wormhole-william@latest
``` 

To move files between machines for testing with a GUI for wormhole.

```sh
go install github.com/Jacalz/rymdport/v3@latest
```




https://github.com/joeblew999.keys

to check ssh keys match your own. this is public BTW.

```sh
cat $HOME/.ssh/joeblew999_github.com.pub
``` 