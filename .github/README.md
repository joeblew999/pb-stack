# Github

https://github.com/joeblew999/pb-stack

! STATUS !

Read the Docs for now to get a feel for the intent. Ask questions if something seems off.

The scaffolds for the [TASK](../mod/task/README.md) files and [TOFU](../mod/tofu/README.md) are going in currently. We need this to rapidly develop. 

Next will be the [GitHub Actions for CI and CD](../.github/workflows/README.md) round-tripping to your own Desktop.

Then simple examples / playgrounds, so we can work up the code generator against real projects. Its the only way for others to understand how to develop this.

## Documentation

See [Doc](../doc/README.md) folder for Project Info.


## Make

We need task and golang installed to bootstrap.

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

Copy the -env-template to .env to suit your own git and github credentials.

```sh

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
