# Github

https://github.com/joeblew999/pb-stack

Task files for cross platform development and runtimes.

https://taskfile.dev

https://github.com/go-task/task


## Terms

**Developer** is a person that works with the compilation of the code.

**Operator** is a person that works with the non compiled code with bun as the interpreter. They are utilising the binaries that the Developer built. An Operator is also a CI and CD system.

**User** is person that uses the final Software Product. They are utilising the code that the Operator build.

**Single sourcing** in data refers to using a single source of information to generate multiple outputs or documents, rather than creating separate content for each need. 

## Philosophy

**Interoperable**:  The system is designed so that the same code can be used by Developers, CI / CD ( in github ), Operators ( on cloud or on premise ), Users.

After many years building systems for large teams, you tend to start using a fair bit of code gen. Golang templates go along way in this regards.

**Heramtic**: This system extends that idea of software productivity to allow different Actors in the Software flow to run using the same system, to reduce the complexity. This is similar to the concept of "Single Sourcing", but applied deeply. 

For example, when you Package a project, the task files, binary dependencies ( for each OS and ARCH) and configuration files are put into the Package. This allows the next Actor in the Software flow to have everything they need to run the system. 

The other aspect is self sovereignty.  This means reducing your dependency on other external systems. NATS Jetstream is a huge help with this, allowing you to for example bootstrap the public SSH Keys onto any device, via a NATS Public system. Of course you do need NATS running somewhere. Another example is data sync between offline apps, which can be easily done using NATS running in a Leaf Node setup.

**GUI** Lastly, there is the GUI aspect. Web, Desktop and Mobile can all be built using simple HTMX constructs with a good quality Webview.  This applies to all Actors in the Software Flow. For example, a Developers laptop can run some tools, an Operator actor can also run some GUI to help with the process, as well as A Web GUI for end users. 

## Base Stack

The basis underneath everything is:

Task provides a cross platform shell.

SQLITE provides the DB. You can scale this out with Master Master copies. YOu can have ephemeral DB's feed off dates, with the DB acting a as Materialised view.

Benthos provides a workflow system.

NATS provides a base connectivity system under everything else.

## Status

Its still very much a WIP, with things changing, but the constructs do work. 

I have built systems like this before and the philosophy does work exceeding well. This is a refactoring of those previous systems that uses Task files and Process Compose. 

## TODO

0. Determine SHELL, and manipulate it.

Task can not determine the shell you have configured, or the path to the shell config, so we need a Golang tool for it, that we install as part of bootstrap, that Task can then call.

1. Manipulate PATH.

Environment PATH manipulation does not exist in Task, where we add .DEP and .BIN at runtime, so we need a Golang tool, that Task calls.

2. File path joins.

File path joins is broken in TASK on Windows, so we need a golang tool for that, which Task will call.

3. SSH Config

We need a cross platform way to manipulate the ssh config, and call it.
When i installed openssh to Darwin, git stopped working in vscode for example.

4. Packaging for golang stuff

https://github.com/imjamesonzeller/tasklight-v3 has a good one for Desktops.

- has service installer and Task file modularity for all OS. NOT ready yet though.

```sh
git clone https://github.com/imjamesonzeller/tasklight-v3

go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest

go install github.com/go-task/task/v3/cmd/task@latest


``` 


## Documentation


You can use the Docs folder for your Project documentation for each Actor type.

This Project uses the top level docs folder.  See [Doc](../doc/README.md) folder for Project Info.

## Boostrapping 

Boostrapping runs at 2 level: 

- OS level setup

- Task level setup


## 0. OS Setup

Developers need golang, git, ssh and bun.

Operators of course need much less. Sometimes none, depending on the Project. 

### For Darwin / Linux

Install Brew, if not already installed

https://docs.brew.sh/Installation


```sh
which brew

# Assuming its not installed then...
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" 

brew doctor

brew analytics off
brew analytics

```


Then, install the os binaries you need by:

```sh

brew install visual-studio-code

brew install go
brew install git
brew install openssh
brew install bun

``` 

Then configure your shell to use them all:


```sh
code ~/.zshrc

# VS Code ( nothing needed )

# GO
# The core golang system
export GOROOT=$(brew --prefix)/opt/go/libexec
# Where you golang projects live, with the bin, pkg, src folders inside. 
export GOPATH=$HOME/workspace/go

# add both to the path.
export PATH=$PATH:$GOROOT/bin
export PATH=$PATH:$GOPATH/bin


## GIT
# The core git system
PATH="$(brew --prefix)/opt/git/libexec/git-core:$PATH"

## BUN ( nothing needed )


```

Then check that your shell can see the binaries:

```sh

which code
/opt/homebrew/bin/code

which go
/opt/homebrew/bin/go

go version
go version go1.24.2 darwin/arm64

which git
opt/homebrew/opt/git/libexec/git-core/git
git version
git version 2.49.0

which openssh
/opt/homebrew/bin/ssh

which bun
/opt/homebrew/bin/bun


```

To uninstall all binaries, do:

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

# Try it with:
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


## 3. Task Setup

Assuming you have golang installed, bootstrap task onto to your laptop...

```sh
go install github.com/go-task/task/v3/cmd/task@latest
```

## 4. Env Setup

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

GIT_VAR_ORG_NAME=xxx

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


## Helpers 

For base OS setup of SSH keys and files for Dev-time and Run-time, golang helpers will eventually be written and less Task based commands will be needed.



