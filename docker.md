
# DOCKER


Build golang locally with task.

Build docker, locally and remote with Task with  TaskFile inside the docker.

See: https://github.com/scaleaq/taskfile-docker

## Usage

docker.taskfile is a docker that loads task. 

---

templating

It copies the default docker and docker.env and docker.yaml GitHub action into your repo, so you do it once and can make changes.

Its designed to be used as a remote taskfile and is useful when you get into Teams using a best practice setup.

---

docker.mk is old stuff.

---

docker.taskfile.yml is a Taskfile dedicated to docker, designed to be reused.

---

DockerFile is a standard Dockerfile that has your stuff in it.

---

Dockerfile.multi is not yet backed.

## Best Practices


Best practice is to use multistage builds. smaller, more secure.
r.

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





