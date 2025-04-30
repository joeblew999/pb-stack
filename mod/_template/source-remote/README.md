# source-remote template

Testing

I am using the Default trick to allow overriding, but i am getting some bugs...

```yaml
  GIT_VAR_SRC_REPO_URL: '{{ .GIT_VAR_SRC_REPO_URL | default "repo-name-default-from-git.taskfile.yml" }}'
```

## Setup

```sh
cp ./.env-template ./.env

```

## BUGS

this folder: 

- BUG: VAR variable does not override ENV variable.

test01-vars folder:

- BUG: Does not get the values from the .env in the parent folder.

test02-envs folder: 

- BUG: Does not pick up ENV variable from the parent folder not working. 


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





