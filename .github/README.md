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

3. SSH Config

We need a cross platform way to manipulate the ssh config, and call it.
When i installed openssh to Darwin, git stopped working in vscode.

4. Packaging for golang stuff

https://github.com/imjamesonzeller/tasklight-v3 has a good one for Desktops.

- has service installer and Task file modularity for all OS. NOT ready yet though.

```sh
git clone https://github.com/imjamesonzeller/tasklight-v3

go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest

go install github.com/go-task/task/v3/cmd/task@latest


``` 


## Documentation

See [Doc](../doc/README.md) folder for Project Info.

## 0. Base OS Setup

You need golang, git, ssh and bun.

## Golang that does not effect your already installed golang.

https://github.com/kevincobain2000/gobrew allows installing golang compiler versions on the fly.

For darwin / linux:

```sh
curl -sL https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | bash

# This needs to be run in each new shell. Can add it to your shell.
code $HOME/.bashrc
code $HOME/.zshrc
export PATH="$HOME/.gobrew/current/bin:$HOME/.gobrew/bin:$PATH"

```

For windows:

```sh
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))
``` 

Test it from any github repo that has a go.mod, which is all golang projects:

```sh

```



### For Darwin / Linux

Install Brew, if not already installed

https://docs.brew.sh/Installation

TODO: Is there a golang brew controller, like what we have for Winget ?

```sh
which brew

# Assuming its not installed then...
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" 

brew doctor

brew analytics off
brew analytics

```


Then, install everything by:


```sh

brew install go
brew install git
brew install openssh
brew install bun
``` 

Then configure your shell to use them all:

For Zsh, do:
```sh
code ~/.zshrc
```

Then check that your shell is using all packages by:

```sh

which go
which git
which openssh
which bun

```

To uninstall all packages, do:

```sh

brew uninstall go
brew uninstall git
brew uninstall openssh
brew uninstall bun

brew cleanup -s

which go
which git
which openssh
which bun

```

And remove them from you shell config...


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

Then install golang, git, ssh and bun:

Check package names have not changed by:

```sh
winget list GnuWin32.Which
winget search GnuWin32.Which

winget list Golang.go
winget search Golang.go

winget list Git.git
winget search Git.git

winget list Microsoft.OpenSSH.Preview
winget search Microsoft.OpenSSH.Preview

winget list Oven-sh.Bun
winget search Oven-sh.Bun

```

Then install them by:

```sh 

winget install --id=GnuWin32.Which

winget install --id=Golang.go

winget install --id=Git.git

winget install --id=Microsoft.OpenSSH.Preview

winget install --id=Oven-sh.Bun

which which

which go

which git

which openssh

which bun

```

If you need to remove then, do:

```sh

winget uninstall --id=GnuWin32.Which

winget uninstall --id=Golang.go

winget uninstall --id=Git.git

winget uninstall --id=Microsoft.OpenSSH.Preview

winget uninstall --id=Oven-sh.Bun

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

