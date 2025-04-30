# source-remote template


We need Operators to easily build off code remotely.

This is a Compiler as a Service (CAS).

https://github.com/joeblew999/pb-stack-example

## BUGS

BUG: Does not get the values from the .env in the parent folder 

```sh
task

- git src
GIT_VAR_SRC_REPO_URL: repo-name-default-from-git.taskfile.yml
GIT_VAR_SRC_REPO_NAME: repo-url-default-from-git.taskfile.ym
```

