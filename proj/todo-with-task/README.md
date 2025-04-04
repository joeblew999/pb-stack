# todo

Testing with a task file here.

Its needs a .taskrc.yml in order to work.

```sh

task base


- shell
BASE_SHELL: /bin/zsh
BASE_SHELL_OS_NAME: darwin
BASE_SHELL_OS_ARCH: arm64
BASE_GOOS_NAME: darwin
BASE_GOOS_ARCH: arm64

- user
BASE_HOME: /Users/apple
BASE_USER: apple

BASE_PATH: ::/Users/apple/Library/pnpm:/opt/homebrew/opt/git/libexec/git-core:/opt/homebrew/opt/make/libexec/gnubin:/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/Library/Apple/usr/bin:/Users/apple/.cargo/bin:/Users/apple/.orbstack/bin:/Users/apple/workspace:/Users/apple/workspace/go/bin:/opt/homebrew/opt/go/libexec/bin

BASE_SRC_NAME: base_name

BASE_SRC: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo
BASE_BIN: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo/.bin
BASE_DEP: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo/.dep
BASE_PACK: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo/.pack

- task

- task bin

- task env
TASK_TEMP_DIR:
TASK_REMOTE_DIR:
TASK_OFFLINE:
FORCE_COLOR:

- task vars
CLI_ARGS:
CLI_FORCE: false
CLI_SILENT: false
CLI_VERBOSE: false
CLI_OFFLINE: false

The name of the current task:
TASK: base:default
The alias used for the current task, otherwise matches TASK:
ALIAS: base
TASK_EXE: task
ROOT_TASKFILE: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo
ROOT_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo
TASKFILE: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/base.taskfile.yml
TASKFILE_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack
TASK_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo
USER_WORKING_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo
CHECKSUM:
TIMESTAMP:
TASK_VERSION: v3.42.1

``` 

It does NOT see the 4 "signing (env)" values from the .env in the root. It think it should.

```sh

task git

- bin
GIT_BIN_NAME: git
GIT_BIN_WHICH: /opt/homebrew/opt/git/libexec/git-core/git
GIT_BIN_VERSION: git version 2.49.0

- signing (env)
GIT_SIGNING_USER_NAME:
GIT_SIGNING_USER_EMAIL:
GIT_SIGNING_KEY_PRIV:
GIT_SIGNING_KEY:
- signing (var)
GIT_SIGNING_PROGRAM: ssh
GIT_SIGNING_FORMAT: ssh

- var
GIT_GITROOT: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack
GIT_GITROOT_VERSION: f30a4bed41a24810a7eb8c798b8e5f07420eca7e

```

