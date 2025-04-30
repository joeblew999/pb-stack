# source-remote template


We need Operators to easily build off code remotely.

This is a Compiler as a Service (CAS).

https://github.com/joeblew999/pb-stack-example

## setup

```sh
cp ./.env-template ./.env

```

## BUGS

BUG: VAR variable does not override ENV variable . See this folder.

BUG: VAR variable overrides are not working at all. See test01-vars folder.

BUG: ENV variable from the parent folder not working. See test02-envs folder.


```sh

task git

- git src
GIT_VAR_SRC_REPO_URL: repo-name-from-local-env
GIT_VAR_SRC_REPO_NAME: repo-url-from-local-env

```

Correct should be:

```sh

task git

- git src
GIT_VAR_SRC_REPO_URL: repo-name-from-local-task
GIT_VAR_SRC_REPO_NAME: repo-url-from-local-task
```





