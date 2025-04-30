# source-remote template


We need Operators to easily build off code remotely.

This is a Compiler as a Service (CAS).

https://github.com/joeblew999/pb-stack-example

## setup

```sh
cp ./.env-template ./.env

```

## BUGS

BUG: VAR overrides are NOT working at all. See test01-vars folder.

BUG: BUG: Does not pick up ENV variable from the parent folder. See test02-envs folder.


```sh

task git

- git src
GIT_VAR_SRC_REPO_URL: repo-name-default-from-git.taskfile.yml
GIT_VAR_SRC_REPO_NAME: repo-url-default-from-git.taskfile.ym
```

Correct should be:

```sh

task git

- git src
GIT_VAR_SRC_REPO_URL: repo-name-from-local-env
GIT_VAR_SRC_REPO_NAME: repo-url-from-local-env
```





