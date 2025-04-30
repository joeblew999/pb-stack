# source-remote template


We need Operators to easily build off code remotely.

This is a Compiler as a Service (CAS).

https://github.com/joeblew999/pb-stack-example

## BUGS

BUG: Var overrides are NOT working at all. Only ENV If the .env is in the same folder.


```sh
# first copy the .env-template to .env as task needs this.

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





