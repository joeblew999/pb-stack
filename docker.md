
# DOCKER

Build golang ( and js ) locally, in CI and inside Docker with Task.

See: https://github.com/scaleaq/taskfile-docker

## Usage

```sh

task docker

```

## Files

docker.github.yaml_template is a github action template. It will be part of the templating.

docker.md is this.

docker.mk is the old way. will be deleted later

docker.taskfile is a docker that loads task. 

docker.taskfile.yml is the Taskfile. It will be part of the templating.

Dockerfile is the standard DockerFile

Dockerfile_template is the standard DockerFile template. It will be part of the templating.

Dockerfile.multi is a multi arch docker, that we will get working.


## Templating

It copies the default docker and docker.env and docker.yaml GitHub action into your repo, so you do it once and can make changes.

Its designed to be used as a remote taskfile and is useful when you get into Teams using a best practice setup.

## Multi stage builds

Best practice is to use multistage builds. smaller, more secure.

https://serverfault.com/questions/960648/use-makefile-to-copy-files-in-docker-multi-stage-builds

https://stackoverflow.com/questions/54761769/docker-and-makefile-build


## Fly Remote builder

https://github.com/fly-apps/docker-daemon

GH or FLY. Dont know yet.

## TODO

slowly move the stuff in docker.mk into docker.taskfile.yml

---

Must have multi stage as we need to deploy to amd64, arm64, rasp pi, fly.

https://docs.docker.com/build/building/multi-stage/

Get a decent example...

---

github actions:

get github actions docker, using task, such that local and github actions use the exact same task file.

---

secrets:

when we get SOPS secrets working, we can add the GITHUB_TOKEN into it and task will suck it up and spit it out to the ENV at runtime. This will allow the SOPS secrets to be reused and shared over nats etc.

---

task inside docker:

Run Task inside Docker itself. Obviously best way to have parity.

examples: 

https://github.com/scaleaq/taskfile-docker

https://dev.to/kameshsampath/simplify-your-dockerfile-1j5k





