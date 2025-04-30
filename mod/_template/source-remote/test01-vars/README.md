# source-remote template


## BUGS

BUG: Does not get the values from the .env in the parent folder.

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
GIT_VAR_SRC_REPO_URL: repo-name-default-from-local-env
GIT_VAR_SRC_REPO_NAME: repo-url-default-from-local-env
```



