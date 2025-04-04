# todo

Testing with a task file here.

Its needs a .taskrc.yml in order to work.

```sh
task --list-all
task: Available tasks for this project:
* bin:                                         
* default:                                     todo print
* base:default:                                base print      (aliases: base)
* base:dep:                                    base dep, installs shell level components.
* base:run:                                    
* git:clone:                                   git clone
* git:config:                                  git config
* git:default:                                 git print      (aliases: git)
* git:pull:                                    git pull
* git:push:                                    git push ( eg: GIT_PUSH_COMMIT_MESSAGE=?? task git:push )
* git:run-dep:                                 git run any deps needed ...
* git:sign-del:                                git sign delete, to delete the signing up in your git config.
* git:sign-get:                                git sign get, to see what settings you have in your git config.
* git:sign-set:                                git sign set, to set the signing up in your git config.
* git:ssh-create:                              git ssh creation. Only once of it not there.
* git:ssh-del:                                 deletes ssh keys and their config. Tricky to be idempotent
* git:ssh-set:                                 git ssh setup. Part of run-dep.
* git:status:                                  git status
* go:bin:                                      go build
* go:default:                                  go print      (aliases: go)
* remote:docker-compose:down:                  Stop and remove containers, networks
* remote:docker-compose:up:                    Create and start containers
* remote:security:sast:checkov:scanner:        Infrastructure as code static analysis
* remote:security:sast:grype:scanner:          A vulnerability scanner for container images, filesystems, and SBOMs
* remote:security:sast:trivy:aws:              [EXPERIMENTAL] Scan AWS account
* remote:security:sast:trivy:config:           Scan config files for misconfigurations
* remote:security:sast:trivy:filesystem:       Scan local filesystem
* remote:security:sast:trivy:image:            Scan a container image
* remote:security:sast:trivy:kubernetes:       [EXPERIMENTAL] Scan kubernetes cluster
* remote:security:sast:trivy:repository:       Scan a repository
* remote:security:sast:trivy:rootfs:           Scan rootfs
* remote:security:sast:trivy:sbom:             Scan SBOM for vulnerabilities and licenses
* remote:security:sast:trivy:vm:               [EXPERIMENTAL] Scan a virtual machine image
* remote:terraform:apply:                      terraform apply -auto-approve
* remote:terraform:destroy:                    terraform destroy
* remote:terraform:doc:                        terraform-docs markdown table
* remote:terraform:fmt:                        terraform fmt
* remote:terraform:init:                       terraform init
* remote:terraform:plan:                       terraform plan
* remote:terraform:terrascan:                  Terrascan static code analyzer
* remote:terraform:tflint:                     tflint
* remote:terraform:validate:                   terraform validate

``` 

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

BASE_PATH: /Users/apple/Library/pnpm:/opt/homebrew/opt/git/libexec/git-core:/opt/homebrew/opt/make/libexec/gnubin:/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:/System/Cryptexes/App/usr/bin:/usr/bin:/bin:/usr/sbin:/sbin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/local/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/bin:/var/run/com.apple.security.cryptexd/codex.system/bootstrap/usr/appleinternal/bin:/Library/Apple/usr/bin:/Users/apple/.cargo/bin:/Users/apple/.orbstack/bin:/Users/apple/workspace:/Users/apple/workspace/go/bin:/opt/homebrew/opt/go/libexec/bin

BASE_SRC_NAME: base_name

BASE_SRC: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all
BASE_BIN: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all/.bin
BASE_DEP: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all/.dep
BASE_PACK: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all/.pack

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
ROOT_TASKFILE: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all
ROOT_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all
TASKFILE: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/base.taskfile.yml
TASKFILE_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack
TASK_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all
USER_WORKING_DIR: /Users/apple/workspace/go/src/github.com/joeblew999/pb-stack/proj/todo-with-all
CHECKSUM:
TIMESTAMP:
TASK_VERSION: v3.42.1
ITEM:
EXIT_CODE:

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

