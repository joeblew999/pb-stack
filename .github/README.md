# Github

https://github.com/joeblew999/pb-stack

! STATUS !

Read the Docs for now to get a feel for the intent. Ask questions if something seems off.

The scaffolds for the [TASK](../mod/task/README.md) files and [TOFU](../mod/tofu/README.md) are going in currently. We need this to rapidly develop. 

Next will be the [GitHub Actions for CI and CD](../.github/workflows/README.md) round-tripping to your own Desktop.

Simple examples / playgrounds, so we can work up the code generator against real projects. Its the only way for others to understand how to develop this.

## Documentation

See [Doc](../doc/README.md) folder for Project Info.


## Task

Bootstrap task onto to your laptop...

```sh

make     
TASK_BIN_NAME:            task
TASK_BIN_VERSION:         recursive-config-search
TASK_BIN_WHICH:           
TASK_BIN_WHICH_VERSION:  

make task 
go install github.com/go-task/task/v3/cmd/task@recursive-config-search

make
TASK_BIN_NAME:            task
TASK_BIN_VERSION:         recursive-config-search
TASK_BIN_WHICH:           /Users/apple/workspace/go/bin/task
TASK_BIN_WHICH_VERSION:   3.42.1


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

THen Task takes over and does the next level of Bootstrap.

https://taskfile.dev/reference/

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

