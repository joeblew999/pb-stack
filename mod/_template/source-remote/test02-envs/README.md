# source-remote template


## BUGS

BUG: Does not pick up ENV variable from the parent folder.

```sh
task

- git src
GIT_VAR_SRC_REPO_URL: repo-name-default-from-git.taskfile.yml
GIT_VAR_SRC_REPO_NAME: repo-url-default-from-git.taskfile.ym
```

