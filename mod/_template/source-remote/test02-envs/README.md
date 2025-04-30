# source-remote template


## BUGS

BUG: Does not pickup  VAR variable from its own task file, to override the ENV variables OR DEFAULT variables.

```sh
task

- git src
GIT_VAR_SRC_REPO_URL: repo-name-default-from-git.taskfile.yml
GIT_VAR_SRC_REPO_NAME: repo-url-default-from-git.taskfile.ym
```

correct is:

```sh
task

- git src
GIT_VAR_SRC_REPO_URL: repo-name-from-local-task
GIT_VAR_SRC_REPO_NAME: repo-url-from-local-task
```