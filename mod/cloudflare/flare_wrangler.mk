### CLOUDFLARE WRANGLER

# Uses Bun to run C3 ( Create-Cloudflare ), so must be in the make include.


FLARE_WRANGLER_DEP=wrangler2
FLARE_WRANGLER_DEP_BIN=wrangler2

FLARE_WRANGLER_DEP_REPO=cloudflare-wrangler-go
FLARE_WRANGLER_DEP_REPO_URL=https://github.com/cloudflare/cloudflare-go
FLARE_WRANGLER_DEP_REPO_DEEP=$(FLARE_WRANGLER_DEP_REPO)/cmd/flarectl

FLARE_WRANGLER_DEP_MOD=github.com/cloudflare/cloudflare-go
FLARE_WRANGLER_DEP_MOD_DEEP=$(FLARE_WRANGLER_DEP_MOD)/cmd/flarectl

# https://github.com/cloudflare/cloudflare-wrangler-go/releases/tag/v0.104.0
# https://github.com/cloudflare/workers-sdk
FLARE_WRANGLER_DEP_VERSION=v0.104.0

FLARE_WRANGLER_DEP_META=flare_wrangler_meta.json

FLARE_WRANGLER_DEP_NATIVE=$(FLARE_WRANGLER_DEP)
FLARE_WRANGLER_DEP_WHICH=$(shell command -v $(FLARE_WRANGLER_DEP_NATIVE))

# CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo@latest
FLARE_WRANGLER_GO_INSTALL_CMD=CGO_ENABLED=0 $(BASE_DEP_BIN_GO_NAME) install
FLARE_WRANGLER_GO_BUILD_CMD=CGO_ENABLED=0 $(BASE_DEP_BIN_GO_NAME) build

# run 
# hardcoded to my account for now to make things easier.
# removed ALL Cloudflae stuff form $HOME/.zshrc

# flare-wrangler-run-auth-whoami tells you BTW
export CLOUDFLARE_WRANGLER_ACCOUNT_ID=xxx


export CLOUDFLARE_WRANGLER_API_TOKEN=???
export CLOUDFLARE_WRANGLER_EMAIL=gedw99@gmail.com
export WRANGLER_SEND_METRICS=false
export CLOUDFLARE_WRANGLER_API_BASE_URL=https://api.cloudflare.com/client/v4
export WRANGLER_LOG=debug

# THIS IS ONLY for golang code: https://github.com/cloudflare/cloudflare-wrangler-go?tab=readme-ov-file#getting-started
export CLOUDFLARE_WRANGLER_API_KEY=xxx
export CLOUDFLARE_WRANGLER_API_EMAIL=xxx@gmail.com



## THIS IS Needed for each web site, i think ...
CF_DOMAIN=?
CF_SUBDOMAIN=?
CF_ZONE_ID=?
CF_ACCOUNT_ID=?


FLARE_WRANGLER_RUN_ZONE=amplify-cms.com

FLARE_WRANGLER_RUN_PATH=$(PWD)
FLARE_WRANGLER_BIN_PATH=$(FLARE_WRANGLER_RUN_PATH)/.bin

FLARE_WRANGLER_RUN_CMD=cd $(FLARE_WRANGLER_RUN_PATH) && $(FLARE_WRANGLER_DEP_NATIVE)

flare-wrangler-print:
	@echo ""
	@echo "- bin"
	@echo "FLARE_WRANGLER_DEP:             $(FLARE_WRANGLER_DEP)"
	@echo "FLARE_WRANGLER_DEP_VERSION:     $(FLARE_WRANGLER_DEP_VERSION)"
	@echo "FLARE_WRANGLER_DEP_WHICH:       $(FLARE_WRANGLER_DEP_WHICH)"
	@echo "FLARE_WRANGLER_DEP_NATIVE:      $(FLARE_WRANGLER_DEP_NATIVE)"
	
	@echo ""
	
	@echo ""
	@echo "- env"
	@echo ""
	@echo "- more env"
	@echo "CLOUDFLARE_WRANGLER_ACCOUNT_ID:             $(CLOUDFLARE_WRANGLER_ACCOUNT_ID)"
	@echo ""
	@echo "CF_DOMAIN:             $(CF_DOMAIN)"
	@echo "CF_DOMAIN:             $(CF_SUBDOMAIN)"
	@echo "CF_API_TOKEN:          $(CF_API_TOKEN)"
	@echo "CF_ZONE_ID:            $(CF_ZONE_ID)"
	@echo "CF_ACCOUNT_ID:         $(CF_ACCOUNT_ID)"
	@echo ""
	@echo ""
	@echo "- run"
	@echo "FLARE_WRANGLER_RUN_ZONE:        $(FLARE_WRANGLER_RUN_ZONE)"

	@echo "FLARE_WRANGLER_RUN_PATH:        $(FLARE_WRANGLER_RUN_PATH)"
	@echo "FLARE_WRANGLER_BIN_PATH:        $(FLARE_WRANGLER_BIN_PATH)"
	@echo "FLARE_WRANGLER_RUN_CMD:         $(FLARE_WRANGLER_RUN_CMD)"


### dep

flare-wrangler-dep-template: base-dep-init
	@echo ""

	@echo ""
	@echo "-version"
	rm -rf $(BASE_MAKE_IMPORT)/$(FLARE_WRANGLER_DEP_META)
	echo $(FLARE_WRANGLER_DEP_VERSION) >> $(BASE_MAKE_IMPORT)/$(FLARE_WRANGLER_DEP_META)

	@echo ""
	cp -r $(BASE_MAKE_IMPORT)/$(FLARE_WRANGLER_DEP_META) $(BASE_CWD_DEP)
	
	#cp -r $(BASE_MAKE_IMPORT)/flare_wrangler.go $(BASE_CWD_DEP)

	cp -r $(BASE_MAKE_IMPORT)/flare_wrangler.md $(BASE_CWD_DEP)
	cp -r $(BASE_MAKE_IMPORT)/flare_wrangler.mk $(BASE_CWD_DEP)

flare-wrangler-dep-del:
	rm -f $(FLARE_WRANGLER_DEP_WHICH)

flare-wrangler-dep: 
	@echo ""
	@echo "$(FLARE_WRANGLER_DEP) dep check ... "
	@echo ""
	@echo "FLARE_WRANGLER_DEP_WHICH: $(FLARE_WRANGLER_DEP_WHICH)"

ifeq ($(FLARE_WRANGLER_DEP_WHICH), )
	@echo ""
	@echo "$(FLARE_WRANGLER_DEP) dep check: failed"
	$(MAKE) flare-wrangler-dep-single
else
	@echo ""
	@echo "$(FLARE_WRANGLER_DEP) dep check: passed"
endif

flare-wrangler-dep-single: flare-wrangler-dep-template
	# Am just using brew and whatever for now.
	@echo ""
	@echo "installed with your OS package manager ..."
	
ifeq ($(BASE_OS_NAME),darwin)
	brew install cloudflare-wrangler2
endif
ifeq ($(BASE_OS_NAME),linux)
	brew install cloudflare-wrangler2
endif
ifeq ($(BASE_OS_NAME),windows)
	$(MAKE) base-bin-windows
endif


flare-wrangler-dep-all: flare-wrangler-dep-template
	@echo ""
	rm -rf $(FLARE_WRANGLER_DEP_REPO)
	$(BASE_DEP_BIN_GIT_NAME) clone $(FLARE_WRANGLER_DEP_REPO_URL) -b $(FLARE_WRANGLER_DEP_VERSION)
	@echo $(FLARE_WRANGLER_DEP_REPO) >> .gitignore
	touch go.work
	go work use $(FLARE_WRANGLER_DEP_REPO)

	@echo ""
	cd $(FLARE_WRANGLER_DEP_REPO_DEEP) && GOOS=darwin GOARCH=amd64 $(FLARE_WRANGLER_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_WRANGLER_DEP)_$(BASE_BIN_SUFFIX_DARWIN_AMD64)
	cd $(FLARE_WRANGLER_DEP_REPO_DEEP) && GOOS=darwin GOARCH=arm64 $(FLARE_WRANGLER_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_WRANGLER_DEP)_$(BASE_BIN_SUFFIX_DARWIN_ARM64)
	
	cd $(FLARE_WRANGLER_DEP_REPO_DEEP) && GOOS=linux GOARCH=amd64 $(FLARE_WRANGLER_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_WRANGLER_DEP)_$(BASE_BIN_SUFFIX_LINUX_AMD64)
	cd $(FLARE_WRANGLER_DEP_REPO_DEEP) && GOOS=linux GOARCH=arm64 $(FLARE_WRANGLER_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_WRANGLER_DEP)_$(BASE_BIN_SUFFIX_LINUX_ARM64)
	
	cd $(FLARE_WRANGLER_DEP_REPO_DEEP) && GOOS=windows GOARCH=amd64 $(FLARE_WRANGLER_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_WRANGLER_DEP)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64)
	cd $(FLARE_WRANGLER_DEP_REPO_DEEP) && GOOS=windows GOARCH=arm64 $(FLARE_WRANGLER_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_WRANGLER_DEP)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64)

	rm -rf $(FLARE_WRANGLER_DEP_REPO)
	rm -f go.work
	rm -f go.work.sum

	#touch go.work
	#go work use $(OS_MOD)

flare-wrangler-run-h: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) -h

flare-wrangler-run-create: bun-dep flare-wrangler-dep
	$(BUN_DEP_NATIVE) create cloudflare@latest
flare-wrangler-run-wrangler: bun-dep flare-wrangler-dep
	$(BUN_DEP_NATIVE) install wrangler@latest
	

flare-wrangler-run-auth-login: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) login
flare-wrangler-run-auth-logout: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) logout
flare-wrangler-run-auth-whoami: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) whoami


flare-wrangler-run-build-del:
	rm -rf $(FLARE_WRANGLER_BIN_PATH)
flare-wrangler-run-build: flare-wrangler-dep
	cd $(FLARE_WRANGLER_RUN_PATH) && go run github.com/syumai/workers/cmd/workers-assets-gen@v0.23.1 -mode=go -o .bin
	cd $(FLARE_WRANGLER_RUN_PATH) && GOOS=js GOARCH=wasm go build -o $(FLARE_WRANGLER_BIN_PATH)/app.wasm .
flare-wrangler-run-dev: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) dev
flare-wrangler-run-deploy: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) deploy

FLARE_WRANGLER_RUN_PROJECT_NAME=aaa

flare-wrangler-run-pages-h: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) pages -h
flare-wrangler-run-pages-project-h: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) pages project -h
flare-wrangler-run-pages-project-list: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) pages project list
flare-wrangler-run-pages-project-create: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) pages project create $(FLARE_WRANGLER_RUN_PROJECT_NAME)
flare-wrangler-run-pages-project-delete: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) pages project delete $(FLARE_WRANGLER_RUN_PROJECT_NAME)
	




flare-wrangler-run-workflows-h: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) workflows -h
flare-wrangler-run-workflows-list: flare-wrangler-dep
	$(FLARE_WRANGLER_RUN_CMD) workflows list






	