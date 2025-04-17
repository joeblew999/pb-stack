# Github

https://github.com/joeblew999/pb-stack


Task files for cross platform development and runtimes.

https://taskfile.dev

https://github.com/go-task/task

## Status


status ! WIP !

0. Determine SHELL, and manipulate it.

Task can not determine the shell you have configured, or the path to the shell config, so we need a Golang tool for it, that we install as part of bootstrap, that Task can then call.

1. Manipulate PATH.

Environment PATH manipulation does not exist in Task, where we add .DEP and .BIN at runtime, so we need a Golang tool, that Task calls.

2. File path joins.

File path joins is broken in TASK, so we need a golang tool for that, which Task will call.


## Documentation

See [Doc](../doc/README.md) folder for Project Info.

## 0. Base OS Setup

You need golang, git and ssh. Nothing else.

### For Darwin / Linux

Install Brew, and then:

```sh
which go
which git
which openssh
```

```sh
brew install go
brew install git
brew install openssh
```

Open your Shell config. For Zsh its:
```sh
code ~/.zshrc
```

### For Windows

Windows Package Manager (aka WinGet) comes pre-installed with Windows 11 (21H2 and later). It can also be found in the Microsoft Store or be installed directly.

https://learn.microsoft.com/en-us/windows/package-manager/winget/

Install Winget from Powershell:

```sh
Add-AppxPackage -RegisterByFamilyName -MainPackage Microsoft.DesktopAppInstaller_8wekyb3d8bbwe
```

Then install which, like we have on Darwin / Linux:

```sh 
winget list GnuWin32.Which
winget search GnuWin32.Which
```

```sh
winget install GnuWin32.Which

# Check it with:
#which which
#C:\Users\admin\AppData\Local\Microsoft\WinGet\Links\which.EXE
```

Then install golang, git and ssh:

```sh
winget list Golang.go
winget search Golang.go

winget list Git.git
winget search Git.git

winget list Openssh
winget search Openssh
```

```sh 
winget install GnuWin32.Which
winget install Golang.go
winget install Git.git
winget install Microsoft.OpenSSH.Preview
```

To easily control winget using golang using https://github.com/mbarbita/go-winget 

```sh
go install github.com/mbarbita/go-winget@latest 
```

## 2. Sync the SSH on all your devices

SSH is all that is needed, and Wormhole makes this easy.

To move ssh config & keys around between systems using a CLI for wormhole.

https://github.com/psanford/wormhole-william

```sh
go install github.com/psanford/wormhole-william@latest
``` 

To move files between machines using a GUI for wormhole.

```sh
go install github.com/Jacalz/rymdport/v3@latest
```




To check ssh keys on Github match your own. This is public BTW.

https://github.com/joeblew999.keys


```sh
cat $HOME/.ssh/joeblew999_github.com.pub
``` 


## 1. Task Setup

Assuming you have golang installed, bootstrap task onto to your laptop...

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

## 2. Env Setup

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

## 3. GIT and SSH setup

In the git task file, are all the SSH and GIT functions.







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


## Helpers for base OS setup of SSH keys and files for Dev-time and Run-time

