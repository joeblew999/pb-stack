
# DOCKER

# assumes docker is installed. 
# assumes fly.mk is included.

# TODO:

# to build docker cross platform i need the "containerd image store enabled"
# https://github.com/orbstack/orbstack/issues/887#issuecomment-2159293284
# ~/.orbstack/config/docker.json

# from https://github.com/greenpau/caddy-lambda/blob/main/Makefile
# go install github.com/greenpau/versioned/cmd/versioned@latest


DOCKER_DEP_BIN=docker

DOCKER_DEP_NATIVE=$(DOCKER_DEP_BIN)
ifeq ($(BASE_OS_NAME),windows)
	DOCKER_DEP_NATIVE=$(DOCKER_DEP_BIN).exe
endif

DOCKER_DEP_WHICH:=$(shell $(BASE_DEP_BIN_WHICH_NAME) $(DOCKER_DEP_NATIVE))
DOCKER_DEP_WHICH_VERSION=$(shell $(DOCKER_DEP_NATIVE) -v)

DOCKER_DEP_META=docker_meta.json

### meta
DOCKER_DEP_META_PREFIX=docker
DOCKER_DEP_META_NAME=$(DOCKER_DEP_META_PREFIX)_meta.json
DOCKER_DEP_META_WHICH=$(BASE_MAKE_IMPORT)/$(DOCKER_DEP_META_NAME)

### compose

DOCKER_COMPOSE_DEP_BIN=docker-compose
DOCKER_COMPOSE_DEP_NATIVE=$(DOCKER_COMPOSE_DEP_BIN)
ifeq ($(BASE_OS_NAME),windows)
	DOCKER_COMPOSE_DEP_NATIVE=$(DOCKER_COMPOSE_DEP_BIN).exe
endif

DOCKER_COMPOSE_DEP_WHICH=$(shell $(BASE_DEP_BIN_WHICH_NAME) $(DOCKER_COMPOSE_DEP_NATIVE))
DOCKER_COMPOSE_DEP_WHICH_VERSION=$(shell $(DOCKER_COMPOSE_DEP_NATIVE) -v)


docker-print: base-dep

	@echo ""
	@echo "--- DOCKER: bin ---"
	@echo "DOCKER_DEP_BIN:                    $(DOCKER_DEP_BIN)"
	@echo "DOCKER_DEP_NATIVE:                 $(DOCKER_DEP_NATIVE)"
	@echo "DOCKER_DEP_WHICH:                  $(DOCKER_DEP_WHICH)"
	@echo "DOCKER_DEP_WHICH_VERSION:          $(DOCKER_DEP_WHICH_VERSION)"
	@echo ""

	@echo "--- DOCKER COMPOSE: bin ---"
	@echo "DOCKER_COMPOSE_DEP_BIN:            $(DOCKER_COMPOSE_DEP_BIN)"
	@echo "DOCKER_COMPOSE_DEP_NATIVE:         $(DOCKER_COMPOSE_DEP_NATIVE)"
	@echo "DOCKER_COMPOSE_DEP_WHICH:          $(DOCKER_COMPOSE_DEP_WHICH)"
	@echo "DOCKER_COMPOSE_DEP_WHICH_VERSION:  $(DOCKER_COMPOSE_DEP_WHICH_VERSION)"
	@echo ""
	@echo "- meta"
	@echo "DOCKER_DEP_META_PREFIX:            $(DOCKER_DEP_META_PREFIX)"
	@echo "DOCKER_DEP_META_NAME:              $(DOCKER_DEP_META_NAME)"
	@echo "DOCKER_DEP_META_WHICH:             $(DOCKER_DEP_META_WHICH)"
	@echo ""
	@echo ""


docker-base-list:
	$(BASE_DEP_BIN_TREE_NAME) -h $(BASE_MAKE_IMPORT)/$(DOCKER_DEP_META_PREFIX)*

docker-base-edit:
	$(VSCODE_BIN_NAME) $(BASE_MAKE_IMPORT)/$(DOCKER_DEP_META_PREFIX)*

### dep

docker-dep-list: 
	ls $(BASE_CWD_DEP)/docker*

docker-dep-template: base-dep-init
	
	@echo ""
	@echo "-version"
	rm -rf $(BASE_MAKE_IMPORT)/$(DOCKER_DEP_META)
	echo $(DOCKER_DEP_VERSION) >> $(BASE_MAKE_IMPORT)/$(DOCKER_DEP_META)

	# templates to dep.

	# templates to dep.
	cp -r $(BASE_MAKE_IMPORT)/docker* $(BASE_CWD_DEP)
	
	#cp -r $(BASE_MAKE_IMPORT)/docker_config_dockerfile $(BASE_CWD_DEP)

	#cp -r $(BASE_MAKE_IMPORT)/$(DOCKER_DEP_META) $(BASE_CWD_DEP)
	#cp -r $(BASE_MAKE_IMPORT)/caddy_test.mk $(BASE_CWD_DEP)
	#cp -r $(BASE_MAKE_IMPORT)/docker.go $(BASE_CWD_DEP)
	#cp -r $(BASE_MAKE_IMPORT)/docker.md $(BASE_CWD_DEP)
	#cp -r $(BASE_MAKE_IMPORT)/docker.mk $(BASE_CWD_DEP)


docker-dep-del:
	# BE CAREFUL !!
	#rm -f $(DOCKER_DEP_WHICH)

docker-dep: 
	@echo ""
	@echo " $(DOCKER_DEP) dep check ... "
	@echo ""
	@echo "DOCKER_DEP_WHICH: $(DOCKER_DEP_WHICH)"

ifeq ($(DOCKER_DEP_WHICH), )
	@echo ""
	@echo "$(DOCKER_DEP) dep check: failed"
	$(MAKE) docker-dep-single
else
	@echo ""
	@echo "$(DOCKER_DEP) dep check: passed"
endif

#export GITHUB_ACTIONS=true

docker-dep-single: docker-dep-template
    # We dont need to ...
	@echo ""
	@echo "Now installing docker, based on OS sniffing"
	@echo ""

# Check if inside Github Actions
# os.getenv("GITHUB_ACTIONS")
# os.getenv("TRAVIS")
# os.getenv("CIRCLECI")
# os.getenv("GITLAB_CI")

ifeq ($(GITHUB_ACTIONS)),true)
	@echo "- github found"
	# NO IDEA YET what to do yet :)
endif

ifeq ($(BASE_OS_NAME),darwin)
	@echo "- darwin found"
    # really better to pull the binary.
    #brew install --cask orbstack
endif
ifeq ($(BASE_OS_NAME),linux)
	@echo "- linux found"
endif
ifeq ($(BASE_OS_NAME),windows)
	@echo "- windows found"
	# https://docs.orbstack.dev/compare/utm
	# Just thinking how we can use UTM later to automate it.
	# IF we do then its really a different .mk file
endif
	@echo ""
	@echo "Finished installing docker"
	@echo ""


## Installs binaries into side docker, to make it easy with fly.

docker-inject:
	@echo ""
	@echo "Copying binary into docker ..."
	# pull out of runnign docker compose ...
	# push binary into docker ...
	@echo ""

# https://earthly.dev/blog/docker-and-makefiles/

### run

# Registry.
# PICK ONE - Github is the one to use. Fly is not public.
# Can then boot fly, hetzner etc off gh registry

DOCKER_VAR_REGISTRY=ghcr.io
# fly: https://til.simonwillison.net/fly/fly-docker-registry
# registry.fly.io/your-app-name:unique-tag-for-your-image
# flyctl deploy --app datasette-demo --image registry.fly.io/datasette-demo:datasette-demo-v0
#DOCKER_VAR_REGISTRY=registry.fly.io

DOCKER_VAR_USERNAME=gedw99
# fly uses its own way, with "fly-run-auth-login". You can check with "fly-run-auth-whoami"
DOCKER_VAR_PASSWORD=$(GITHUB_TOKEN)
# https://github.com/gedw99?tab=packages
# HELP: https://github.com/cli/cli/packages
# https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry
DOCKER_VAR_PACKAGES=https://github.com/$(DOCKER_VAR_USERNAME)?tab=packages


DOCKER_VAR_APPNAME=$(BASE_BIN_NAME)
DOCKER_VAR_TAG=$(BASE_GITROOT_TAG)
DOCKER_VAR_LABEL=label-x

# var
DOCKER_VAR_CONTAINER_NAME=$(DOCKER_VAR_APPNAME)
DOCKER_VAR_IMAGE_NAME=$(DOCKER_VAR_APPNAME)
# DOCKER_REG=ghcr.io/ripienaar/choria-compose:latest
DOCKER_VAR_IMAGE_WHICH=$(DOCKER_VAR_REGISTRY)/$(DOCKER_VAR_USERNAME)/$(DOCKER_VAR_IMAGE_NAME)
DOCKER_VAR_IMAGE_WHICH_TAGGED=$(DOCKER_VAR_IMAGE_WHICH):$(DOCKER_VAR_TAG)

DOCKER_RUN_NETWORK_NAME=docker-network

DOCKER_RUN_PATH=$(BASE_CWD_DEP)
DOCKER_RUN_DOCKER_FILE_NAME=Dockerfile
DOCKER_RUN_DOCKER_FILE_WHICH=$(DOCKER_RUN_PATH)/$(DOCKER_RUN_DOCKER_FILE_NAME)

DOCKER_COMPOSE_RUN_PATH=$(BASE_CWD_DEP)
DOCKER_COMPOSE_RUN_FILE_NAME=compose.yml
DOCKER_COMPOSE_RUN_FILE_WHICH=$(DOCKER_COMPOSE_RUN_PATH)/$(DOCKER_COMPOSE_RUN_FILE_NAME)

DOCKER_RUN_DOCKERSWARM_FILE=$(DOCKER_COMPOSE_RUN_PATH)/compose-swarm.yml
DOCKER_RUN_DOCKERSWARM_STACK_NAME=email-server

DOCKER_RUN_CMD=cd $(CADDY_RUN_PATH) && $(DOCKER_DEP_NATIVE)


docker-run-print: base-dep
	@echo ""
	@echo "--- DOCKER: run ---"
	@echo "DOCKER_VAR_REGISTRY:               $(DOCKER_VAR_REGISTRY)"
	@echo "- Option:"
	@echo "  DOCKER_VAR_REGISTRY=ghcr.io"    
	@echo "  DOCKER_VAR_REGISTRY=registry.fly.io"               
	@echo "DOCKER_VAR_USERNAME:               $(DOCKER_VAR_USERNAME)"
	@echo "DOCKER_VAR_PASSWORD:               $(DOCKER_VAR_PASSWORD)"
	@echo "DOCKER_VAR_PACKAGES:               $(DOCKER_VAR_PACKAGES)"
	
	@echo ""
	@echo "DOCKER_VAR_APPNAME:                $(DOCKER_VAR_APPNAME)"
	@echo "DOCKER_VAR_CONTAINER_NAME:         $(DOCKER_VAR_CONTAINER_NAME)"
	@echo "DOCKER_VAR_IMAGE_NAME:             $(DOCKER_VAR_IMAGE_NAME)"
	@echo "DOCKER_VAR_IMAGE_WHICH:            $(DOCKER_VAR_IMAGE_WHICH)"
	@echo "DOCKER_VAR_IMAGE_WHICH_TAGGED:     $(DOCKER_VAR_IMAGE_WHICH_TAGGED)"
	@echo "DOCKER_RUN_NETWORK_NAME:           $(DOCKER_RUN_NETWORK_NAME)"
	
	@echo ""
	@echo "--- DOCKER: run cmd ---"
	@echo "DOCKER_RUN_PATH:                   $(DOCKER_RUN_PATH)"
	@echo "DOCKER_RUN_DOCKER_FILE_NAME:       $(DOCKER_RUN_DOCKER_FILE_NAME)"
	@echo "DOCKER_RUN_DOCKER_FILE_WHICH:      $(DOCKER_RUN_DOCKER_FILE_WHICH)"
	@echo ""
	@echo "--- DOCKER COMPOSE: run cmd ---"
	@echo "DOCKER_COMPOSE_RUN_PATH:           $(DOCKER_COMPOSE_RUN_PATH)"
	@echo "DOCKER_COMPOSE_RUN_FILE_NAME:      $(DOCKER_COMPOSE_RUN_FILE_NAME)"
	@echo "DOCKER_COMPOSE_RUN_FILE_WHICH:     $(DOCKER_COMPOSE_RUN_FILE_WHICH)"
	@echo ""
	@echo "--- DOCKER SWARM: run cmd ---"
	@echo "DOCKER_RUN_DOCKERSWARM_FILE:       $(DOCKER_RUN_DOCKERSWARM_FILE)"
	@echo "DOCKER_RUN_DOCKERSWARM_STACK_NAME: $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)"
	
	@echo ""

docker-run-h:
	$(DOCKER_DEP_NATIVE) -h


### clean

## Docker clean list
docker-clean:
	# clear out all docker images and containers
	@echo ""
	@echo "1: deleting docker containers ... "
	$(DOCKER_DEP_NATIVE) rmi -f $(shell $(DOCKER_DEP_NATIVE) images -qa)

	@echo ""
	@echo "2. deleting docker images ... "
	$(DOCKER_DEP_NATIVE) rm -f $(shell $(DOCKER_DEP_NATIVE) ps -qa)

	@echo ""
	@echo "3. deleting docker volumes ... "
	# If not volume then just errors out
	$(DOCKER_DEP_NATIVE) volume rm $(docker volume ls -qf dangling=true)

## Docker clean system
docker-clean-system:
	# works
	# system level. saves lots of space
	# Deletes build cache objects
	@echo ""
	@echo "1: removing unused docker containers, networks, everything ... "
	docker system prune -f

## Docker clean list
docker-clean-list:
	# works
	# shows disk used for all things.
	@echo ""
	@echo "1. docker containers ... "
	$(DOCKER_DEP_NATIVE) ps
	$(DOCKER_DEP_NATIVE) ps --format json

	@echo ""
	@echo "2. docker images ... "
	$(DOCKER_DEP_NATIVE) images
	$(DOCKER_DEP_NATIVE) images --format json
	
	@echo ""
	@echo "3. docker volumes ... "
	$(DOCKER_DEP_NATIVE) volume ls
	$(DOCKER_DEP_NATIVE) volume ls --format json

	@echo ""
	@echo "4. docker summary ... "
	$(DOCKER_DEP_NATIVE) system df
	$(DOCKER_DEP_NATIVE) system df --format json

docker-watch:
	# Outpust in json.
	# when a build happens it tells me as JSON.
	# Will take a while to use it, as the putput is a bit cryptic.
	$(DOCKER_DEP_NATIVE) system events --format json
	

### login 


docker-run-login-h:
	$(DOCKER_DEP_NATIVE) login --help

## docker-run-login
docker-run-login:
	# cant do docker pulls unless auth done.

	@echo ""
	@echo "Loging in to Container Registry ..."
	@echo ""
	@echo "- registry:      $(DOCKER_VAR_REGISTRY)" 
	@echo ""
	
    # Have to switch based on registry
	@echo ""
	@echo "- switching login mathod based on registry..."
	@echo ""
ifeq ($(DOCKER_VAR_REGISTRY),ghcr.io)
	@echo ""
	@echo "- detected github ..."
	# NOTES: https://gist.github.com/yokawasa/841b6db379aa68b2859846da84a9643c
	echo "$(GITHUB_TOKEN)" | $(DOCKER_DEP_NATIVE) login $(DOCKER_VAR_REGISTRY) -u $(DOCKER_VAR_USERNAME) --password-stdin
else
	@echo ""
	@echo "- detected fly.io ..."
	# flyctl auth docker
	# resultis: Authentication successful. You can now tag and push images to registry.fly.io/{your-app}
	$(FLY_DEP_NATIVE) auth docker
endif

docker-run-login-actions:
	# Github actions need the same token for it to login too.
	# setup github to build and publish the docker package too.
	# 1. Make a github secret at: https://github.com/gedw99/gio-htmx/settings/secrets/actions
	# NAME:    GH_PAT
	# TOKEN:   



### bin 

docker-run-build:
	# works
	@echo ""
	@echo "Building with Docker image..."
    # Add "--progress=plain --no-cache" to see all things Docker is doing.
	# From mac arm64, so its linux arm64
	#cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) build --tag $(DOCKER_VAR_IMAGE_WHICH_TAGGED) .

	# https://docs.orbstack.dev/docker/images#multiplatform
	# Create a parallel multi-platform builder
	$(DOCKER_DEP_NATIVE) buildx create --name mybuilder --use

	# Make "buildx" the default
	$(DOCKER_DEP_NATIVE) buildx install

	# Build for multiple platforms
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) build --platform linux/amd64,linux/arm64 --tag $(DOCKER_VAR_IMAGE_WHICH_TAGGED) .

docker-run-build-del:
	# works
	@echo ""
	@echo "Deleting with Docker image..."
	$(DOCKER_DEP_NATIVE) rmi -f $(DOCKER_VAR_IMAGE_WHICH_TAGGED)
	
docker-run-build-tag:
	# will use this soon.
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) buildx build \
		--platform linux/amd64,linux/arm64 \
		--tag basecamp/kamal-proxy:${TAG} \
		--tag basecamp/kamal-proxy:latest \
		--label "org.opencontainers.image.title=kamal-proxy" \
		--push .


## run

docker-run-run: 
	# works
	@echo ""
	@echo "Running with Docker container ..."
	# You can name a Container whatever you want. We want to track the container to the Image, and so use a defaults structure
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) run --name=$(DOCKER_VAR_CONTAINER_NAME) $(DOCKER_VAR_IMAGE_WHICH_TAGGED)

docker-run-run-del:
	# works
	@echo ""
	@echo "Deleting with Docker container ..."
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) container stop $(DOCKER_VAR_CONTAINER_NAME)
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) container rm $(DOCKER_VAR_CONTAINER_NAME)

docker-run-run-detached:
	# not useful for me..
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) run --detach --name=$(DOCKER_VAR_CONTAINER_NAME) $(DOCKER_VAR_IMAGE_WHICH_TAGGED)


docker-run-exec:
	# not working... 
	# https://docs-stage.docker.com/engine/reference/commandline/container_exec/
	@echo ""
	@echo "Running command inside Docker container ..."
	$(DOCKER_DEP_NATIVE) container exec -it $(DOCKER_VAR_CONTAINER_NAME) $(DOCKER_VAR_IMAGE_WHICH_TAGGED)
	@echo ""

docker-run-tag:
	@echo ""
	@echo "Tagging with Docker ..."
	# Create a tag TARGET_IMAGE that refers to SOURCE_IMAGE
	#$(DOCKER_DEP_NATIVE) tag $(DOCKER_VAR_REGISTRY)/$(DOCKER_VAR_CONTAINER_NAME) $(DOCKER_VAR_IMAGE_WHICH_TAGGED)
	$(DOCKER_DEP_NATIVE) tag $(DOCKER_VAR_CONTAINER_NAME) $(DOCKER_VAR_IMAGE_NAME)

### push 

## docker-run-push
docker-run-push: docker-run-login
	# https://docs.docker.com/engine/reference/commandline/push/

	# ! NEEDS a REPO. It uses the DOCKER_VAR_IMAGE_WHICH_TAGGED as the rep !!
	# So we cant have more than 1 docker image for a repo ?

	@echo ""
	@echo "Pushing with Docker ..."
	@echo "URL: https://github.com/users/$(DOCKER_VAR_USERNAME)/packages/container/package/$(DOCKER_VAR_APPNAME)"

	$(DOCKER_DEP_NATIVE) push $(DOCKER_VAR_IMAGE_WHICH_TAGGED)

	@echo "You can find it on github at:"
	@echo https://github.com/users/$(DOCKER_VAR_USERNAME)/packages/container/package/$(DOCKER_VAR_APPNAME)
	@echo ""

docker-run-pull:
	@echo ""
	@echo "Pulling with Docker ..."
	$(DOCKER_DEP_NATIVE) pull $(DOCKER_VAR_IMAGE_WHICH_TAGGED)

docker-run-inspect: docker-run-pull
	# Needs to be local to inspect it, so pull always.
	@echo ""
	@echo "Inspecting with Docker ..."
	$(DOCKER_DEP_NATIVE) inspect $(DOCKER_VAR_IMAGE_WHICH_TAGGED)



### network

docker-run-network-list:
	$(DOCKER_DEP_NATIVE) network ls
docker-run-network-create:
	# ex: docker network create caddy
	$(DOCKER_DEP_NATIVE) network create $(DOCKER_RUN_NETWORK_NAME)
docker-run-network-inspect:
	$(DOCKER_DEP_NATIVE) network inspect $(DOCKER_RUN_NETWORK_NAME)
docker-run-network-prune:
	# only deletes the unused ones, which is ok. I think Orb needs certain ones.
	$(DOCKER_DEP_NATIVE) network prune -f



### compose

docker-compose-run-h:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_COMPOSE_DEP_NATIVE) -h
docker-compose-run-build:
	@echo ""
	@echo "Starting docker compose and build ..."
	@echo "DOCKER_COMPOSE_RUN_FILE_WHICH: $(DOCKER_COMPOSE_RUN_FILE_WHICH)"
	@echo ""

	cd $(DOCKER_RUN_PATH) && $(DOCKER_COMPOSE_DEP_NATIVE) -f $(DOCKER_COMPOSE_RUN_FILE_WHICH) up --build
docker-compose-run-run:
	@echo ""
	@echo "Starting docker compose and run ..."
	@echo "DOCKER_COMPOSE_RUN_FILE_WHICH: $(DOCKER_COMPOSE_RUN_FILE_WHICH)"
	@echo ""
	cd $(DOCKER_RUN_PATH) && $(DOCKER_COMPOSE_DEP_NATIVE) -f $(DOCKER_COMPOSE_RUN_FILE_WHICH) up

docker-compose-run-exec:
	# to pass commands into docker compose as its running.
	# make ARG=xxx docker-compose-run-exec
	cd $(DOCKER_RUN_PATH) && $(DOCKER_COMPOSE_DEP_NATIVE) exec $(ARG)



### swarm config
# https://www.pulumi.com/what-is/what-are-docker-configs/
# You will also need to initialize a swarm since Docker Configs are a feature of Docker Swarm.


docker-run-config-h:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) config --help
docker-run-config-list:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) config ls
docker-run-config-create-h:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) config create --help
docker-run-config-create-domains-from-file:
# Create a Docker config for the domains
# docker config create domains_config domains.yml
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) config create domains_config domains.yml
docker-run-config-create-domains-from-args:
# Create a Docker config for the domains
# docker config create domains_config domains.yml
	#cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) config create domains_config domains.yml
docker-run-config-remove:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) config rm domains_config


### swarm / stack

docker-swarm-run-h:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) swarm --help
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) stack --help

docker-swarm-run-init:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) swarm init
docker-swarm-run-leave:
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) swarm leave --force


docker-swarm-run-deploy:
	@echo "Deploying the stack to Docker Swarm ..."
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) stack deploy -c $(DOCKER_RUN_DOCKERSWARM_FILE) $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)
	@echo "Stack deployed successfully."

	# docker stack deploy -c $(SWARM_COMPOSE_FILE) $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)
	

docker-swarm-run-update:
	@echo "Updating stack deployment in Docker Swarm..."
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) stack deploy -c $(DOCKER_RUN_DOCKERSWARM_FILE) $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)
	@echo "Stack updated successfully."

	#docker stack deploy -c $(SWARM_COMPOSE_FILE) $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)
	

docker-swarm-run-remove:
	@echo "Removing stack from Docker Swarm..."
	cd $(DOCKER_RUN_PATH) && $(DOCKER_DEP_NATIVE) stack rm $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)
	@echo "Stack removed successfully."

	# docker stack rm $(DOCKER_RUN_DOCKERSWARM_STACK_NAME)


