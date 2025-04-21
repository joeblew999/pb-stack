# Cloudflare tunnel client

# Status: installs, but have not tried much yet.
# Will be powerful for NATS Leaf nodes and Deck Desktop perhaps.

# STATUS: Have not got it all the way yet, but close. Structure is looking good.

# What i need to use it for:

# 1. Connect to my servers over this private network, for Ops being highly secure.
# ex: https://community.hetzner.com/tutorials/connect-over-pvt-net-with-cloudflare-access

# 2. Public connects to my Servers via this, so its very secure.

#3. User can expose their Desktop to Cloud Servers, so that i can sync
# and other users can then interoperate with them Desktop.

# PRICE: its free, so all good.


# https://github.com/cloudflare/cloudflared
# https://github.com/cloudflare/cloudflared/tree/master/cmd/cloudflared

# docs: https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/

# They have Solid Service install code for all OS.
# mac
# https://github.com/cloudflare/cloudflared/blob/master/cmd/cloudflared/macos_service.go
# linux
# https://github.com/cloudflare/cloudflared/blob/master/cmd/cloudflared/linux_service.go
# windows
# https://github.com/cloudflare/cloudflared/blob/master/cmd/cloudflared/windows_service.go

FLARE_TUNNEL_DEP=flare-tunnel
FLARE_TUNNEL_DEP_BIN=cloudflared

FLARE_TUNNEL_DEP_REPO=cloudflared
FLARE_TUNNEL_DEP_REPO_URL=https://github.com/cloudflare/cloudflared
FLARE_TUNNEL_DEP_REPO_DEEP=$(FLARE_TUNNEL_DEP_REPO)/cmd/cloudflared

FLARE_TUNNEL_DEP_MOD=github.com/cloudflare/cloudflared
FLARE_TUNNEL_DEP_MOD_DEEP=$(FLARE_TUNNEL_DEP_MOD)/cmd/cloudflared

# https://github.com/cloudflare/cloudflared/releases/tag/2024.9.1
# https://github.com/cloudflare/cloudflared/releases/tag/2025.1.1
# https://github.com/cloudflare/cloudflared/releases/tag/2025.2.0
# https://github.com/cloudflare/cloudflared/releases/tag/2025.2.1
FLARE_TUNNEL_DEP_VERSION=2025.2.1

FLARE_TUNNEL_DEP_META=flare_tunnel_meta.json

FLARE_TUNNEL_DEP_NATIVE=$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_NATIVE)
FLARE_TUNNEL_DEP_WHICH=$(shell $(BASE_AMP_WHICH_BIN_NAME) $(FLARE_TUNNEL_DEP_NATIVE))

# matches tag version :)
# https://github.com/cloudflare/cloudflared/blob/2025.2.0/cmd/cloudflared/main.go#L34
# we only need "main.version", because we cd into the cmd/cloudflared

# We need special build tags for Step 0 to work...
# https://github.com/cloudflare/cloudflared/blob/master/Makefile#L37
# https://github.com/cloudflare/cloudflared/issues/1414

#export GOFIPS140=latest
FLARE_TUNNEL_GO_BUILD_DATE:=$(shell date -u '+%Y-%m-%d-%H%M')
FLARE_TUNNEL_GO_BUILD_LDFLAGS="-s -w -X main.Version=$(FLARE_TUNNEL_DEP_VERSION) -X main.BuildTime=$(FLARE_TUNNEL_GO_BUILD_DATE)"
FLARE_TUNNEL_GO_BUILD_CMD=CGO_ENABLED=0 $(BASE_DEP_BIN_GO_NAME) build -tags osusergo,netgo -ldflags $(FLARE_TUNNEL_GO_BUILD_LDFLAGS)
# -tags osusergo,netgo,fips
# -tags osusergo,netgo

### meta

FLARE_TUNNEL_DEP_META_PREFIX=flare_tunnel
FLARE_TUNNEL_DEP_META_NAME=$(FLARE_TUNNEL_DEP_META_PREFIX)_meta.json
FLARE_TUNNEL_DEP_META_WHICH=$(BASE_MAKE_IMPORT)/$(FLARE_TUNNEL_DEP_META_NAME)

FLARE_TUNNEL_DEP_TEMPLATE=$(PWD)/$(FLARE_TUNNEL_DEP_META_PREFIX).mk



flare-tunnel-print:
	@echo ""
	@echo "- bin"
	@echo "FLARE_TUNNEL_DEP:              $(FLARE_TUNNEL_DEP)"
	@echo "FLARE_TUNNEL_DEP_VERSION:      $(FLARE_TUNNEL_DEP_VERSION)"
	@echo "FLARE_TUNNEL_DEP_WHICH:        $(FLARE_TUNNEL_DEP_WHICH)"
	@echo "FLARE_TUNNEL_DEP_NATIVE:       $(FLARE_TUNNEL_DEP_NATIVE)"
	@echo ""
	@echo "FLARE_TUNNEL_DEP_VERSION:      $(FLARE_TUNNEL_DEP_VERSION)"
	@echo ""
	@echo "FLARE_TUNNEL_GO_BUILD_DATE:    $(FLARE_TUNNEL_GO_BUILD_DATE)"
	@echo "FLARE_TUNNEL_GO_BUILD_LDFLAGS: $(FLARE_TUNNEL_GO_BUILD_LDFLAGS)"
	@echo "FLARE_TUNNEL_GO_BUILD_CMD:     $(FLARE_TUNNEL_GO_BUILD_CMD)"
	@echo ""
	@echo ""
	@echo "- meta"
	@echo "FLARE_TUNNEL_DEP_META_PREFIX:  $(FLARE_TUNNEL_DEP_META_PREFIX)"
	@echo "FLARE_TUNNEL_DEP_META_NAME:    $(FLARE_TUNNEL_DEP_META_NAME)"
	@echo "FLARE_TUNNEL_DEP_META_WHICH:   $(FLARE_TUNNEL_DEP_META_WHICH)"
	@echo "FLARE_TUNNEL_DEP_TEMPLATE:     $(FLARE_TUNNEL_DEP_TEMPLATE)"
	@echo ""
	

### base

## list base files
flare-tunnel-base-list:
	$(BASE_AMP_TREE_BIN_NAME) -h $(BASE_MAKE_IMPORT)/$(FLARE_TUNNEL_DEP_META_PREFIX)*
## edit base files
flare-tunnel-base-edit:
	$(VSCODE_BIN_NAME) $(BASE_MAKE_IMPORT)/$(FLARE_TUNNEL_DEP_META_PREFIX)*

## edit deployed files
flare-tunnel-dep-edit:
	$(VSCODE_BIN_NAME) $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP_META_PREFIX)*

## list deployed files 
flare-tunnel-dep-list:
	$(BASE_AMP_TREE_BIN_NAME) $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP_META_PREFIX)*


### dep template

## copies a config your cwd, so you can customise it for your project.
flare-tunnel-dep-template-config:

	@echo ""
	@echo "copying template to your cwd ..."

	# create a meta in base cwd
	#rm -rf $(FLARE_TUNNEL_DEP_META_WHICH)
	#echo $(FLARE_TUNNEL_DEP_VERSION) >> $(FLARE_TUNNEL_DEP_META_WHICH)
	
	# copy the configs to base cwd
	cp $(BASE_MAKE_IMPORT)/$(FLARE_TUNNEL_DEP_META_PREFIX)_config* $(BASE_CWD)

flare-tunnel-dep-template: base-dep-init
	@echo ""
	@echo "-version"
	rm -rf $(FLARE_TUNNEL_DEP_META_WHICH)
	echo $(FLARE_TUNNEL_DEP_VERSION) >> $(FLARE_TUNNEL_DEP_META_WHICH)

	# templates to dep.
	cp $(BASE_MAKE_IMPORT)/$(FLARE_TUNNEL_DEP_META_PREFIX)* $(BASE_CWD_DEP)


### dep bin

flare-tunnel-dep-del:
	rm -f $(FLARE_TUNNEL_DEP_WHICH)

flare-tunnel-dep: 
	@echo ""
	@echo "$(FLARE_TUNNEL_DEP) dep check ... "
	@echo ""
	@echo "FLARE_TUNNEL_DEP_WHICH: $(FLARE_TUNNEL_DEP_WHICH)"

ifeq ($(FLARE_TUNNEL_DEP_WHICH), )
	@echo ""
	@echo "$(FLARE_TUNNEL_DEP) dep check: failed"
	$(MAKE) flare-tunnel-dep-single
else
	@echo ""
	@echo "$(FLARE_TUNNEL_DEP) dep check: passed"
endif

flare-tunnel-dep-start:
	rm -rf $(BASE_CWD_DEPTMP)
	mkdir -p $(BASE_CWD_DEPTMP)

	cd $(BASE_CWD_DEPTMP) && $(BASE_DEP_BIN_GIT_NAME) clone $(FLARE_TUNNEL_DEP_REPO_URL) -b $(FLARE_TUNNEL_DEP_VERSION) --single-branch
	cd $(BASE_CWD_DEPTMP) && echo $(FLARE_TUNNEL_DEP_REPO) >> .gitignore
	cd $(BASE_CWD_DEPTMP) && touch go.work
	cd $(BASE_CWD_DEPTMP) && $(BASE_DEP_BIN_GO_NAME) work use $(FLARE_TUNNEL_DEP_REPO)

flare-tunnel-dep-end:
	rm -rf $(BASE_CWD_DEPTMP)

flare-tunnel-dep-single: flare-tunnel-dep-template

	$(MAKE) flare-tunnel-dep-start

ifeq ($(BASE_OS_NAME),darwin)
	@echo "--- darwin ---"
	$(MAKE) flare-tunnel-dep-darwin
endif
ifeq ($(BASE_OS_NAME),linux)
	@echo "--- linux ---"
	$(MAKE) flare-tunnel-dep-linux
endif
ifeq ($(BASE_OS_NAME),windows)
	@echo "--- windows ---"
	$(MAKE) flare-tunnel-dep-windows
endif

	$(MAKE) flare-tunnel-dep-end

flare-tunnel-dep-all: flare-tunnel-dep-template

	$(MAKE) flare-tunnel-dep-start
	
	$(MAKE) flare-tunnel-dep-darwin
	$(MAKE) flare-tunnel-dep-linux
	$(MAKE) flare-tunnel-dep-windows
	
	$(MAKE) flare-tunnel-dep-end

flare-tunnel-dep-darwin:
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_TUNNEL_DEP_REPO_DEEP) && GOOS=darwin GOARCH=amd64 $(FLARE_TUNNEL_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_DARWIN_AMD64)
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_TUNNEL_DEP_REPO_DEEP) && GOOS=darwin GOARCH=arm64 $(FLARE_TUNNEL_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_DARWIN_ARM64)
flare-tunnel-dep-linux:
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_TUNNEL_DEP_REPO_DEEP) && GOOS=linux GOARCH=amd64 $(FLARE_TUNNEL_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_LINUX_AMD64)
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_TUNNEL_DEP_REPO_DEEP) && GOOS=linux GOARCH=arm64 $(FLARE_TUNNEL_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_LINUX_ARM64)
flare-tunnel-dep-windows:
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_TUNNEL_DEP_REPO_DEEP) && GOOS=windows GOARCH=amd64 $(FLARE_TUNNEL_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64)
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_TUNNEL_DEP_REPO_DEEP) && GOOS=windows GOARCH=arm64 $(FLARE_TUNNEL_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64)


BASE_RUN=$(BASE_CWD_DEP)/$(FLARE_TUNNEL_DEP_NATIVE)

flare-tunnel-dep-inspect:
	# Inspect 
	@echo ""
	@echo "BASE_RUN: $(BASE_RUN)"
	@echo ""
	$(MAKE) base-bin-inspect
	$(MAKE) base-amp-redress


	$(MAKE) base-bin-inspect-size

	


# 
# /Users/apple/.cloudflared/cert.pem
FLARE_TUNNEL_VAR_CONFIG_HOME=$(HOME)/.cloudflared


FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_NAME=cert.pem
FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_WHICH=$(FLARE_TUNNEL_VAR_CONFIG_HOME)/$(FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_NAME)


FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_NAME=$(FLARE_TUNNEL_VAR_HOST)_cred.yml
FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_WHICH=$(FLARE_TUNNEL_VAR_CONFIG_HOME)/$(FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_NAME)

FLARE_TUNNEL_VAR_HOST=ubuntu.com
FLARE_TUNNEL_VAR_TUNNEL_NAME=tunnel_$(FLARE_TUNNEL_VAR_HOST)

FLARE_TUNNEL_VAR_TUNNEL_NAME=apple
FLARE_TUNNEL_VAR_TUNNEL_LISTEN=localhost:80


FLARE_TUNNEL_RUN_CMD=$(FLARE_TUNNEL_DEP_NATIVE)

### run

flare-tunnel-run-print:
	@echo ""
	@echo "- run"
	@echo "- run paths"
	@echo "FLARE_TUNNEL_VAR_CONFIG_HOME:    $(FLARE_TUNNEL_VAR_CONFIG_HOME)"
	@echo ""
	@echo "FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_NAME:    $(FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_NAME)"
	@echo "FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_WHICH:   $(FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_WHICH)"
	@echo ""
	@echo ""
	@echo "FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_NAME:    $(FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_NAME)"
	@echo "FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_WHICH:   $(FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_WHICH)"
	@echo ""
	
	@echo ""
	@echo "- run vars"
	@echo "FLARE_TUNNEL_VAR_HOST:           $(FLARE_TUNNEL_VAR_HOST)"
	@echo "FLARE_TUNNEL_VAR_TUNNEL_NAME:    $(FLARE_TUNNEL_VAR_TUNNEL_NAME)"
	@echo "FLARE_TUNNEL_VAR_TUNNEL_LISTEN:  $(FLARE_TUNNEL_VAR_TUNNEL_LISTEN)"
	FLARE_TUNNEL_VAR_TUNNEL_URL
	@echo "- run cmd"
	@echo "FLARE_TUNNEL_RUN_CMD:            $(FLARE_TUNNEL_RUN_CMD)"
	@echo ""

flare-tunnel-run-print-vscode:
	# open tunnel config file in vscode to help see what its doing.
	code $(FLARE_TUNNEL_VAR_CONFIG_HOME)

flare-tunnel-run-print-tree:
	# print, so i can then cat.
	#$(BASE_DEP_BIN_TREE_NAME) --help
	$(BASE_AMP_TREE_BIN_NAME) -h -f $(FLARE_TUNNEL_VAR_CONFIG_HOME)

# STEP 0. Try tunnel for easy test.



flare-tunnel-run-try: flare-tunnel-dep
	# https://try.cloudflare.com
	# assumes port 8080 !
	# brew install cloudflare/cloudflare/cloudflared
	# Nice and easy.
	
	@echo ""
	@echo "Assumes web server on http://localhost:8080"
	@echo ""

	$(FLARE_TUNNEL_DEP_NATIVE) tunnel --url $(FLARE_TUNNEL_VAR_TUNNEL_LISTEN)


# STEP 1. Download 

flare-tunnel-run-h: flare-tunnel-dep
	$(FLARE_TUNNEL_RUN_CMD) -h
flare-tunnel-run-version: flare-tunnel-dep
	$(FLARE_TUNNEL_RUN_CMD) --version
flare-tunnel-run-version-update: flare-tunnel-dep
	# fails with:
	# flare-tunnel_bin_darwin_arm64 update
	# 2025-03-24T11:00:04Z ERR update check failed error="no release found"
	$(FLARE_TUNNEL_RUN_CMD) update

# STEP 2. login and pick a domain.
# https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/#2-authenticate-cloudflared

flare-tunnel-run-tunnel-h: flare-tunnel-dep
	$(FLARE_TUNNEL_RUN_CMD) tunnel -h
flare-tunnel-run-tunnel-login: flare-tunnel-dep
	# Opens Browser to the CF Dash, so i can to pick a Domain / hostname
	# The CLI gives me NO options to control this.
	$(FLARE_TUNNEL_RUN_CMD) tunnel login
	# Certificate is then saved in config home !
	# default: "/Users/apple/.cloudflared/cert.pem"
flare-tunnel-run-tunnel-login-del: flare-tunnel-dep
	# The ONLY way is to delete the Cert file, and then login again.
	rm -f $(FLARE_TUNNEL_VAR_CONFIG_CERT_FILE_WHICH)


# STEP 3. Create a tunnel and give it a name
# https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/#3-create-a-tunnel-and-give-it-a-name


flare-tunnel-run-tunnel-create-with-config: flare-tunnel-dep
	# CONFIG that is used later.
	$(FLARE_TUNNEL_RUN_CMD) tunnel create --config $(FLARE_TUNNEL_VAR_CONFIG_HOME)/config.yml $(FLARE_TUNNEL_VAR_TUNNEL_NAME)

flare-tunnel-run-tunnel-create: flare-tunnel-dep
	# --cred-file  :  Filepath at which to read/write the tunnel credentials [$TUNNEL_CRED_FILE]
	$(FLARE_TUNNEL_RUN_CMD) tunnel create --cred-file $(FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_WHICH) $(FLARE_TUNNEL_VAR_TUNNEL_NAME)
	# Tunnel credentials written to ...
	
flare-tunnel-run-tunnel-delete: flare-tunnel-dep
	$(FLARE_TUNNEL_RUN_CMD) tunnel delete --cred-file $(FLARE_TUNNEL_VAR_CONFIG_CRED_FILE_WHICH) $(FLARE_TUNNEL_VAR_TUNNEL_NAME)
flare-tunnel-run-tunnel-list: flare-tunnel-dep
	$(FLARE_TUNNEL_RUN_CMD) tunnel list
flare-tunnel-run-tunnel-info: flare-tunnel-dep
	# displays details about the active connectors for a given tunnel (identified by name or uuid).
	#$(FLARE_TUNNEL_RUN_CMD) tunnel info -h
	$(FLARE_TUNNEL_RUN_CMD) tunnel info $(FLARE_TUNNEL_VAR_TUNNEL_NAME)
flare-tunnel-run-tunnel-cleanup: flare-tunnel-dep
	# just closes any connections a named tunnel is using
	$(FLARE_TUNNEL_RUN_CMD) tunnel cleanup $(FLARE_TUNNEL_VAR_TUNNEL_NAME)


# STEP 4. Create a configuration file
# https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/get-started/create-local-tunnel/#4-create-a-configuration-file



flare-tunnel-run-access: flare-tunnel-dep
	$(FLARE_TUNNEL_RUN_CMD) access
flare-tunnel-run-access-login: flare-tunnel-dep
	# no idea yet.
	$(FLARE_TUNNEL_RUN_CMD) access login --app ??

flare-tunnel-run-service: flare-tunnel-dep
	# Not sure i need this ? Since this binary itself is invoked via my control plane,
	# I dont think i need it. Lets see.
	$(FLARE_TUNNEL_RUN_CMD) service

#export TUNNEL_ORIGIN_CERT=?
flare-tunnel-run-tail: flare-tunnel-dep
	# Stream logs from a remote cloudflared
	# --origincert value    Path to the certificate generated for your origin when you run cloudflared login. [$TUNNEL_ORIGIN_CERT]
	$(FLARE_TUNNEL_RUN_CMD) tail -h
	$(FLARE_TUNNEL_RUN_CMD) tail --origincert=$(BASE_CWD)/.dep/flare_tunnel_certificate.pem




### start

flare-tunnel-start-proxy-dns: flare-tunnel-dep
	# DNS upstream url=https://1.1.1.1/dns-query
	# DNS upstream url=https://1.0.0.1/dns-query
	# DNS over HTTPS proxy server address=dns://localhost:53
	# Metrics server on 127.0.0.1:57590/metrics
	sudo $(FLARE_TUNNEL_RUN_CMD) proxy-dns







	