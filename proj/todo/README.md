# todo

Testing with a task file here.

Its needs a .taskrc.yml in order to work.

It does not see the values from the .env in the root.

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

