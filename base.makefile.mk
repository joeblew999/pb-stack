# base.
#
# We use special folders for containment, so that we can fully control whats on the the disk.
# 
# Design for many binaries from 1 single repo.

# Builder using https://github.com/wombatwisdom/wombat-builder
# I do a local git clone, and its time consuming.
# Wombat will do it remotely, put it in NATS Bucket, and i can get it from there FAST !!
# OR https://github.com/superfly/rchab

## TODO
# DONE: Replace command with https://github.com/hairyhenderson/go-which. Only use this inside Modules, and not in base.mk, because base still needs to use "command -v" for now.
# DONE: Replace any ls with TREE
# Replace any mkdir with ?
# Replace rm with ?

# Replace cp with ?
# Replace File Search and Replace, so we can that "base_mod_makefile" and "base_mod_taskfile" will be possible.
  # It will also be useful for Deck, where we can use this to Update Deck lines, without needing line numbers ?
  # It will also be useful for Files ( in the Monaco IDE) 



base-edit:
	$(VSCODE_BIN_NAME) $(BASE_MAKE_IMPORT)/base*
base-list:
	$(BASE_AMP_TREE_BIN_NAME) -h $(BASE_MAKE_IMPORT)/base*

## TODOs
# Add version file that gets put into the .bin named "BIN_NAME_VERSION", with the git bits inside
 # We can use git because ALL file systsm are git file systems.
 # I dont care about ldflags. The file is enough for me. When it goes into NATS the Wildcards will be based on NAME, VERSION, OS, ARCH

### base 
#  BASE_MAKE_IMPORT must be used in the original Makefile ( no exceptions )

BASE_SHELL_OS_NAME := $(shell uname -s | tr A-Z a-z)
BASE_SHELL_OS_ARCH := $(shell uname -m | tr A-Z a-z)

# os
BASE_OS_NAME := $(shell go env GOOS)
BASE_OS_ARCH := $(shell go env GOARCH)

# env we are on. If not empty we are inside github.
BASE_ENV_CI := $(CI)

# makefile
BASE_MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))

# git stuff for now. Doubt i need anything else.
include $(BASE_MAKE_IMPORT)/base.env

# Seems to work. We can reference this from any Makefile, without it including it.
# Will allow me and devs to pick their IDE.
include $(BASE_MAKE_IMPORT)/vscode.mk


# git 
BASE_GITROOT=$(shell git rev-parse --show-toplevel)
BASE_GITROOT_BIN := $(BASE_GITROOT)/.bin
BASE_GITROOT_DATA := $(BASE_GITROOT)/.data

BASE_GITROOT_VERSION=$(shell git rev-parse HEAD)
BASE_GITROOT_TAG=$(shell git describe --abbrev=0 --tags)
# If there is no tag, then call it latest
# Or maybe use the VERSION which is the sha.
# Dont knwo whats best yet.  Kind of depends...
ifeq ($(BASE_GITROOT_TAG), )
	BASE_GITROOT_TAG=latest
endif

# Check if inside Github Actions
export GITHUB_ACTIONS=true
ifeq ($(GITHUB_ACTIONS)),true)
	@echo " ! Detected inside Github CI !"
	BASE_GITROOT=${{ github.workspace }}
	BASE_GITROOT_BIN := ${{ github.workspace }}/.bin
	BASE_GITROOT_DATA := ${{ github.workspace }}/.data
endif

# pwd 
# for a mono repo, we have a top level for stuff
BASE_PWD_NAME := kanka
BASE_PWD:=$(HOME)/$(BASE_PWD_NAME)
BASE_PWD_BIN:=$(BASE_PWD)/.bin
BASE_PWD_DATA:=$(BASE_PWD)/.data
BASE_PWD_META:=$(BASE_PWD)/.meta

# cdw ( aka dev ) to install next to src.
# inside a mono repo, we have project and this is what this is.
BASE_CWD:=$(PWD)
BASE_CWD_BIN_NAME=.bin
BASE_CWD_BIN:=$(PWD)/$(BASE_CWD_BIN_NAME)
BASE_CWD_DATA:=$(PWD)/.data

# deps live in this folder.
# we can toggle the .dep folder to use to be at CWD, MODULE or PWD level,
# based on how we are composing a project together.
# I default to Module level, so that I can share them and so blow up the disk.

# The 4 Options then are:
# A. OS Level in Kanka 
#BASE_CWD_DEP:=$(BASE_PWD)/.dep
#BASE_CWD_DEPTMP:=$(BASE_PWD)/.deptmp
# B. GIT ROOT 
#BASE_CWD_DEP:=$(BASE_GITROOT)/.dep
#BASE_CWD_DEPTMP:=$(BASE_GITROOT)/.deptmp
# C. MODULE level for shared deps.
BASE_CWD_DEP:=$(BASE_MAKE_IMPORT)/.dep
BASE_CWD_DEPTMP:=$(BASE_MAKE_IMPORT)/.deptmp
# D. CWD level for where your src is.
#BASE_CWD_DEP:=$(PWD)/.dep
#BASE_CWD_DEPTMP:=$(PWD)/.deptmp

BASE_CWD_META:=$(PWD)/.meta
BASE_CWD_SRC:=$(PWD)/.src
BASE_CWD_PACK:=$(PWD)/.pack



# .mk files ( and maybe others later), that we reuse, and are copied ot the BASE_GITROOT_BIN
BASE_ARTIFACTS:=?

# so "local binaries" are loaded automagically. 
# EX BASE_CWD:               /Users/apple/workspace/go/src/github.com/gedw99/kanka-cloudflare/modules/fly-apps__docker-daemon
export PATH:=$(PATH):$(BASE_CWD_BIN)

# so "desktop shared binaries" are loaded automatically. These are yours and others. Like aqua, etc
# EX BASE_PWD:              /Users/apple/kanka
export PATH:=$(PATH):$(BASE_PWD_BIN)

# so "repo shared binaries" are load automatically. These are in my repo root.
# EX BASE_GITROOT:           /Users/apple/workspace/go/src/github.com/gedw99/kanka-cloudflare
export PATH:=$(PATH):$(BASE_GITROOT_BIN)

export PATH:=$(PATH):$(BASE_CWD_DEP)

# for tinygo
export PATH:=$(PATH):$(BASE_CWD_DEP)/tinygo/bin

# for gobrew, so we can easily flip go versions.
export PATH:=$(PATH):$(HOME)/.gobrew/current/bin:$(HOME)/.gobrew/bin
#export GOROOT="$(HOME)/.gobrew/current/go"

### cwd
# ops that we need to clean a project.

base-cwd-bin-del:
	@echo ""
	@echo "BASE_CWD_BIN:   $(BASE_CWD_BIN)"
	@echo ""
	rm -rf $(BASE_CWD_BIN)
	@echo ""
base-cwd-data-del:
	@echo ""
	@echo "BASE_CWD_DATA:   $(BASE_CWD_DATA)"
	@echo ""
	rm -rf $(BASE_CWD_DATA)
	@echo ""
base-cwd-dep-del:
	@echo ""
	@echo "BASE_CWD_DEP:   $(BASE_CWD_DEP)"
	@echo ""
	rm -rf $(BASE_CWD_DEP)
	@echo ""
base-cwd-deptmp-del:
	@echo ""
	@echo "BASE_CWD_DEPTMP:   $(BASE_CWD_DEPTMP)"
	@echo ""
	rm -rf $(BASE_CWD_DEPTMP)
	@echo ""
base-cwd-src-del:
	@echo ""
	@echo "BASE_CWD_SRC:   $(BASE_CWD_SRC)"
	@echo ""
	rm -rf $(BASE_CWD_SRC)
	@echo ""
base-cwd-pack-del:
	@echo ""
	@echo "BASE_CWD_PACK:  $(BASE_CWD_PACK)"
	@echo ""
	rm -rf $(BASE_CWD_PACK)
	@echo ""



## deletes our generated binaries

### clean

base-clean-print:
	@echo ""
	@echo ""
	@echo "- This prints the sizes only to get a feel for whats there."
	@echo ""
	@echo "- modules"
	# this is the root "modules" folder. Be VERY caereful...
	#du -shc $(BASE_MAKE_IMPORT)
	@echo ""

	@echo ""
	@echo "- natives"
	#du -shc $(BASE_MAKE_IMPORT)/*$(BASE_BIN_SUFFIX_DARWIN_AMD64)
	@echo ""

	@echo ""
	@echo "- BASE_CWD_BIN"
	du -shc $(BASE_CWD_BIN)	
	@echo ""

	@echo ""
	@echo "- BASE_CWD_DATA"
	du -shc $(BASE_CWD_DATA)
	@echo ""

	@echo ""
	@echo "- BASE_CWD_DEP"
	du -shc $(BASE_CWD_DEP)
	@echo ""

	@echo ""
	@echo "- BASE_CWD_DEPTMP"
	du -shc $(BASE_CWD_DEPTMP)
	@echo ""

	@echo ""
	@echo "- BASE_CWD_SRC"
	du -shc $(BASE_CWD_SRC)
	@echo ""

	@echo ""
	@echo "- BASE_CWD_PACK"
	du -shc $(BASE_CWD_PACK)
	@echo ""
	
base-clean-module:
	# for a single module, from that place.
	# leave all source.

	@echo ""
	@echo "- delting project things."
	@echo "- .bin"
	#rm -rf $(BASE_CWD_BIN)
	@echo "- .dep"
	#rm -rf $(BASE_CWD_DEP)
	@echo "- .deptmp"
	#rm -rf $(BASE_CWD_DEPTMP)
	@echo "- .pack"
	#rm -rf $(BASE_CWD_PACK)

base-clean-module-all:
	# reaches down into all modules ...

	@echo ""
	@echo "- del project things."
	@echo "- .bin"
	ls -al $(BASE_CWD_BIN)
	#rm -rf $(BASE_CWD_BIN)
	@echo ""
	
	@echo ""
	@echo "- .dep"
	#ls -al $(BASE_CWD_DEP)
	#rm -rf $(BASE_CWD_DEP)
	@echo ""

	@echo ""
	@echo "- .deptmp"
	#rm -rf $(BASE_CWD_DEPTMP)
	@echo ""
	
	@echo ""
	@echo "- .pack"
	#rm -rf $(BASE_CWD_PACK)
	@echo ""
	

base-temp:
	exit
	rm -rf $(BASE_MAKE_IMPORT)/*$(BASE_BIN_SUFFIX_DARWIN_AMD64)

	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_DARWIN_ARM64)
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_LINXU_AMD64)
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_LINUX_ARM64)
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_WINDOWS_AMD64)
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_WINDOWS_ARM64)
	# go wasm / wasi
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_WASM_GO)
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_WASI_GO)
	# tinygo wasm / wasi
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_WASM_TINY)
	rm -rf $(OS_DEP_ROOT)/*$(BASE_BIN_SUFFIX_WASI_TINY)

	
	@echo ""
	@echo "ls..."
	ls -al $(BASE_MAKE_IMPORT)

	@echo ""
	@echo "size..."
	du -shc $(BASE_MAKE_IMPORT)


base-clean-os:
	# clean all go
	$(BASE_DEP_BIN_GO_NAME) clean -cache -testcache -modcache -fuzzcache

	# clean all docker
	docker buildx prune --all
	docker system prune --all --volumes

### pre

base-print-pre:
	@echo ""
	
	@echo "--- base : make ---"
	@echo "BASE_MAKE_IMPORT:       $(BASE_MAKE_IMPORT)"
	@echo "BASE_MAKEFILE:          $(BASE_MAKEFILE)"
	@echo "BASE_ARTIFACTS:         $(BASE_ARTIFACTS)"
	@echo ""
	@echo "--- base : shell ---"
	@echo "BASE_SHELL_OS_NAME:     $(BASE_SHELL_OS_NAME)"
	@echo "BASE_SHELL_OS_ARCH:     $(BASE_SHELL_OS_ARCH)"
	@echo ""
	@echo "--- base : os ---"
	@echo "BASE_OS_NAME:           $(BASE_OS_NAME)"
	@echo "BASE_OS_ARCH:           $(BASE_OS_ARCH)"
	@echo ""
	@echo "--- base : paths ---"
	@echo "BASE_ENV_PATH:          $(PATH)"
	@echo ""
	@echo "--- base : env ---"
	@echo "BASE_ENV_CI:            $(BASE_ENV_CI)"
	@echo ""
	@echo ""
	@echo "--- base : Git Root directory ---"
	@echo "BASE_GITROOT:           $(BASE_GITROOT)"
	@echo "BASE_GITROOT_BIN:       $(BASE_GITROOT_BIN)"
	@echo "BASE_GITROOT_DATA:      $(BASE_GITROOT_DATA)"
	@echo ""
	@echo "--- base : Git Root Version---"
	@echo "BASE_GITROOT_VERSION:   $(BASE_GITROOT_VERSION)"
	@echo "BASE_GITROOT_TAG:       $(BASE_GITROOT_TAG)"
	@echo ""
	@echo "--- base : Personal Working folders ---"
	@echo "Personal Working directory:"
	@echo "BASE_PWD:              $(BASE_PWD)"
	@echo "BASE_PWD_BIN:          $(BASE_PWD_BIN)"
	@echo "BASE_PWD_DATA:         $(BASE_PWD_DATA)"
	@echo "BASE_PWD_META:         $(BASE_PWD_META)"
	@echo ""
	@echo "--- base : Current Working folders ---"
	@echo ":"
	@echo "BASE_CWD:               $(BASE_CWD)"
	@echo "BASE_CWD_BIN:           $(BASE_CWD_BIN)"
	@echo "BASE_CWD_DATA:          $(BASE_CWD_DATA)"
	@echo "BASE_CWD_DEP:           $(BASE_CWD_DEP)"
	@echo "BASE_CWD_DEPTMP:        $(BASE_CWD_DEPTMP)"
	@echo "BASE_CWD_META:          $(BASE_CWD_META)"
	@echo "BASE_CWD_SRC:           $(BASE_CWD_SRC)"
	@echo "BASE_CWD_PACK:          $(BASE_CWD_PACK)"
	@echo ""

## base-print
base-print: base-print-pre base-dep-print base-src-print base-bin-print base-run-print

### dep 

# load base stuff we need for CI / CD. These are identical to what will happen inside Github acrtions, so that we can just use a Makefile to bootstrap.


# git 
# Have not done a "dep" step for it yet. Not sure i need it.
BASE_DEP_BIN_GIT_NAME=git
ifeq ($(BASE_OS_NAME),windows)
	BASE_DEP_BIN_GIT_NAME=git.exe
endif
BASE_DEP_BIN_GIT_WHICH=$(shell command -v $(BASE_DEP_BIN_GIT_NAME))
BASE_DEP_BIN_GIT_VERSION_WHICH=$(shell $(BASE_DEP_BIN_GIT_NAME) version)

# gobrew ( via nothing !! )
# https://github.com/kevincobain2000/gobrew/releases/tag/v1.10.10
# https://github.com/kevincobain2000/gobrew/releases/tag/v1.10.10
BASE_DEP_BIN_GOBREW_NAME=gobrew
BASE_DEP_BIN_GOBREW_VERSION=v1.10.12
ifeq ($(BASE_OS_NAME),windows)
	BASE_DEP_BIN_GOBREW_NAME=gobrew.exe
endif
# gobrew version
BASE_DEP_BIN_GOBREW_VERSION_WHICH=$(shell $(BASE_DEP_BIN_GOBREW_NAME) version)
BASE_DEP_BIN_GOBREW_WHICH=$(shell command -v $(BASE_DEP_BIN_GOBREW_NAME))

# go ( via gobrew )
# versions dynamically managed with gobrew.
BASE_DEP_BIN_GO_NAME=go
ifeq ($(BASE_OS_NAME),windows)
	BASE_DEP_BIN_GO_NAME=go.exe
endif
BASE_DEP_BIN_GO_WHICH=$(shell command -v $(BASE_DEP_BIN_GO_NAME))
BASE_DEP_BIN_GO_VERSION_WHICH=$(shell $(BASE_DEP_BIN_GO_NAME) version)



### tinygo

BASE_DEP_BIN_TINYGO_NAME=tinygo
# https://github.com/tinygo-org/tinygo/releases/tag/v0.33.0
# https://github.com/tinygo-org/tinygo/releases/tag/v0.36.0
# https://github.com/tinygo-org/tinygo/releases/tag/v0.37.0
# Drop the V
BASE_DEP_BIN_TINYGO_VERSION=0.37.0

BASE_DEP_TINYGO_META_PREFIX=tinygo

ifeq ($(BASE_OS_NAME),windows)
	BASE_DEP_BIN_TINYGO_NAME=tinygo.exe
endif
BASE_DEP_BIN_TINYGO_VERSION_WHICH=$(shell $(BASE_DEP_BIN_TINYGO_NAME) version)
BASE_DEP_BIN_TINYGO_WHICH=$(shell command -v $(BASE_DEP_BIN_TINYGO_NAME))

# Good reference
# https://github.com/tinygo-org/homebrew-tools/blob/master/tinygo.rb

# ! Depends on binaryen: https://github.com/WebAssembly/binaryen/releases/tag/version_122

# https://github.com/tinygo-org/tinygo/releases/download/v0.36.0/tinygo0.36.0.darwin-amd64.tar.gz
BASE_DEP_BIN_TINYGO_DOWNLOAD_DARWIN_URL=https://github.com/tinygo-org/tinygo/releases/download/v$(BASE_DEP_BIN_TINYGO_VERSION)/tinygo$(BASE_DEP_BIN_TINYGO_VERSION).darwin-arm64.tar.gz
# tinygo0.36.0.darwin-arm64.tar.gz
BASE_DEP_BIN_TINYGO_DOWNLOAD_DARWIN_IMAGE=tinygo$(BASE_DEP_BIN_TINYGO_VERSION).darwin-arm64.tar.gz

# https://github.com/tinygo-org/tinygo/releases/download/v0.36.0/tinygo0.36.0.windows-amd64.zip
BASE_DEP_BIN_TINYGO_DOWNLOAD_WINDOWS_URL=https://github.com/tinygo-org/tinygo/releases/download/v$(BASE_DEP_BIN_TINYGO_VERSION)/tinygo$(BASE_DEP_BIN_TINYGO_VERSION).windows-amd64.zip

# tinygo0.36.0.windows-amd64.zip
BASE_DEP_BIN_TINYGO_DOWNLOAD_WINDOWS_IMAGE=tinygo$(BASE_DEP_BIN_TINYGO_VERSION).windows-amd64.zip


# gojq
# https://github.com/itchyny/gojq
# go install github.com/itchyny/gojq/cmd/gojq@latest
# can curl using wgot ?
# ex: curl "http://localhost:2019/config/" | gojq
# 

## base-dep-print


base-dep-print:
	@echo ""
	@echo "BASE_CWD_DEP:                            $(BASE_CWD_DEP)"
	@echo ""
	@echo "--- dep git ---"
	@echo "BASE_DEP_BIN_GIT_NAME:                   $(BASE_DEP_BIN_GIT_NAME)"
	@echo "BASE_DEP_BIN_GIT_VERSION_WHICH:          $(BASE_DEP_BIN_GIT_VERSION_WHICH)"
	@echo "BASE_DEP_BIN_GIT_WHICH:                  $(BASE_DEP_BIN_GIT_WHICH)"

	@echo ""
	@echo ""
	@echo "--- dep gobrew ---"
	@echo "BASE_DEP_BIN_GOBREW_NAME:                $(BASE_DEP_BIN_GOBREW_NAME)"
	@echo "BASE_DEP_BIN_GOBREW_VERSION:             $(BASE_DEP_BIN_GOBREW_VERSION)"
	@echo "BASE_DEP_BIN_GOBREW_VERSION_WHICH:       $(BASE_DEP_BIN_GOBREW_VERSION_WHICH)"
	@echo "BASE_DEP_BIN_GOBREW_WHICH:               $(BASE_DEP_BIN_GOBREW_WHICH)"
	@echo ""
	@echo ""
	@echo "--- dep golang ---"
	@echo "BASE_DEP_BIN_GO_NAME:                    $(BASE_DEP_BIN_GO_NAME)"
	@echo "BASE_DEP_BIN_GO_VERSION_WHICH:           $(BASE_DEP_BIN_GO_VERSION_WHICH)"
	@echo "BASE_DEP_BIN_GO_WHICH:                   $(BASE_DEP_BIN_GO_WHICH)"
	@echo ""
	@echo ""
	@echo "--- dep tinygo ---"
	@echo "BASE_DEP_BIN_TINYGO_NAME:                $(BASE_DEP_BIN_TINYGO_NAME)"
	@echo "BASE_DEP_BIN_TINYGO_VERSION:             $(BASE_DEP_BIN_TINYGO_VERSION)"
	@echo "BASE_DEP_BIN_TINYGO_VERSION_WHICH:       $(BASE_DEP_BIN_TINYGO_VERSION_WHICH)"
	@echo "BASE_DEP_BIN_TINYGO_WHICH:               $(BASE_DEP_BIN_TINYGO_WHICH)"
	@echo "BASE_DEP_TINYGO_META_PREFIX:             $(BASE_DEP_TINYGO_META_PREFIX)"
	

	@echo "TINYGOROOT:                              $(TINYGOROOT)"
	@echo "BASE_DEP_BIN_TINYGO_INFO:                $(shell $(BASE_DEP_BIN_TINYGO_NAME) info)"
	@echo ""

	@echo ""
	
	
base-dep-init:
	@echo ""
	@echo "base-dep-init - start"
	@echo ""

	mkdir -p $(BASE_CWD_DEP)

	@echo ""
	@echo "base-dep-init - end"
	@echo ""
base-dep-init-del:
	rm -rf $(BASE_CWD_DEP)

base-dep-template: base-dep-init
	@echo ""
	@echo "base-dep-template - start"
	@echo ""

# If same then do nothing.
ifeq ($(BASE_MAKE_IMPORT), $(BASE_CWD_DEP))
	
else
	# copy templates to dep.
	cp -r $(BASE_MAKE_IMPORT)/base* $(BASE_CWD_DEP)
endif


# if BASE_ENV_CI is empty, then we are not inside github action, so do the copy of deps.
#ifeq ($(BASE_ENV_CI), )
	# copy templates to dep.
#cp -r $(BASE_MAKE_IMPORT)/base* $(BASE_CWD_DEP)
#endif

	@echo ""
	@echo "base-dep-template - end"
	@echo ""


## Commented off for now, until we work out things a bit more, to not break developer existing setup of Make, git, go, etc
base-dep-del: #base-dep-git-del base-dep-gobrew-del base-amp-wgot-del base-amp-tree-dep-del base-amp-garble-dep-del base-dep-redress-del

base-dep: base-dep-template
	@echo ""
	@echo " checking base deps"
	@echo ""
	@echo "BASE_DEP_BIN_GIT_WHICH:         $(BASE_DEP_BIN_GIT_WHICH)"
	@echo "BASE_DEP_BIN_GIT_VERSION:       $(shell $(BASE_DEP_BIN_GIT_WHICH) -v)"
	@echo ""
	@echo "BASE_DEP_BIN_GOBREW_WHICH:      $(BASE_DEP_BIN_GOBREW_WHICH)"
	@echo "BASE_DEP_BIN_GOBREW_VERSION:    $(shell $(BASE_DEP_BIN_GOBREW_WHICH) version)"

	#TODO: do checking perhaps...


base-dep-gobrew-del:
	rm -rf $(BASE_DEP_BIN_TINYGO_WHICH)
base-dep-gobrew:
	# Keep as it is. 
	# I have no use case for Users needeing this yet.
	# I might need to make it be installed using wgot later for github CI.

# gobrew
# # https://raw.githubusercontent.com/kevincobain2000/gobrew/v1.10.10/git.io.sh
# curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | sh
	@echo ""
	@echo "gobrew dep :"
	@echo ""
ifeq ($(BASE_OS_NAME),darwin)
	@echo "--- gobrew: darwin ---"
	curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/$(BASE_DEP_BIN_GOBREW_VERSION)/git.io.sh | sh
endif
ifeq ($(BASE_OS_NAME),linux)
	@echo "--- gobrew: linux ---"
	curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/$(BASE_DEP_BIN_GOBREW_VERSION)/git.io.sh | sh
endif
ifeq ($(BASE_OS_NAME),windows)
	@echo "--- gobrew: windows ---"
	curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/$(BASE_DEP_BIN_GOBREW_VERSION)/git.io.sh | sh
	# Below uses Powershell and we dont use powersehll because bash works fine.
	#Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))
endif



## edit deployed files
base-dep-tinygo-edit:
	$(VSCODE_BIN_NAME) $(BASE_CWD_DEP)/$(BASE_DEP_TINYGO_META_PREFIX)*

## list deployed files 
base-dep-tinygo-list:
	$(BASE_AMP_TREE_BIN_NAME) $(BASE_CWD_DEP)/$(BASE_DEP_TINYGO_META_PREFIX)*


base-dep-tinygo-del:
	@echo ""
	@echo "tinygo del ! :"
	@echo ""
ifeq ($(BASE_OS_NAME),darwin)
	@echo "--- tinygo: darwin ---"
	rm -rf $(BASE_CWD_DEP)/tinygo
endif
ifeq ($(BASE_OS_NAME),linux)
	@echo "--- tinygo: linux ---"
	
endif
ifeq ($(BASE_OS_NAME),windows)
	@echo "--- tinygo: windows ---"
	$(BASE_CWD_DEP)/tinygo
endif


base-dep-tinygo: 
	@echo ""
	@echo "tinygo dep :"
	@echo ""
	# TODO: do a dep check, and onyl do it if needed..
ifeq ($(BASE_OS_NAME),darwin)
	$(MAKE) base-dep-tinygo-darwin
endif
ifeq ($(BASE_OS_NAME),linux)
	$(MAKE) base-dep-tinygo-linux
endif
ifeq ($(BASE_OS_NAME),windows)
	$(MAKE) base-dep-tinygo-windows
endif


base-dep-tinygo-all: base-dep-tinygo-darwin base-dep-tinygo-linux base-dep-tinygo-windows

base-dep-tinygo-darwin:
	@echo "--- tinygo: darwin ---"
	# https://tinygo.org/getting-started/install/macos/

	# grab
	$(MAKE) ARG=$(BASE_DEP_BIN_TINYGO_DOWNLOAD_DARWIN_URL) base-amp-wgot-run 
	
	# decompress to BASE_CWD_DEPTMP
	$(MAKE) ARG=$(BASE_DEP_BIN_TINYGO_DOWNLOAD_DARWIN_IMAGE) base-amp-arc-run-decompress

	# Print contents
	$(BASE_AMP_TREE_BIN_NAME) -h $(BASE_CWD_DEPTMP)*
	
	# move it ALL across 
	mv $(BASE_CWD_DEPTMP)/tinygo $(BASE_CWD_DEP)/tinygo
	# OR
	# move ONLY the bin file into shared .deps thats always mapped
	#mv $(BASE_CWD_DEPTMP)/tinygo/bin/tinygo $(BASE_CWD_DEP)/

	rm -rf $(BASE_CWD_DEPTMP)

	rm -f $(BASE_DEP_BIN_TINYGO_DOWNLOAD_DARWIN_IMAGE)

base-dep-tinygo-linux:
	@echo "--- tinygo: linux ---"
	# https://tinygo.org/getting-started/install/linux/
	# SAME as we do for darwin
	@echo "NOT SUPPORTED: Contact support if you need this and we will add it."


base-dep-tinygo-windows:
	@echo "--- tinygo: windows ---"
	# https://tinygo.org/getting-started/install/windows/
	# https://github.com/tinygo-org/tinygo/releases/download/v0.36.0/tinygo0.36.0.windows-amd64.zip

	# grab
	#$(MAKE) ARG=$(BASE_DEP_BIN_TINYGO_DOWNLOAD_WINDOWS_URL) base-amp-wgot-run 
	
	# decompress to BASE_CWD_DEPTMP
	$(MAKE) ARG=$(BASE_DEP_BIN_TINYGO_DOWNLOAD_WINDOWS_IMAGE) base-amp-arc-run-decompress

	# move it ALL across 
	mv $(BASE_CWD_DEPTMP)/tinygo $(BASE_CWD_DEP)/tinygo
	# OR
	# move ONLY the bin file into shared .deps thats always mapped
	#mv $(BASE_CWD_DEPTMP)/tinygo/bin/tinygo $(BASE_CWD_DEP)/

	rm -rf $(BASE_CWD_DEPTMP)

	rm -rf $(BASE_DEP_BIN_TINYGO_DOWNLOAD_WINDOWS_IMAGE)





### amp ( our core deps )

base-amp-print:
	@echo ""
	@echo "-- amp --"
	@echo "-- chronologically --"
	@echo ""
	@echo "BASE_AMP_ARC_BIN_NAME:          $(BASE_AMP_ARC_BIN_NAME)"
	@echo "BASE_AMP_ARC_BIN_WHICH:         $(BASE_AMP_ARC_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_ARC_COM_BIN_NAME:      $(BASE_AMP_ARC_COM_BIN_NAME)"
	@echo "BASE_AMP_ARC_COM_BIN_WHICH:     $(BASE_AMP_ARC_COM_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_GARBLE_BIN_NAME:       $(BASE_AMP_GARBLE_BIN_NAME)"
	@echo "BASE_AMP_GARBLE_BIN_WHICH:      $(BASE_AMP_GARBLE_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_REDRESS_BIN_NAME:      $(BASE_AMP_REDRESS_BIN_NAME)"
	@echo "BASE_AMP_REDRESS_BIN_WHICH:     $(BASE_AMP_REDRESS_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_TREE_BIN_NAME:         $(BASE_AMP_TREE_BIN_NAME)"
	@echo "BASE_AMP_TREE_BIN_WHICH:        $(BASE_AMP_TREE_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_REPLACE_BIN_NAME:      $(BASE_AMP_REPLACE_BIN_NAME)"
	@echo "BASE_AMP_REPLACE_BIN_WHICH:     $(BASE_AMP_REPLACE_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_WGOT_BIN_NAME:         $(BASE_AMP_WGOT_BIN_NAME)"
	@echo "BASE_AMP_WGOT_BIN_WHICH:        $(BASE_AMP_WGOT_BIN_WHICH)"
	@echo ""
	@echo "BASE_AMP_WHICH_BIN_NAME:        $(BASE_AMP_WHICH_BIN_NAME)"
	@echo "BASE_AMP_WHICH_BIN_WHICH:       $(BASE_AMP_WHICH_BIN_WHICH)"
	@echo ""
	

base-amp-dep-all-del:
	# smoke test that it all works
	$(MAKE) base-amp-arc-dep-del
	$(MAKE) base-amp-arc-com-dep-del
	$(MAKE) base-amp-garble-dep-del
	$(MAKE) base-amp-gh-dep-del
	$(MAKE) base-amp-redress-dep-del
	$(MAKE) base-amp-replace-dep-del
	$(MAKE) base-amp-tree-dep-del
	$(MAKE) base-amp-which-dep-del
	$(MAKE) base-amp-wgot-dep-del

base-amp-dep-all:
	# smoke test that it all works
	$(MAKE) base-amp-arc-dep-single
	$(MAKE) base-amp-arc-com-dep-single
	$(MAKE) base-amp-garble-dep-single
	$(MAKE) base-amp-gh-dep-single
	$(MAKE) base-amp-redress-dep-single
	$(MAKE) base-amp-replace-dep-single
	$(MAKE) base-amp-tree-dep-single
	$(MAKE) base-amp-which-dep-single
	$(MAKE) base-amp-wgot-dep-single





#### amp_arc

# https://github.com/jm33-m0/arc
# https://github.com/gedw99/arc
BASE_AMP_ARC_BIN_NAME=amp_arc
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_ARC_BIN_NAME=amp_arc.exe
endif
# https://github.com/jm33-m0/arc/releases/tag/v1.0.2
# https://github.com/jm33-m0/arc/releases/tag/v2.0.0
BASE_AMP_ARC_BIN_VERSION=v1.0.2
#BASE_AMP_ARC_BIN_VERSION=v2.0.0
BASE_AMP_ARC_BIN_VERSION_WHICH=$(shell $(BASE_AMP_ARC_BIN_NAME) version)
BASE_AMP_ARC_BIN_WHICH=$(shell command -v $(BASE_AMP_ARC_BIN_NAME))

base-amp-arc-dep-print:
	@echo ""
	@echo "--- dep arc ---"
	@echo "BASE_AMP_ARC_BIN_NAME:             $(BASE_AMP_ARC_BIN_NAME)"
	@echo "BASE_AMP_ARC_BIN_VERSION:          $(BASE_AMP_ARC_BIN_VERSION)"
	@echo "BASE_AMP_ARC_BIN_VERSION_WHICH:    $(BASE_AMP_ARC_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_ARC_BIN_WHICH:            $(BASE_AMP_ARC_BIN_WHICH)"
	@echo ""
base-amp-arc-dep-del:
	@echo ""
	@echo "deleting $(BASE_AMP_ARC_BIN_NAME) dep ..."
	rm -f $(BASE_AMP_ARC_BIN_WHICH)
	@echo ""
base-amp-arc-dep:
ifeq ($(BASE_AMP_ARC_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_ARC_BIN_NAME) dep check: failed, so installing ..."
	$(MAKE) base-amp-arc-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_ARC_BIN_NAME) dep check: passed"
endif

base-amp-arc-dep-single:
	go install github.com/jm33-m0/arc/cmd/arc@$(BASE_AMP_ARC_BIN_VERSION)
	mv $(GOPATH)/bin/arc $(GOPATH)/bin/$(BASE_AMP_ARC_BIN_NAME) 

BASE_AMP_ARC_RUN_VAR_FILE=?file?

base-amp-arc-run-print:
	@echo ""
	@echo "BASE_AMP_ARC_RUN_VAR_FILE:    $(BASE_AMP_ARC_RUN_VAR_FILE)"
	@echo ""

base-amp-arc-run-h: base-amp-arc-dep
	$(BASE_AMP_ARC_BIN_NAME) -h

base-amp-arc-run-compress: base-amp-arc-dep
	$(BASE_AMP_ARC_BIN_NAME) -h

base-amp-arc-run-decompress: base-amp-arc-dep
	# ARG=file to decompress 
	$(BASE_AMP_ARC_BIN_NAME) -x -f $(ARG) $(BASE_CWD_DEPTMP)





#### amp_arc_compressor

# https://github.com/jm33-m0/arc
# https://github.com/gedw99/arc
BASE_AMP_ARC_COM_BIN_NAME=amp_arc_compressor
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_ARC_COM_BIN_NAME=amp_arc_compressor.exe
endif
# https://github.com/jm33-m0/arc/releases/tag/v1.0.2
BASE_AMP_ARC_COM_BIN_VERSION=v1.0.2
BASE_AMP_ARC_COM_BIN_VERSION_WHICH=$(shell $(BASE_AMP_ARC_COM_BIN_NAME) version)
BASE_AMP_ARC_COM_BIN_WHICH=$(shell command -v $(BASE_AMP_ARC_COM_BIN_NAME))

base-amp-arc-com-dep-print:
	@echo ""
	@echo "--- dep arc ---"
	@echo "BASE_AMP_ARC_COM_BIN_NAME:             $(BASE_AMP_ARC_COM_BIN_NAME)"
	@echo "BASE_AMP_ARC_COM_BIN_VERSION:          $(BASE_AMP_ARC_COM_BIN_VERSION)"
	@echo "BASE_AMP_ARC_COM_BIN_VERSION_WHICH:    $(BASE_AMP_ARC_COM_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_ARC_COM_BIN_WHICH:            $(BASE_AMP_ARC_COM_BIN_WHICH)"
	@echo ""
base-amp-arc-com-dep-del:
	@echo ""
	@echo "deleting $(BASE_AMP_ARC_COM_BIN_NAME) dep ..."
	rm -f $(BASE_AMP_ARC_COM_BIN_WHICH)
	@echo ""
base-amp-arc-com-dep:
ifeq ($(BASE_AMP_ARC_COM_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_ARC_COM_BIN_NAME) dep check: failed, so installing ..."
	$(MAKE) base-amp-arc-com-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_ARC_COM_BIN_NAME) dep check: passed"
endif

base-amp-arc-com-dep-single:
	go install github.com/jm33-m0/arc/cmd/compressor@$(BASE_AMP_ARC_COM_BIN_VERSION)
	mv $(GOPATH)/bin/compressor $(GOPATH)/bin/$(BASE_AMP_ARC_COM_BIN_NAME) 

BASE_AMP_ARC_COM_RUN_VAR_FILE=?file?

base-amp-arc-com-run-print:
	@echo ""
	@echo "BASE_AMP_ARC_COM_RUN_VAR_FILE:    $(BASE_AMP_ARC_COM_RUN_VAR_FILE)"
	@echo ""
base-amp-arc-com-run-h: base-amp-arc-com-dep
	$(BASE_AMP_ARC_COM_BIN_NAME) -h
base-amp-arc-com-run-compress-brotli: base-amp-arc-com-dep
	@echo ""
	@echo "br compress ..."
	$(BASE_AMP_ARC_COM_BIN_NAME) -t br -c $(BASE_AMP_ARC_COM_RUN_VAR_FILE) -o $(BASE_AMP_ARC_COM_RUN_VAR_FILE).br
base-amp-arc-com-run-decompress-brotli: base-amp-arc-com-dep
	@echo ""
	@echo "br decompress ..."
	$(BASE_AMP_ARC_COM_BIN_NAME) -t br -f $(BASE_AMP_ARC_COM_RUN_VAR_FILE).br -o $(BASE_AMP_ARC_COM_RUN_VAR_FILE)



#### garble

# https://github.com/burrowers/garble
BASE_AMP_GARBLE_BIN_NAME=amp_garble
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_GARBLE_BIN_NAME=amp_garble.exe
endif
# https://github.com/burrowers/garble/releases/tag/v0.13.0
# https://github.com/burrowers/garble/releases/tag/v0.14.1
BASE_AMP_GARBLE_BIN_VERSION=v0.14.1
BASE_AMP_GARBLE_BIN_VERSION_WHICH=$(shell $(BASE_AMP_GARBLE_BIN_NAME) version)
BASE_AMP_GARBLE_BIN_WHICH=$(shell command -v $(BASE_AMP_GARBLE_BIN_NAME))

base-amp-garble-dep-print:
	@echo ""
	@echo "--- dep garble ---"
	@echo "BASE_AMP_GARBLE_BIN_NAME:            $(BASE_AMP_GARBLE_BIN_NAME)"
	@echo "BASE_AMP_GARBLE_BIN_VERSION:         $(BASE_AMP_GARBLE_BIN_VERSION)"
	@echo "BASE_AMP_GARBLE_BIN_VERSION_WHICH:   $(BASE_AMP_GARBLE_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_GARBLE_BIN_WHICH:           $(BASE_AMP_GARBLE_BIN_WHICH)"
	@echo ""
base-amp-garble-dep-del:
	rm -rf $(BASE_AMP_GARBLE_BIN_WHICH)
base-amp-garble-dep:
	@echo ""
	@echo "garble dep"
	@echo ""
ifeq ($(BASE_AMP_GARBLE_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_GARBLE_BIN_NAME) check: failed"
	$(MAKE) base-amp-garble-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_GARBLE_BIN_NAME) check: passed"
endif

base-amp-garble-dep-single:
	$(BASE_DEP_BIN_GO_NAME) install mvdan.cc/garble@$(BASE_AMP_GARBLE_BIN_VERSION)
	mv $(GOPATH)/bin/garble $(GOPATH)/bin/$(BASE_AMP_GARBLE_BIN_NAME)


### github cli

# https://github.com/cli/cli
BASE_AMP_GH_BIN_NAME=amp_gh
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_GH_BIN_NAME=amp_gh.exe
endif
# https://github.com/cli/cli/releases/tag/v2.55.0
BASE_AMP_GH_BIN_VERSION=v2.55.0
BASE_AMP_GH_BIN_VERSION_WHICH=$(shell $(BASE_AMP_GH_BIN_NAME) version)
BASE_AMP_GH_BIN_WHICH=$(shell command -v $(BASE_AMP_GH_BIN_NAME))

base-amp-gh-dep-print:
	@echo ""
	@echo "--- dep garble ---"
	@echo "BASE_AMP_GH_BIN_NAME:            $(BASE_AMP_GH_BIN_NAME)"
	@echo "BASE_AMP_GH_BIN_VERSION:         $(BASE_AMP_GH_BIN_VERSION)"
	@echo "BASE_AMP_GH_BIN_VERSION_WHICH:   $(BASE_AMP_GH_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_GH_BIN_WHICH:           $(BASE_AMP_GH_BIN_WHICH)"
	@echo ""
base-amp-gh-dep-del:
	rm -rf $(BASE_AMP_GH_BIN_WHICH)
base-amp-gh-dep:
	@echo ""
	@echo "garble dep"
	@echo ""
ifeq ($(BASE_AMP_GH_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_GH_BIN_NAME) check: failed"
	$(MAKE) base-amp-garble-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_GH_BIN_NAME) check: passed"
endif

base-amp-gh-dep-single:
	$(BASE_DEP_BIN_GO_NAME) install github.com/cli/cli/v2/cmd/gh@$(BASE_AMP_GH_BIN_VERSION)
	mv $(GOPATH)/bin/gh $(GOPATH)/bin/$(BASE_AMP_GH_BIN_NAME)


#### redress

## base-amp-redress-run is a basic decompiler
# https://github.com/goretk/redress
BASE_AMP_REDRESS_BIN_NAME=amp_redress
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_REDRESS_BIN_NAME=amp_redress.exe
endif
# https://github.com/goretk/redress/releases/tag/v1.2.0
# https://github.com/goretk/redress/releases/tag/v1.2.8
# https://github.com/goretk/redress/releases/tag/v1.2.14
BASE_AMP_REDRESS_BIN_VERSION=v1.2.14
BASE_AMP_REDRESS_BIN_VERSION_WHICH=$(shell $(BASE_AMP_REDRESS_BIN_NAME) version)
BASE_AMP_REDRESS_BIN_WHICH=$(shell command -v $(BASE_AMP_REDRESS_BIN_NAME))

base-amp-redress-dep-print:
	@echo ""
	@echo "--- dep redress ---"
	@echo "BASE_AMP_REDRESS_BIN_NAME:           $(BASE_AMP_REDRESS_BIN_NAME)"
	@echo "BASE_AMP_REDRESS_BIN_VERSION:        $(BASE_AMP_REDRESS_BIN_VERSION)"
	@echo "BASE_AMP_REDRESS_BIN_VERSION_WHICH:  $(BASE_AMP_REDRESS_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_REDRESS_BIN_WHICH:          $(BASE_AMP_REDRESS_BIN_WHICH)"
	@echo ""
base-amp-redress-dep-del:
	rm -rf $(BASE_AMP_REDRESS_BIN_WHICH)
base-amp-redress-dep:
	@echo ""
	@echo "redress dep"
	@echo ""
ifeq ($(BASE_AMP_REDRESS_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_REDRESS_BIN_NAME) check: failed"
	$(MAKE) base-amp-redress-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_REDRESS_BIN_NAME) check: passed"
endif

base-amp-redress-dep-single:
	$(BASE_DEP_BIN_GO_NAME) install github.com/goretk/redress@$(BASE_AMP_REDRESS_BIN_VERSION)
	mv $(GOPATH)/bin/redress $(GOPATH)/bin/$(BASE_AMP_REDRESS_BIN_NAME)


# default to the BASE_RUN
BASE_AMP_REDRESS_RUN_VAR_BINARY=$(BASE_RUN)

base-amp-redress-run-print:
	@echo ""
	@echo "BASE_AMP_REDRESS_RUN_VAR_BINARY: $(BASE_AMP_REDRESS_RUN_VAR_BINARY)"
	@echo ""
	@echo ""

base-amp-redress-run-h:
	$(BASE_AMP_REDRESS_BIN_NAME) -h

	# github.com/mandiant/GoReSym is pure go and another inspector. It also gets shot back when garbles :)
	# go install github.com/mandiant/GoReSym@latetst
	# GoReSym -t -d -p ./gophemeral_bin_darwin_arm64

	# https://github.com/radareorg/radare2 works with REDRESS and looks freaking useful.


base-amp-redress-run: base-amp-redress-dep
	@echo ""
	@echo "using redress to see how well its obfuscates ..."
	@echo ""
	@echo "-- version ---"
	@echo ""
	$(BASE_AMP_REDRESS_BIN_NAME) version $(BASE_AMP_REDRESS_RUN_VAR_BINARY)
	@echo ""
	@echo "-- info ---"
	@echo ""
	$(BASE_AMP_REDRESS_BIN_NAME) info $(BASE_AMP_REDRESS_RUN_VAR_BINARY)
	@echo ""
	@echo "-- packages  --std --vendor ---"
	@echo ""
	$(BASE_AMP_REDRESS_BIN_NAME) packages $(BASE_AMP_REDRESS_RUN_VAR_BINARY) --std --vendor
	@echo ""
	@echo "-- moduledata ---"
	@echo ""
	$(BASE_AMP_REDRESS_BIN_NAME) moduledata $(BASE_AMP_REDRESS_RUN_VAR_BINARY)



#### replace

# https://github.com/webdevops/go-replace
# https://github.com/nuvolaris/go-replace is BETTER
# https://github.com/nuvolaris/go-replace/network/dependents uses it for a FULL DEV tooling !

# Does find replace on text. Really useful.
#go install github.com/webdevops/go-replace@latest
BASE_AMP_REPLACE_BIN_NAME=amp_replace
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_REPLACE_BIN_NAME=amp_replace.exe
endif
# https://github.com/webdevops/go-replace/releases/tag/22.10.0
BASE_AMP_REPLACE_BIN_VERSION=22.10.0
BASE_AMP_REPLACE_BIN_VERSION_WHICH=$(shell $(BASE_AMP_REPLACE_BIN_NAME) --version)
BASE_AMP_REPLACE_BIN_WHICH=$(shell command -v $(BASE_AMP_REPLACE_BIN_NAME))

base-amp-replace-dep-print:
	@echo ""
	@echo "--- dep replace ---"
	@echo "BASE_AMP_REPLACE_BIN_NAME:             $(BASE_AMP_REPLACE_BIN_NAME)"
	@echo "BASE_AMP_REPLACE_BIN_VERSION:          $(BASE_AMP_REPLACE_BIN_VERSION)"
	@echo "BASE_AMP_REPLACE_BIN_VERSION_WHICH:    $(BASE_AMP_REPLACE_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_REPLACE_BIN_WHICH:            $(BASE_AMP_REPLACE_BIN_WHICH)"
	@echo ""
base-amp-replace-dep-del:
	rm -rf $(BASE_AMP_REPLACE_BIN_WHICH)
base-amp-replace-dep:
	@echo ""
	@echo "replace dep"
	@echo ""
ifeq ($(BASE_AMP_REPLACE_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_REPLACE_BIN_NAME) check: failed"
	$(MAKE) base-amp-replace-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_REPLACE_BIN_NAME) check: passed"
endif

base-amp-replace-dep-single:
	$(BASE_DEP_BIN_GO_NAME) install github.com/webdevops/go-replace@$(BASE_AMP_REPLACE_BIN_VERSION)
	mv $(GOPATH)/bin/go-replace $(GOPATH)/bin/$(BASE_AMP_REPLACE_BIN_NAME)

# --search
BASE_AMP_REPLACE_RUN_VAR_SEARCH=?search?
# --replace
BASE_AMP_REPLACE_RUN_VAR_REPLACE=?replace?

# --path=   use files in this path
BASE_AMP_REPLACE_RUN_VAR_PATH=?file-pattern?
# --path-pattern  file pattern (* for wildcard, only basename of file)
BASE_AMP_REPLACE_RUN_VAR_FILE_PATTERN=?file-pattern?

BASE_AMP_REPLACE_RUN_VAR_FILE=?file?

base-amp-replace-run-print:
	@echo ""
	@echo "BASE_AMP_REPLACE_RUN_VAR_SEARCH:         $(BASE_AMP_REPLACE_RUN_VAR_SEARCH)"
	@echo "BASE_AMP_REPLACE_RUN_VAR_REPLACE:        $(BASE_AMP_REPLACE_RUN_VAR_REPLACE)"
	@echo "BASE_AMP_REPLACE_RUN_VAR_FILE_PATTERN:   $(BASE_AMP_REPLACE_RUN_VAR_FILE_PATTERN)"

	@echo "BASE_AMP_REPLACE_RUN_VAR_FILE:           $(BASE_AMP_REPLACE_RUN_VAR_FILE)"
	@echo ""
base-amp-replace-run-h: base-amp-replace-dep
	$(BASE_AMP_REPLACE_BIN_NAME) -h

	
base-amp-replace-run-dry: base-amp-replace-dep
	$(BASE_AMP_REPLACE_BIN_NAME) --dry-run --search $(BASE_AMP_REPLACE_RUN_VAR_SEARCH) --replace $(BASE_AMP_REPLACE_RUN_VAR_REPLACE) --path-pattern $(BASE_AMP_REPLACE_RUN_VAR_FILE_PATTERN)  $(BASE_AMP_REPLACE_RUN_VAR_FILE)
base-amp-replace-run: base-amp-replace-dep
	$(BASE_AMP_REPLACE_BIN_NAME) --search $(BASE_AMP_REPLACE_RUN_VAR_SEARCH) --replace $(BASE_AMP_REPLACE_RUN_VAR_REPLACE) $(BASE_AMP_REPLACE_RUN_VAR_FILE)


#### tree

# https://github.com/a8m/tree
# https://github.com/a8m/s3tree is for S3. Might be useful.
# Use https://github.com/kk17/s3tree instead, https://github.com/gedw99/s3tree
BASE_AMP_TREE_BIN_NAME=amp_tree
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_TREE_BIN_NAME=amp_tree.exe
endif
BASE_AMP_TREE_BIN_VERSION=latest
BASE_AMP_TREE_BIN_VERSION_WHICH=$(shell $(BASE_AMP_TREE_BIN_NAME) version)
BASE_AMP_TREE_BIN_WHICH=$(shell command -v $(BASE_AMP_TREE_BIN_NAME))

base-amp-tree-dep-print:
	@echo ""
	@echo "--- dep tree ---"
	@echo "BASE_AMP_TREE_BIN_NAME:              $(BASE_AMP_TREE_BIN_NAME)"
	@echo "BASE_AMP_TREE_BIN_VERSION:           $(BASE_AMP_TREE_BIN_VERSION)"
	@echo "BASE_AMP_TREE_BIN_VERSION_WHICH:     $(BASE_AMP_TREE_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_TREE_BIN_WHICH:             $(BASE_AMP_TREE_BIN_WHICH)"
	@echo ""
base-amp-tree-dep-del:
	rm -rf $(BASE_AMP_TREE_BIN_WHICH)
base-amp-tree-dep:
ifeq ($(BASE_AMP_TREE_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_TREE_BIN_NAME) dep check: failed, so installing ..."
	@echo ""
	$(MAKE) base-amp-tree-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_TREE_BIN_NAME) dep check: passed"
	@echo ""
endif

base-amp-tree-dep-single:
	$(BASE_DEP_BIN_GO_NAME) install github.com/a8m/tree/cmd/tree@$(BASE_AMP_TREE_BIN_VERSION)
	mv $(GOPATH)/bin/tree $(GOPATH)/bin/$(BASE_AMP_TREE_BIN_NAME) 



#### which

# go-which
# https://github.com/hairyhenderson/go-which
# go install github.com/hairyhenderson/go-which/cmd/which
# Replaces "which" and "command -v"
BASE_AMP_WHICH_BIN_NAME=amp_which
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_WHICH_BIN_NAME=amp_which.exe
endif
# https://github.com/hairyhenderson/go-which/releases/tag/v0.2.0
BASE_AMP_WHICH_BIN_VERSION=v0.2.0
BASE_AMP_WHICH_BIN_VERSION_WHICH=$(shell $(BASE_AMP_WHICH_BIN_NAME) version)
BASE_AMP_WHICH_BIN_WHICH=$(shell command -v $(BASE_AMP_WHICH_BIN_NAME))

base-amp-which-dep-print:
	@echo ""
	@echo "--- dep which ---"
	@echo "BASE_AMP_WHICH_BIN_NAME:             $(BASE_AMP_WHICH_BIN_NAME)"
	@echo "BASE_AMP_WHICH_BIN_VERSION:          $(BASE_AMP_WHICH_BIN_VERSION)"
	@echo "BASE_AMP_WHICH_BIN_VERSION_WHICH:    $(BASE_AMP_WHICH_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_WHICH_BIN_WHICH:            $(BASE_AMP_WHICH_BIN_WHICH)"
	@echo ""
base-amp-which-dep-del:
	rm -rf $(BASE_AMP_WHICH_BIN_WHICH)
base-amp-which-dep:
ifeq ($(BASE_AMP_WHICH_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_WHICH_BIN_NAME) dep check: failed, so installing ..."
	@echo ""
	$(MAKE) base-amp-which-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_ARC_COM_BIN_NAME) dep check: passed"
	@echo ""
endif

base-amp-which-dep-single:
	# https://github.com/hairyhenderson/go-which/tree/main/cmd
	$(BASE_DEP_BIN_GO_NAME) install github.com/hairyhenderson/go-which/cmd/which@$(BASE_AMP_WHICH_BIN_VERSION)
	mv $(GOPATH)/bin/which $(GOPATH)/bin/$(BASE_AMP_WHICH_BIN_NAME) 

BASE_AMP_WHICH_RUN_VAR_FILE=?file?

base-amp-which-run-print:
	@echo ""
	@echo "BASE_AMP_WHICH_RUN_VAR_FILE:    $(BASE_AMP_WHICH_RUN_VAR_FILE)"
	@echo ""
base-amp-which-run-h: base-amp-which-dep
	$(BASE_AMP_WHICH_BIN_NAME) -h
base-amp-which-run: base-amp-which-dep
	$(BASE_AMP_WHICH_BIN_NAME)  $(BASE_AMP_WHICH_RUN_VAR_FILE)



#### wgot 

# like wget, got cross platform.
# https://github.com/bitrise-io/got
BASE_AMP_WGOT_BIN_NAME=amp_wgot
ifeq ($(BASE_OS_NAME),windows)
	BASE_AMP_WGOT_BIN_NAME=amp_wgot.exe
endif
# no tags
# https://github.com/bitrise-io/got/blob/master/go.mod
BASE_AMP_WGOT_BIN_VERSION=latest
BASE_AMP_WGOT_BIN_VERSION_WHICH=$(shell $(BASE_AMP_WGOT_BIN_NAME) version)
BASE_AMP_WGOT_BIN_WHICH=$(shell command -v $(BASE_AMP_WGOT_BIN_NAME))

base-amp-wgot-dep-print:
	@echo ""
	@echo "--- dep wgot ---"
	@echo "BASE_AMP_WGOT_BIN_NAME:              $(BASE_AMP_WGOT_BIN_NAME)"
	@echo "BASE_AMP_WGOT_BIN_VERSION:           $(BASE_AMP_WGOT_BIN_VERSION)"
	@echo "BASE_AMP_WGOT_BIN_VERSION_WHICH:     $(BASE_AMP_WGOT_BIN_VERSION_WHICH)"
	@echo "BASE_AMP_WGOT_BIN_WHICH:             $(BASE_AMP_WGOT_BIN_WHICH)"
	@echo ""
base-amp-wgot-dep-del:
	rm -rf $(BASE_AMP_WGOT_BIN_WHICH)
base-amp-wgot-dep:
ifeq ($(BASE_AMP_WGOT_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_AMP_WGOT_BIN_NAME) dep check: failed, so installing ..."
	@echo ""
	$(MAKE) base-amp-wgot-dep-single
else
	@echo ""
	@echo "$(BASE_AMP_WGOT_BIN_NAME) dep check: passed"
	@echo ""
endif
base-amp-wgot-dep-single:
	$(BASE_DEP_BIN_GO_NAME) install github.com/bitrise-io/got/cmd/wgot@$(BASE_AMP_WGOT_BIN_VERSION)
	mv $(GOPATH)/bin/wgot $(GOPATH)/bin/$(BASE_AMP_WGOT_BIN_NAME) 

base-amp-wgot-run-h: base-amp-wgot-dep
	$(BASE_AMP_WGOT_BIN_NAME) -h


base-amp-wgot-run: base-amp-wgot-dep
	# download the files using ARG passed
	$(BASE_AMP_WGOT_BIN_NAME) $(ARG)

	

	
base-dep-gio:
	# GIO CMD
	# https://github.com/gioui/gio-cmd/releases/tag/v0.7.0
	# GIO code
	# https://github.com/gioui/gio/releases/tag/v0.7.0
	# GIO-X code
	# https://github.com/gioui/gio-x/releases/tag/v0.7.0

	@echo ""
	@echo "gio deps:"
ifeq ($(BASE_OS_NAME),darwin)
	@echo "--- gio: darwin ---"
	#curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | sh
endif
ifeq ($(BASE_OS_NAME),linux)
	@echo "--- gio: linux ---"
	#https://gioui.org/doc/install/linux
	#sudo apt install gcc pkg-config libwayland-dev libx11-dev libx11-xcb-dev libxkbcommon-x11-dev libgles2-mesa-dev libegl1-mesa-dev libffi-dev libxcursor-dev libvulkan-dev
endif
ifeq ($(BASE_OS_NAME),windows)
	@echo "--- gio: windows ---"
	#curl -sLk https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.sh | sh
	# Below uses Powershell and we dont use powersehll because bash works fine.
	#Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/kevincobain2000/gobrew/master/git.io.ps1'))
endif





### src

BASE_SRC_NAME=.
BASE_SRC_URL=???
BASE_SRC_VERSION=latest
# !!! This is also givng me a Fatal ?
BASE_SRC_VERSION_WHICH=$(shell cd $(BASE_SRC) && $(BASE_SRC_VERSION_CMD))
# git rev-parse --short HEAD
BASE_SRC_VERSION_CMD=$(BASE_DEP_BIN_GIT_NAME) rev-parse HEAD
#BASE_SRC_VERSION_CMD=$(BASE_DEP_BIN_GIT_NAME) rev-parse --short HEAD
# !!! This is causing the "fatal", because $BASE_SRC is NOT always there.
# We can make this better.
BASE_SRC_TAG_WHICH=$(shell cd $(BASE_SRC) && $(BASE_SRC_TAG_CMD))
# $(shell git describe --abbrev=0 --tags)
BASE_SRC_TAG_CMD=$(BASE_DEP_BIN_GIT_NAME) describe --abbrev=0 --tags


BASE_SRC=$(BASE_CWD_SRC)/$(BASE_SRC_NAME)
#BASE_SRC_WHICH=$(shell test -d $(BASE_SRC))
# WORKS !
#BASE_SRC_WHICH=$(wildcard $(BASE_SRC))
#BASE_SRC_WHICH=$(shell $(BASE_AMP_WHICH_BIN_NAME) $(BASE_SRC))
BASE_SRC_WHICH=$(shell fsize $(BASE_SRC))

#BASE_SRC_NAME=go-brew
#BASE_SRC_URL=https://github.com/kevincobain2000/gobrew
#BASE_SRC_VERSION=v1.2
#BASE_SRC=$(BASE_CWD_SRC)/$(BASE_SRC_NAME)

BASE_SRC_CMD=cd $(BASE_SRC) &&

## base-src-print
base-src-print:
	@echo ""
	@echo "--- base ---"
	@echo "BASE_CWD_SRC:            $(BASE_CWD_SRC)"
	@echo ""
	@echo "--- src ---"
	@echo "BASE_SRC_NAME:           $(BASE_SRC_NAME)"
	@echo "BASE_SRC_URL:            $(BASE_SRC_URL)"
	@echo "BASE_SRC_VERSION:        $(BASE_SRC_VERSION)"
	@echo "BASE_SRC_VERSION_WHICH:  $(BASE_SRC_VERSION_WHICH)"
	@echo "BASE_SRC_VERSION_CMD:    $(BASE_SRC_VERSION_CMD)"
	@echo "BASE_SRC_TAG_WHICH:      $(BASE_SRC_TAG_WHICH)"
	@echo "BASE_SRC_TAG_CMD:        $(BASE_SRC_TAG_CMD)"
	
	@echo ""
	@echo "--- src calculated ---"
	@echo "BASE_SRC:             $(BASE_SRC)"
	@echo "BASE_SRC_WHICH:       $(BASE_SRC_WHICH)"
	@echo "BASE_SRC_CMD:         $(BASE_SRC_CMD)"
	@echo ""

	@echo ""
	@echo "--- src signing ---"
	@echo ""
	@echo "- bin" 
	@echo "BASE_SRC_SIGNING_BIN_NAME:      $(BASE_SRC_SIGNING_BIN_NAME)"
	@echo "BASE_SRC_SIGNING_BIN_WHICH:     $(BASE_SRC_SIGNING_BIN_WHICH)"
	@echo "" 
	@echo "" 
	@echo "- config" 
	@echo "BASE_SRC_SIGNING_CONFIG:        $(BASE_SRC_SIGNING_CONFIG)"
	@echo "BASE_SRC_SIGNING_CONFIG_LOCAL:  $(BASE_SRC_SIGNING_CONFIG_LOCAL)"
	@echo "" 
	@echo "" 
	@echo "- var" 
	@echo "" 
	@echo "BASE_SRC_SIGNING_USER_NAME:     $(BASE_SRC_SIGNING_USER_NAME)"
	@echo "BASE_SRC_SIGNING_USER_EMAIL:    $(BASE_SRC_SIGNING_USER_EMAIL)"
	@echo "" 

	@echo "BASE_SRC_SIGNING_KEY_CONFIG:    $(HOME)/.ssh/config"
	@echo "BASE_SRC_SIGNING_KEY_PRIV:      $(BASE_SRC_SIGNING_KEY_PRIV)"
	@echo "BASE_SRC_SIGNING_KEY:           $(BASE_SRC_SIGNING_KEY)"
	@echo "BASE_SRC_SIGNING_PROGRAM:       $(BASE_SRC_SIGNING_PROGRAM)"
	@echo "BASE_SRC_SIGNING_FORMAT:        $(BASE_SRC_SIGNING_FORMAT)"
	@echo ""


## base-src-ls-h
base-src-ls-h:
	$(BASE_AMP_TREE_BIN_NAME) --help
	# -d		Directories only.
	# -h 		Directories and files.
	# -P		List only those files that match the pattern given.
	# -C		Colour turned on.
	  # PURPLE = folder
	  # GREEN = executable ( .sh )
## base-src-ls
base-src-ls:
	# only folders, no files.
	# if not there does not blow up at all.
	$(BASE_AMP_TREE_BIN_NAME) -d -C $(BASE_CWD_SRC)
## base-src-ls-all
base-src-ls-all:
	# all files and folders
	$(BASE_AMP_TREE_BIN_NAME) -h -C $(BASE_CWD_SRC)


## Delete this LATER ..
base-src-init-deep:
	# just the .src/name for my testing
	mkdir -p $(BASE_SRC)
	touch $(BASE_SRC)/test.md

## base-src-init
base-src-init:
	# just the .src
	@echo "Creating .src ..."
	mkdir -p $(BASE_CWD_SRC)

## base-src-del
base-src-del:
	# just the .src
	@echo "Deleting .src ..."
	rm -rf $(BASE_CWD_SRC)

## base-src-dep ( will clone OR pull )
base-src-dep:
    # Smart like Dep, in that it first checks if its there.

    # Check is src code is there, to decide is i need to do a full clone.
    # TODO: Once this works, add a base-src-force, cause we might want to overide when we change VERSION, etc
	# Later can also check BASE_SRC_VERSION_WHICH to be smart and only flip if version changed in config.
	# I checked and it works correctly.
	@echo ""
	@echo "Checking for src existance ... "
	@echo ""

ifeq ($(BASE_SRC_VERSION_WHICH), )
	@echo "Not found ... "
	# $(MAKE) base-src
else
	@echo "Found ... "
	# $(MAKE) base-src-pull
endif

base-src: base-src-del base-src-init
	@echo "So, doing git clone ..."
	@echo ""
	# --single-branch, so dont get github pages, etc. https://github.com/pojntfx/hydrapp hit me.
	# I should look at pushing to Github Pages, basedon his code, because he go so much into it.
	cd $(BASE_CWD_SRC) && $(BASE_DEP_BIN_GIT_NAME) clone $(BASE_SRC_URL) -b $(BASE_SRC_VERSION) --single-branch
	cd $(BASE_CWD_SRC) && echo $(BASE_SRC_NAME) >> .gitignore


base-src-private: base-src-del base-src-init
	@echo ""
	@echo "doing clone with ssh credentials."
	@echo "these signatures are in the base.env"
	@echo ""
	
	cd $(BASE_CWD_SRC) && $(BASE_DEP_BIN_GIT_NAME) clone $(BASE_SRC_URL) -b $(BASE_SRC_VERSION) --single-branch
	cd $(BASE_CWD_SRC) && echo $(BASE_SRC_NAME) >> .gitignore


## base-src-status
base-src-status:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) status

## base-src-version-tag
base-src-version-tag:
	# CANT work out how to just return the result !!
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) describe --abbrev=4 --dirty --always --tags

## base-src-pull
base-src-pull:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) pull --rebase

SRC_PUSH_MESSAGE:=

## base-src-push
base-src-push: base-src-sign-set
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) add .
ifeq ($(SRC_PUSH_MESSAGE), )
	@echo ""
	@echo ""
	@echo "you need a commit message. Like this:"
	@echo "make SRC_PUSH_MESSAGE='im a duffus' base-src-push"
	@echo ""
	@echo ""
else
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) commit -S -am '$(SRC_PUSH_MESSAGE)'
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) push
endif
	

### signing of src code	

BASE_SRC_SIGNING_CONFIG=$(HOME)/.gitconfig
BASE_SRC_SIGNING_CONFIG_LOCAL=$(BASE_SRC)/.git/config

#BASE_SRC_SIGNING_USER_NAME=gedw99
#BASE_SRC_SIGNING_USER_EMAIL=gedw99@gmail.com

# Thre is ssh and smime options. Pick one

# ssh based
# https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key#telling-git-about-your-ssh-key

BASE_SRC_SIGNING_KEY_PRIV=$(HOME)/.ssh/gedw99_github.com
BASE_SRC_SIGNING_KEY=$(HOME)/.ssh/gedw99_github.com.pub
BASE_SRC_SIGNING_PROGRAM=ssh
BASE_SRC_SIGNING_FORMAT=ssh

# PKI based using smime
# https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key#telling-git-about-your-x509-key-1
#BASE_SRC_SIGNING_PATH=""
#BASE_SRC_SIGNING_KEY=e9abcabcef55053ead1dda7c297c6efd6c146a86
#BASE_SRC_SIGNING_PROGRAM=smimesign
#BASE_SRC_SIGNING_FORMAT=x509

BASE_SRC_SIGNING_BIN_NAME=smimesign
ifeq ($(BASE_OS_NAME),windows)
	BASE_SRC_SIGNING_BIN_NAME=smimesign.exe
endif
BASE_SRC_SIGNING_BIN_WHICH=$(shell command -v $(BASE_SRC_SIGNING_BIN_NAME))


base-src-sign-print-scope:
	@echo ""
	@echo "- git config scope"
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --list --show-scope --show-origin
	@echo ""

base-src-sign-dep-del:
	rm -f $(BASE_SRC_SIGNING_BIN_WHICH)
	#brew uninstall smimesign

base-src-sign-dep:
	@echo ""
	@echo " $(BASE_SRC_SIGNING_BIN_NAME) dep check ... "


ifeq ($(BASE_SRC_SIGNING_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_SRC_SIGNING_BIN_NAME) dep check: failed"
	$(MAKE) base-src-sign-dep-single
else
	@echo ""
	@echo "$(BASE_SRC_SIGNING_BIN_NAME) dep check: passed"
endif



base-src-sign-dep-single:
	# https://github.com/github/smimesign
	# Smimesign is an S/MIME signing utility for macOS and Windows that is compatible with Git. 
	# This allows developers to sign their Git commits and tags using X.509 certificates issued by public certificate authorities or their organization's internal certificate authority. 
	# Smimesign uses keys and certificates already stored in the macOS Keychain or the Windows Certificate Store.
	go install github.com/github/smimesign@latest
	#brew install smimesign


base-src-sign-run-h: base-src-sign-dep
	$(BASE_SRC_SIGNING_BIN_NAME) -h
	$(BASE_SRC_SIGNING_BIN_NAME) -v
base-src-sign-run-list: base-src-sign-dep
	$(BASE_SRC_SIGNING_BIN_NAME) --list-keys
base-src-sign-run-sign: base-src-sign-dep
	#$(BASE_SRC_SIGNING_BIN_NAME) --sign -h
	$(BASE_SRC_SIGNING_BIN_NAME) --sign USER-ID=$(BASE_SRC_SIGNING_USER_NAME)
	

base-src-sign-set:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.name $(BASE_SRC_SIGNING_USER_NAME)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.email  $(BASE_SRC_SIGNING_USER_EMAIL)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.signingkey $(BASE_SRC_SIGNING_KEY)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.program $(BASE_SRC_SIGNING_PROGRAM)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.format $(BASE_SRC_SIGNING_FORMAT)
base-src-sign-get:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --get user.name
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --get user.email
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.signingkey
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.program
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.format
base-src-sign-del:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.name ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.email ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.signingkey ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.program ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.format ""


### mod

base-mod-tidy:
	cd $(BASE_CWD_SRC)/$(BASE_SRC_NAME) && $(BASE_DEP_BIN_GO_NAME) mod tidy
base-mod-upgrade:
	$(BASE_DEP_BIN_GO_NAME) install github.com/oligot/go-mod-upgrade@latest
	cd $(BASE_CWD_SRC)/$(BASE_SRC_NAME) && go-mod-upgrade

### bin 



BASE_BIN_NAME=$(BASE_SRC_NAME)

# name with BASE_OS_ARCH
BASE_BIN_TARGET=$(BASE_BIN_NAME)_bin_$(BASE_OS_NAME)_$(BASE_OS_ARCH)
ifeq ($(BASE_OS_NAME),windows)
	BASE_BIN_TARGET=$(BASE_BIN_NAME)_bin_$(BASE_OS_NAME)_$(BASE_OS_ARCH).exe
endif


# suffix to the go.mod of the src, used by the go.work file
BASE_BIN_MOD := .

# suffix to the main.go of the src
BASE_BIN_ENTRY := .

# BASE_SRC_VERSION_WHICH gives back the version, so i can then add it to the LDFLags.
# E:G BASE_BIN_GO_BUILD_CMD:=CGO_ENABLED=0 go build -a -ldflags "-w -X '$(THIS_PKG)/cmd.Version=$(BASE_SRC_VERSION_WHICH)'" 

# name of pkg. useful for versioning.
# https://tip.golang.org/doc/go1.24
# The go build command now sets the main module’s version in the compiled binary based on the version control system tag and/or commit. A +dirty suffix will be appended if there are uncommitted changes. Use the -buildvcs=false flag to omit version control information from the binary.
# This will make versioing automatic i think...
# # name of pkg. useful for versioning.

BASE_BIN_PKG := .


# used for naming all binaries suffix.
BASE_BIN_SUFFIX_NATIVE=bin_$(BASE_OS_NAME)_$(BASE_OS_ARCH)
ifeq ($(BASE_OS_NAME),windows)
	BASE_BIN_SUFFIX_NATIVE=bin_$(BASE_OS_NAME)_$(BASE_OS_ARCH).exe
endif

# constants for bin targets
# ALL NAMING is based on the GOOS and GOARCH on purpose.
BASE_BIN_SUFFIX_DARWIN_AMD64=bin_darwin_amd64
BASE_BIN_SUFFIX_DARWIN_ARM64=bin_darwin_arm64
BASE_BIN_SUFFIX_LINUX_AMD64=bin_linux_amd64
BASE_BIN_SUFFIX_LINUX_ARM64=bin_linux_arm64
BASE_BIN_SUFFIX_WINDOWS_AMD64=bin_windows_amd64.exe
BASE_BIN_SUFFIX_WINDOWS_ARM64=bin_windows_arm64.exe
# wasm and wasi
BASE_BIN_SUFFIX_WASM_GO=bin_js_wasmgo.wasm
BASE_BIN_SUFFIX_WASI_GO=bin_wasip1_wasmgo.wasi
BASE_BIN_SUFFIX_WASM_TINY=bin_js_wasmtiny.wasm
BASE_BIN_SUFFIX_WASI_TINY=bin_wasip1_wasmtiny.wasi
# bundles
BASE_BIN_SUFFIX_EXT_ANDROID_AMD64=android_amd64.apk
BASE_BIN_SUFFIX_EXT_ANDROID_ARM64=android_arm64.apk
BASE_BIN_SUFFIX_EXT_DARWIN_AMD64=darwin_amd64.app
BASE_BIN_SUFFIX_EXT_DARWIN_ARM64=darwin_arm64.app
BASE_BIN_SUFFIX_EXT_IOS_AMD64=ios_amd64.app
BASE_BIN_SUFFIX_EXT_IOS_ARM64=ios_arm64.ipa
BASE_BIN_SUFFIX_EXT_WINDOWS_AMD64=windows_amd64.msi
BASE_BIN_SUFFIX_EXT_WINDOWS_ARM64=windows_arm64.msi



BASE_BIN_GO_INSTALL_CMD=$(BASE_DEP_BIN_GO_NAME) install -tags osusergo,netgo -ldflags '-extldflags "-static"'
# commented off for now. We can get more specific later.
#BASE_BIN_GO_BUILD_CMD=$(BASE_DEP_BIN_GO_NAME) build -tags osusergo,netgo -ldflags '-extldflags "-static"'
BASE_BIN_GO_BUILD_CMD=$(BASE_DEP_BIN_GO_NAME) build -buildvcs=false
BASE_BIN_GO_GARBLE_BUILD_CMD=$(BASE_AMP_GARBLE_BIN_NAME) build
BASE_BIN_GO_WASM_CMD=$(BASE_DEP_BIN_GO_NAME) build


# run time..
#BASE_BIN_MOD_FSPATH=$(BASE_SRC_NAME)/$(BASE_BIN_MOD)
BASE_BIN_MOD_FSPATH=$(BASE_CWD_SRC)/$(BASE_SRC_NAME)/$(BASE_BIN_MOD)
BASE_BIN_ENTRY_FSPATH=$(BASE_CWD_SRC)/$(BASE_SRC_NAME)/$(BASE_BIN_ENTRY)

# meta to help with binary updates.
# File Suffix, to make it easy to stay orgnaised.
# ! It will change to JSON later maybe
BASE_BIN_META_FILE_PREFIX=base_bin
BASE_BIN_META_FILE_SUFFIX=meta
BASE_BIN_META_FILE_EXT=.meta
# version of binaries (one)
# base_bin_version.meta
BASE_BIN_META_VERSION_FILENAME=base_bin_version.meta
# list of binaries (many)
# base_bin_list.meta
BASE_BIN_META_LIST_FILENAME=base_bin_list.meta

# boot of darwin binaries (one)
# base_bin_boot_darwin.sh
BASE_BIN_META_BOOT_DARWIN_FILENAME=base_bin_boot_darwin.sh
# boot of linxu binaries (one)
# base_bin_boot_linux.sh
BASE_BIN_META_BOOT_LINUX_FILENAME=base_bin_boot_linux.sh
# boot of windows binaries (one)
# base_bin_boot_windows.bat
BASE_BIN_META_BOOT_WINDOWS_FILENAME=base_bin_boot_windows.bat
# boot of wasm binaries (one)
# base_bin_boot_wasm
BASE_BIN_META_BOOT_WASM_FILENAME=base_bin_boot_wasm.sh
# boot of wasm binaries (one)
# base_bin_boot_wasm
BASE_BIN_META_BOOT_WASI_FILENAME=base_bin_boot_wasi.sh


BASE_BIN_META_WASM_EXE_NAME=wasm_exec.js
# $(GOROOT)/misc/wasm/wasm_exec.js
# $(GOROOT)/lib/wasm_exec.js
BASE_BIN_META_WASM_EXE_WHICH=$(GOROOT)/lib/wasm/wasm_exec.js

## base-bin-print
base-bin-print:
	@echo ""
	@echo "--- bin inputs ---"
	@echo "BASE_BIN_MOD:            $(BASE_BIN_MOD)"
	@echo ""
	@echo "BASE_BIN_ENTRY:          $(BASE_BIN_ENTRY)"
	@echo ""
	@echo "BASE_BIN_NAME:           $(BASE_BIN_NAME)"
	@echo ""
	@echo "BASE_BIN_PKG:            $(BASE_BIN_PKG)"
	@echo ""

	@echo ""
	@echo ""
	@echo "--- bin calculated ---"
	@echo ""
	@echo "BASE_CWD_SRC:            $(BASE_CWD_SRC)"
	@echo "BASE_SRC_CMD:            $(BASE_SRC_CMD)"
	@echo ""
	@echo "BASE_CWD_BIN_NAME:       $(BASE_CWD_BIN_NAME)"
	@echo "BASE_CWD_BIN:            $(BASE_CWD_BIN)"
	@echo ""
	@echo "BASE_BIN_MOD_FSPATH:     $(BASE_BIN_MOD_FSPATH)"
	@echo "BASE_BIN_ENTRY_FSPATH:   $(BASE_BIN_ENTRY_FSPATH)"
	@echo ""
	@echo "BASE_BIN_SUFFIX_NATIVE:  $(BASE_BIN_SUFFIX_NATIVE)"
	@echo "BASE_BIN_TARGET:         $(BASE_BIN_TARGET)"
	@echo ""
	@echo "BASE_BIN_GO_INSTALL_CMD:        $(BASE_BIN_GO_INSTALL_CMD)"
	@echo "BASE_BIN_GO_BUILD_CMD:          $(BASE_BIN_GO_BUILD_CMD)"
	@echo "BASE_BIN_GO_GARBLE_BUILD_CMD:   $(BASE_BIN_GO_GARBLE_BUILD_CMD)"
	@echo "BASE_BIN_GO_WASM_CMD:           $(BASE_BIN_GO_WASM_CMD)"
	@echo ""


## base-src-ls-h
base-bin-ls-h:
	$(BASE_AMP_TREE_BIN_NAME) --help
	# -d = only directories.
	# -h = directories and files
	# -C = colour
	  # PURPLE = folder
	  # GREEN = executable ( .sh )
## base-src-ls
base-bin-ls:
	# only folders, no files.
	# if not there does not blow up at all.
	$(BASE_AMP_TREE_BIN_NAME) -d -C $(BASE_CWD_BIN)
## base-src-ls-all
base-bin-ls-all:
	# all files and folders
	$(BASE_AMP_TREE_BIN_NAME) -h -C $(BASE_CWD_BIN)

## base-bin-ls-base ( just base, which is meta and base files as they all start with base. )
base-bin-ls-base:
	@echo ""
	$(BASE_AMP_TREE_BIN_NAME) -h -C -P base.* $(BASE_CWD_BIN) 

## base-bin-ls-native list the native files.
base-bin-ls-native:
	@echo ""
	@echo ""
	# darwin
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_DARWIN_AMD64) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_DARWIN_AMD64) $(BASE_CWD_BIN)
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_DARWIN_ARM64) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_DARWIN_ARM64) $(BASE_CWD_BIN)

	# linux
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_LINUX_AMD64) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_LINUX_AMD64) $(BASE_CWD_BIN)
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_LINUX_ARM64) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_LINUX_ARM64) $(BASE_CWD_BIN)

	# windows
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_WINDOWS_AMD64) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_WINDOWS_AMD64) $(BASE_CWD_BIN)
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_WINDOWS_ARM64) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_WINDOWS_ARM64) $(BASE_CWD_BIN)

## base-bin-list-wasm list the wasm files.
base-bin-ls-wasm:	
	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_WASM_GO) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_WASM_GO) $(BASE_CWD_BIN)

	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_WASI_GO) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_WASI_GO) $(BASE_CWD_BIN)

	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_WASM_TINY) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_WASM_TINY) $(BASE_CWD_BIN) 

	@echo ""
	@echo "-- list $(BASE_BIN_SUFFIX_WASI_TINY) --"
	$(BASE_AMP_TREE_BIN_NAME) -h -P $(BASE_BIN_SUFFIX_WASI_TINY) $(BASE_CWD_BIN) 




base-bin-init:
	@echo ""
	@echo "Creating .bin ..."
	@echo ""
	mkdir -p $(BASE_CWD_BIN)

	# del and create the meta version
	rm -f $(BASE_CWD_BIN)/$(BASE_BIN_META_VERSION_FILENAME)
	@echo $(BASE_SRC_VERSION_WHICH) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_VERSION_FILENAME)

	# delete meta list of binaries, so its always correct.
	# For CI with CGO builds, we will have to think of another way. Probably NATS as a global system at CI time,
	# with files going into NATS Object store.
	
	rm -f $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)

	# delete meta boot files so always correct.
	rm -f $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_DARWIN_FILENAME)
	rm -f $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_LINUX_FILENAME)
	rm -f $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WINDOWS_FILENAME)

	@echo ""
	@echo "checking if go project ..."

ifneq ($(wildcard $(BASE_BIN_MOD_FSPATH)/go.mod), )
	@echo ""
	@echo "go.mod found"
	@echo ""
	$(MAKE) base-bin-init-golang
else
	@echo ""
	@echo "NO go.mod found - Fix config please !!"
	@echo ""
endif


base-bin-init-golang:
	@echo ""
	@echo "Installing go using gobrew ..."
	@echo ""
	#cd $(BASE_CWD_SRC)/$(BASE_SRC_NAME) && $(BASE_DEP_BIN_GOBREW_NAME) use mod
	#cd $(BASE_CWD_SRC)/$(BASE_BIN_ENTRY_FSPATH) && $(BASE_DEP_BIN_GOBREW_NAME) use mod
	# fucker pulls over http to github to get JSON mapping.
	# can we stop this ?
	cd $(BASE_BIN_MOD_FSPATH) && $(BASE_DEP_BIN_GOBREW_NAME) -h
	#cd $(BASE_BIN_MOD_FSPATH) && $(BASE_DEP_BIN_GOBREW_NAME) use mod
	
	@echo ""
	@echo "- config go work file"
	@echo ""
	cd $(BASE_CWD_SRC) && touch go.work
	#cd $(BASE_CWD_SRC) && $(BASE_DEP_BIN_GO_NAME) work use $(BASE_SRC_NAME)/$(BASE_BIN_ENTRY)
	cd $(BASE_CWD_SRC) && $(BASE_DEP_BIN_GO_NAME) work use $(BASE_SRC_NAME)/$(BASE_BIN_MOD)

	@echo ""
	@echo "-bin mod and gen"
	@echo ""
	cd $(BASE_BIN_ENTRY_FSPATH) && $(BASE_DEP_BIN_GO_NAME) mod tidy
	cd $(BASE_BIN_ENTRY_FSPATH) && $(BASE_DEP_BIN_GO_NAME) generate -v .


## base-bin-del
base-bin-del:
	rm -rf $(BASE_CWD_BIN)



## base-bin builds your local bin
base-bin: base-bin-init
	@echo ""
	
	
ifeq ($(BASE_OS_NAME),darwin)
	$(MAKE) base-bin-darwin
endif
ifeq ($(BASE_OS_NAME),linux)
	$(MAKE) base-bin-linux
endif
ifeq ($(BASE_OS_NAME),windows)
	$(MAKE) base-bin-windows
endif
	
	# TODO echo out the binaries to a file, so we can pass it to the Artifacts thing like:
	# https://github.com/actions/upload-artifact/tree/v4/?tab=readme-ov-file#upload-an-individual-file
	# .meta file:
	# eg: pocketbase_meta holds only the version
	# so pocketbase_meta needs to hold "binary name, version name" in a CSV style system...

	# If we uplaod the meta, how the hell to merge it. Needs to pul it and then merge and then uplload.


	# then in CI, after the base-bin has pumped out the binary_meta to disk, the artifacts system
	# can use that to uplaod the binaries
	# then we can finally adapt all .mk files to get their binary from github releases that WE control using our naming convention.


	$(MAKE) base-bin-ls

## base-bin-all builds for Darin, Linux and Windows for amd64 and arm64
base-bin-all: base-bin-init base-bin-darwin base-bin-linux base-bin-windows

base-bin-darwin: #base-bin-init
	@echo ""
	@echo "-bin darwin"
	@echo ""

	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_AMD64) .
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_ARM64) .

	$(MAKE) base-bin-darwin-meta
	
base-bin-darwin-meta:
	# meta
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_AMD64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_ARM64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)

	# meta boot
	# Apple quarantine runner
	@echo xattr -dr com.apple.quarantine $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_AMD64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_DARWIN_FILENAME)
	@echo xattr -dr com.apple.quarantine $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_ARM64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_DARWIN_FILENAME)


base-bin-linux: #base-bin-init
	@echo ""
	@echo "-bin linux"
	@echo ""

	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) .
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) .

	
	#cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux-musl" $(BASE_BIN_GO_BUILD_CMD)  -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) .
	#cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux-musl" $(BASE_BIN_GO_BUILD_CMD)  -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) .

	$(MAKE) base-bin-linux-meta

base-bin-linux-meta:
	# meta
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)

	# meta boot
	@echo "hi linux" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_LINUX_FILENAME)
	@echo "hi linux" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_LINUX_FILENAME)

base-bin-windows: #base-bin-init
	@echo ""
	@echo "-bin windows"
	@echo ""

	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64) .
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64) .

	$(MAKE) base-bin-windows-meta
	
base-bin-windows-meta:
	# meta list
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)

	# meta boot
	@echo "hi windows" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WINDOWS_FILENAME)
	@echo "hi windows" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WINDOWS_FILENAME)


## base-bin-obf
base-bin-obf: base-bin-init
	# same as base-bin, but built using garble.
	@echo ""
	@echo "-bin with garble"
	@echo ""

ifeq ($(BASE_OS_NAME),darwin)
	$(MAKE) base-bin-darwin-obf
endif
ifeq ($(BASE_OS_NAME),linux)
	$(MAKE) base-bin-linux-obf
endif
ifeq ($(BASE_OS_NAME),windows)
	$(MAKE) base-bin-windows-obf
endif

	$(MAKE) base-bin-ls

base-bin-all-obf: base-bin-init base-bin-darwin-obf base-bin-linux-obf base-bin-windows-obf

base-bin-darwin-obf: #base-bin-init
	@echo ""
	@echo "-bin darwin obf"
	@echo ""
	
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_ARM64) *.go

	$(MAKE) base-bin-darwin-meta

base-bin-linux-obf: #base-bin-init
	@echo ""
	@echo "-bin linux obf"
	@echo ""
	
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) *.go

	$(MAKE) base-bin-linux-meta

base-bin-windows-obf: #base-bin-init
	@echo ""
	@echo "-bin windows obf"
	@echo ""
	
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64) *.go

	$(MAKE) base-bin-windows-meta

base-bin-wasm-obf:
	# garble and standard wasm build.
	@echo ""
	@echo "-bin wasm obf"
	@echo ""
	
	cd $(BASE_BIN_ENTRY_FSPATH) && GOOS=js GOARCH=wasm $(BASE_BIN_GO_GARBLE_BUILD_CMD) build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO)
	@echo "--- golang wasm wasm_exec.js ---"
	cp $(BASE_BIN_META_WASM_EXE_WHICH) $(BASE_CWD_BIN)

	

## base-bin-cgo
base-bin-cgo: base-bin-init
	# for gio stuff
	@echo ""
	@echo "-bin with cgo"
	@echo ""

ifeq ($(BASE_OS_NAME),darwin)
	$(MAKE) base-bin-darwin-cgo
endif
ifeq ($(BASE_OS_NAME),linux)
	$(MAKE) base-bin-linux-cgo
endif
ifeq ($(BASE_OS_NAME),windows)
	$(MAKE) base-bin-windows-cgo
endif
	
	$(MAKE) base-bin-ls

## base-bin-all-cgo build all GOOS, GOARCH with CGO. Needed for GIO and SQLite, etc
base-bin-all-cgo: base-bin-init base-bin-darwin-cgo base-bin-linux-cgo base-bin-windows-cgo
	# for gio stuff ...

base-bin-darwin-cgo:
	@echo ""
	@echo "-bin darwin cgo"
	@echo ""
	
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_ARM64) *.go

	$(MAKE) base-bin-darwin-meta

base-bin-linux-cgo:
	@echo ""
	@echo "-bin dafwin cgo"
	@echo ""
	
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) *.go

	$(MAKE) base-bin-linux-meta

base-bin-windows-cgo:
	@echo ""
	@echo "-bin windows cgo"
	@echo ""
	cd $(BASE_BIN_ENTRY_FSPATH) && $(BASE_DEP_BIN_GO_NAME) mod tidy
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=windows GOARCH=arm64 $(BASE_BIN_GO_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64) *.go

	$(MAKE) base-bin-windows-meta

## base-bin-cgo-obf
base-bin-cgo-obf: base-bin-init 
	# cgo and garbled
	# for gio stuff ...
	@echo ""
	@echo "-bin with cgo and garble"
	@echo ""

ifeq ($(BASE_OS_NAME),darwin)
	$(MAKE) base-bin-darwin-cgo-obf
endif
ifeq ($(BASE_OS_NAME),linux)
	$(MAKE) base-bin-linux-cgo-obf
endif
ifeq ($(BASE_OS_NAME),windows)
	$(MAKE) base-bin-windows-cgo-obf
endif

	$(MAKE) base-bin-ls

## base-bin-all-cgo-obf
base-bin-all-cgo-obf: base-bin-init base-bin-darwin-cgo-obf base-bin-linux-cgo-obf base-bin-windows-cgo-obf

base-bin-darwin-cgo-obf:
	@echo ""
	@echo "-bin darwin cgo obf"
	@echo ""
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_DARWIN_ARM64) *.go

	$(MAKE) base-bin-darwin-meta

base-bin-linux-cgo-obf:
	@echo ""
	@echo "-bin linux cgo obf"
	@echo ""
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_LINUX_ARM64) *.go

	$(MAKE) base-bin-linux-meta

base-bin-windows-cgo-obf:
	@echo ""
	@echo "-bin windows cgo obf"
	@echo ""
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64) *.go
	cd $(BASE_BIN_ENTRY_FSPATH) && CGO_ENABLED=1 GOOS=windows GOARCH=arm64 $(BASE_BIN_GO_GARBLE_BUILD_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64) *.go

	$(MAKE) base-bin-windows-meta
	

base-bin-inspect-size: base-bin-inspect-size-text

base-bin-inspect-size-start:
	@echo ""
	@echo "Using go-size-analyzer to see how well its obfuscates ..."
	@echo ""
	@echo "BASE_RUN: $(BASE_RUN)"
	@echo ""
	# https://github.com/Zxilly/go-size-analyzer
	go install github.com/Zxilly/go-size-analyzer/cmd/gsa@latest
	#gsa $(BASE_RUN)
	#gsa -h
	#gsa --version

base-bin-inspect-size-text: base-bin-inspect-size-start
	# simpel etxt output
	gsa $(BASE_RUN)
base-bin-inspect-size-web: base-bin-inspect-size-start
	# works and holds open Server.
	# http://localhost:8090
	gsa --web --listen=:8090 $(BASE_RUN)

base-bin-inspect-size-svg: base-bin-inspect-size-start
	# svg
	# gsa cockroach-darwin-amd64 -f svg -o data.svg --hide-sections
	gsa $(BASE_RUN) -f svg -o $(BASE_RUN).svg --hide-sections
	@echo ""
	@echo "Here is your SVG:"
	@echo "$(BASE_RUN).svg"
	@echo ""
base-bin-inspect-size-tui: base-bin-inspect-size-start
	# tui, and holds open
	# gsa --tui golang-compiled-binary
	gsa --tui $(BASE_RUN)


base-test-src:
    # runs stardard golang tests on src.
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_GO_NAME) test -v ./...
	#$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_GO_NAME) test -v ./... --update

base-test-golden:
	# TODO: Get golang tests outputting to JSON....
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_GO_NAME) tool test2json

base-bin-pub:
	# This copies the binaries to the Root git repo.
	# Its kind is useful locally...
	@echo ""
	@echo "publishing bin to:"
	@echo "BASE_GITROOT_BIN:   $(BASE_GITROOT_BIN)"
	@echo ""
	mkdir -p $(BASE_GITROOT_BIN)
	cp -r $(BASE_CWD_BIN)/* $(BASE_GITROOT_BIN)

	@echo ""
	@echo "publishing .mk's to:"
	@echo "BASE_GITROOT_BIN:   $(BASE_GITROOT_BIN)"
	cp $(BASE_MAKEFILE) $(BASE_GITROOT_BIN)/base.mk
	cp ./$(BASE_ARTIFACTS) $(BASE_GITROOT_BIN)/

	@echo ""
	@echo "listing BASE_GITROOT_BIN"
	ls -al $(BASE_GITROOT_BIN)
	@echo ""

## base-bin-pub-nats

# The idea is to pub to NATS Object Store and nats pushes to all my Servers. 
# Each Server just needs the same NATS CLI polling for Bucket updates.
# 1 bucket in 1 Project in this repo, with the exact same Makefile. 

# NATS 

# Hetzner and fly, and maybe CF later.

# Synadia Cloud: can use their Cloud to manage many of my nats system
# Synadia Plaform: is for actual storage. TOO EXPENSIVE. LETS just run NATS on Fly
# Hetzner to run the binaries. 

BASE_DEP_BIN_NATS_CLI_NAME=nats
ifeq ($(BASE_OS_NAME),windows)
	BASE_DEP_BIN_NATS_CLI_NAME=nats.exe
endif
BASE_DEP_BIN_NATS_CLI_WHICH=$(shell command -v $(BASE_DEP_BIN_NATS_CLI_NAME))

BASE_DEP_BIN_NATS_CLI_RUN_FILE=Makefile
# Synadia
BASE_DEP_BIN_NATS_CLI_RUN_CONSOLE=https://cloud.synadia.com/teams/2XrIt5ApHyjVq8XkELhTaP4vfO3
BASE_DEP_BIN_NATS_CLI_RUN_TOKEN=uat_9ZcVS8aBRtanIGsgz7k0cfFnbIZhEcZjvxcduqWhKHZcDzP7igA3QeBm0Y23q97j
# Hetzner
BASE_DEP_BIN_NATS_CLI_RUN_HCLOUD_TOKEN=?
BASE_DEP_BIN_NATS_CLI_RUN_HETZNER=https://console.hetzner.cloud/projects/2331110/servers/54146093/overview

base-bin-pub-nats-print:
	@echo ""

	@echo "BASE_DEP_BIN_NATS_CLI_NAME:    $(BASE_DEP_BIN_NATS_CLI_NAME)"
	@echo "BASE_DEP_BIN_NATS_CLI_WHICH:   $(BASE_DEP_BIN_NATS_CLI_WHICH)"
	@echo ""
	@echo "BASE_DEP_BIN_NATS_CLI_RUN_FILE:     $(BASE_DEP_BIN_NATS_CLI_RUN_FILE)"
	@echo "BASE_DEP_BIN_NATS_CLI_RUN_CONSOLE:  $(BASE_DEP_BIN_NATS_CLI_RUN_CONSOLE)"
	@echo "BASE_DEP_BIN_NATS_CLI_RUN_TOKEN:    $(BASE_DEP_BIN_NATS_CLI_RUN_TOKEN)"
	@echo ""

base-bin-pub-nats-dep:
    # nats cli
    # https://github.com/nats-io/natscli
    # so we can upload and download binaires.
    # can change each .mk to check if its in nats, during dep check

ifeq ($(BASE_DEP_BIN_NATS_CLI_WHICH), )
	@echo ""
	@echo "$(BASE_DEP_BIN_NATS_CLI_NAME) dep check: failed"
	go install github.com/nats-io/natscli/nats@latest
else
	@echo ""
	@echo "$(BASE_DEP_BIN_NATS_CLI_NAME) dep check: passed"
endif

base-bin-pub-nats-run-upload: base-bin-pub-nats-dep
	@echo ""
	@echo "publishing bin of:" $(BASE_DEP_BIN_NATS_CLI_RUN_FILE)
	$(BASE_DEP_BIN_NATS_CLI_NAME) object put $(BASE_DEP_BIN_NATS_CLI_RUN_FILE)

base-bin-pub-nats-run-download: base-bin-pub-nats-dep
	@echo ""
	@echo "downloading bin to:" 
	$(BASE_DEP_BIN_NATS_CLI_NAME) object get $(BASE_DEP_BIN_NATS_CLI_RUN_FILE)

base-bin-web:
	# 1. Build using standard golang wasm
	$(BASE_BIN_GO_BUILD_CMD)
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && GOOS=js GOARCH=wasm $(BASE_BIN_GO_WASM_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO)

	# 2: Build using gio.
	# https://gioui.org/doc/install/wasm
	$(BASE_DEP_BIN_GO_NAME) install gioui.org/cmd/gogio@latest
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && gogio -target js -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web . 

	# 3. Copy normal wasm inside the GIO packging as main.wasm
	cp $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO) $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web/main.wasm

## base-bin-web-ofs
base-bin-web-ofs:
	# All works

	# 1. Build using garble 
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && GOOS=js GOARCH=wasm garble build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO)

	# 2: Build using gio.
	# https://gioui.org/doc/install/wasm
	$(BASE_DEP_BIN_GO_NAME) install gioui.org/cmd/gogio@latest
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && gogio -target js -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web . 

	# 3. Copy garbled wasm inside the GIO packging as main.wasm
	cp $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO) $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web/main.wasm
	
## base-bin-web-ofs
base-bin-web-tiny:
	# use tinygo ...
	# TO garble we need to copy first and then build... FUCK.

	# 1. Build using tinygo 
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_TINYGO_NAME) build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_js_tinywasm.wasm -target wasm .

	# 2: Build using gio.
	# https://gioui.org/doc/install/wasm
	$(BASE_DEP_BIN_GO_NAME) install gioui.org/cmd/gogio@latest
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && gogio -target js -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web . 

	# 3. Copy garbled wasm inside the GIO packging as main.wasm
	cp $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_js_tinywasm.wasm $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web/main.wasm
	
## base-bin-wasm
base-bin-wasm:
	# golang
	@echo ""
	@echo "--- golang wasm ---"
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && GOOS=js GOARCH=wasm $(BASE_BIN_GO_WASM_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO) -ldflags="-s -w" -gcflags="-trimpath=${PWD}" -trimpath 
	@echo ""
	@echo "--- golang wasm wasm_exec.js ---"
	cp $(BASE_BIN_META_WASM_EXE_WHICH) $(BASE_CWD_BIN)
	
	$(MAKE) base-bin-wasm-meta

## base-bin-wasi
base-bin-wasi:
	# Broke this out, because it rarely works...
	# I dont need it.
	@echo ""
	@echo "--- golang wasi ---"
	# USEFUL: https://go.dev/blog/wasi
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && GOOS=wasip1 GOARCH=wasm $(BASE_BIN_GO_WASM_CMD) -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_GO)
	
	$(MAKE) base-bin-wasm-meta

base-bin-wasm-meta:
	@echo ""
	@echo "--- golang wasm / wasi meta ---"

	# meta
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_GO) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)

	# meta boot
	@echo "hi go wasm booter" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_GO) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WASM_FILENAME)
	@echo "hi go wasi booter" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_GO) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WASI_FILENAME)

	
base-bin-wasm-tiny:
	@echo ""
	@echo "--- tinygo wasm ---"
	@echo ""
	@echo "Off until go 1.24 works with tinygo !!"
	@echo ""

	# FAILS: deckd.go:57:11: undefined: os.Chdir
	# NOTE: works with golang wasm and wasip1 compile.
	# I made an Issue here: https://github.com/tinygo-org/tinygo/issues/4675

	# cryptop rand bug: https://github.com/tinygo-org/tinygo/issues/4777#issuecomment-2704563283
	# wasi works, but not wasm.

	# debug	
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_TINYGO_NAME) build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_TINY) -target wasm .
	# release
	#$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_TINYGO_NAME) build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_TINY) -target wasm -no-debug .
	
	# tinygo wasi
	@echo ""
	@echo "--- tinygo wasi ---"
	@echo ""
	@echo "Off until go 1.24 works with tinygo !!"
	@echo ""

	# debug	
	$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_TINYGO_NAME) build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_TINY) -target wasi .
	# release
	#$(BASE_SRC_CMD) cd $(BASE_BIN_ENTRY) && $(BASE_DEP_BIN_TINYGO_NAME) build -o $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_TINY) -target wasi -no-debug .
	
	$(MAKE) base-bin-wasm-tiny-meta

base-bin-wasm-tiny-meta:
	@echo ""
	@echo "--- tinygo wasm / wasi meta  ---"

	# meta
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_TINY) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)
	@echo $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_TINY) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_LIST_FILENAME)

	# meta boot
	@echo "hi tinygo wasm booter" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASM_TINY) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WASM_FILENAME)
	@echo "hi tinygo wasi booter" $(BASE_BIN_NAME)_$(BASE_BIN_SUFFIX_WASI_TINY) >> $(BASE_CWD_BIN)/$(BASE_BIN_META_BOOT_WASI_FILENAME)




### bin-sign

BASE_DEP_BIN_QUILL_NAME=quill
ifeq ($(BASE_OS_NAME),windows)
	BASE_DEP_BIN_QUILL_NAME=quill.exe
endif
BASE_DEP_BIN_QUILL_WHICH=$(shell command -v $(BASE_DEP_BIN_QUILL_NAME))

# meta to run stuff
# for signing
BASE_QUILL_SIGN_P12=[path-to-p12]
BASE_QUILL_SIGN_PASSWORD=[p12-password]
# for notorisation
BASE_QUILL_NOTARY_KEY=[path-to-private-key-file-from-apple]   # can also be base64 encoded contents instead of a file path
BASE_QUILL_NOTARY_KEY_ID=[apple-private-key-id]               # e.g. XS319FABCD
BASE_QUILL_NOTARY_ISSUER=[apple-notary-issuer-id]             # e.g. a1234b5-1234-5f5d-b0c8-1234bedc5678

base-bin-sign-print:
	# just like base-src-sign-print, which is about signing the source code,
	# we sign the binaries of course

	# this is done BEFORE packing. Packaging is ZIP packing, not OS packaging.

	@echo ""
	@echo "BASE_DEP_BIN_QUILL_NAME:    $(BASE_DEP_BIN_QUILL_NAME)"
	@echo "BASE_DEP_BIN_QUILL_WHICH:   $(BASE_DEP_BIN_QUILL_WHICH)"
	@echo ""
	@echo ""
	@echo "- run vars"
	@echo "BASE_QUILL_SIGN_P12:        $(BASE_QUILL_SIGN_P12)"
	@echo "BASE_QUILL_SIGN_PASSWORD:   $(BASE_QUILL_SIGN_PASSWORD)"
	@echo ""
	@echo "BASE_QUILL_NOTARY_KEY:      $(BASE_QUILL_NOTARY_KEY)"
	@echo "BASE_QUILL_NOTARY_KEY_ID:   $(BASE_QUILL_NOTARY_KEY_ID)"
	@echo "BASE_QUILL_NOTARY_ISSUER:   $(BASE_QUILL_NOTARY_ISSUER)"
	@echo ""

base-bin-sign-dep-del:
	rm -f $(BASE_DEP_BIN_QUILL_WHICH)
base-bin-sign-dep:
    # https://github.com/anchore/quill
    # signs and notorised MAC stuff from ANY OS, and so works in Linux CI, etc 

    # we do it on the fly.

ifeq ($(BASE_DEP_BIN_QUILL_WHICH), )
	@echo ""
	@echo "$(BASE_DEP_BIN_QUILL_NAME) dep check: failed"
	go install github.com/anchore/quill/cmd/quill@latest
else
	@echo ""
	@echo "$(BASE_DEP_BIN_QUILL_NAME) dep check: passed"
endif
	
base-bin-sign-run-h: base-bin-sign-dep
	quill -h
	quill version
	#quill embedded-certificates
	#mkdir -p .bin

base-bin-sign-run-describe: base-bin-sign-dep
	quill describe $(BASE_DEP_BIN_QUILL_WHICH)

base-bin-sign-run-extract: base-bin-sign-dep
	quill extract $(BASE_DEP_BIN_QUILL_WHICH)

base-bin-sign-run-sign: base-bin-sign-dep
	# This is stage 1
	quill sign $(BASE_DEP_BIN_QUILL_WHICH)
base-bin-sign-run-notorise: base-bin-sign-dep
	# This is stage 2
	quill sign-and-notarize $(BASE_DEP_BIN_QUILL_WHICH)
	



### run

BASE_RUN=$(BASE_CWD_BIN)/$(BASE_BIN_TARGET)

BASE_RUN_NAME := ?
BASE_RUN_ENTRY := ?

## base-run-print:
base-run-print:
	@echo ""
	@echo "--- run ---"
	@echo "BASE_RUN:             $(BASE_RUN)"
	@echo "BASE_RUN_NAME:        $(BASE_RUN_NAME)"
	@echo "BASE_RUN_ENTRY:       $(BASE_RUN_ENTRY)"
	@echo ""
	@echo "- meta"
	@echo "BASE_BIN_META_VERSION_FILENAME:       $(BASE_BIN_META_VERSION_FILENAME)"
	@echo "BASE_BIN_META_LIST_FILENAME:          $(BASE_BIN_META_LIST_FILENAME)"
	@echo "- meta data"
	# github.com/benhoyt/goawk
	# also in benthos, so IF i include benthos i get it for free and can at runtime to inspections.
	# so benthos is part of base install here AND for any cloud.
	# https://github.com/redpanda-data/connect/blob/main/docs/modules/components/pages/processors/awk.adoc
	# $(shell awk '// { if ($$2 = 'APP_VERSION_MAJOR) { print $$3 } }' < Version.h)
	@echo "BASE_BIN_META_VERSION:                $(BASE_BIN_META_VERSION_FILENAME)"

## base-run:
base-run:
	$(BASE_BIN_TARGET)
## base-run-h:
base-run-h:
	$(BASE_BIN_TARGET) -h
base-run-args:

## base-run-web
base-run-web:
	#https://github.com/eliben/static-server
	#go run github.com/eliben/static-server@latest
	#static-server  $(BASE_CWD_BIN)

	# https://github.com/caddyserver/caddy/releases/tag/v2.8.4
	# https://github.com/caddyserver/caddy/releases/tag/v2.9.0-beta.3
	# has lots of bug fixes apparently.
	$(BASE_DEP_BIN_GO_NAME) install github.com/caddyserver/caddy/v2/cmd/caddy@v2.9.0-beta.3


	#cd $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_web && caddy file-server --browse

	# https://localhost
	# http://localhost will auto redirect to https :)

	#cd $(BASE_CWD_BIN) && caddy file-server --domain localhost --browse --templates
	# https://localhost
	caddy file-server --domain localhost --browse --reveal-symlinks --templates

	



## base-run-wasm:
base-run-wasm:
	@echo "TODO: Need wazero"
	## wasm
	Need wazero to run the wasm, and not the OS.

ifeq ($(BASE_OS_NAME),linux)
ifeq ($(BASE_OS_ARCH),amd64)
	@echo "linux_amd64"
	$(BASE_CWD_BIN)/$(BASE_SRC_NAME)_linux_amd64
endif
ifeq ($(BASE_OS_ARCH),arm64)
	@echo "linux_arm64"
endif
endif


### deplpy

# so we can share binaries. Maybe a better name is activate ?
# then deploy is when NATS and File Watcher scoops then
# will involve goreman and psutil and gops...

# TASK File global: https://taskfile.dev/usage/#running-a-global-taskfile

## base-pwd-print
base-pwd-print:
	@echo ""
	@echo "--- base : pwd ---"
	@echo "BASE_PWD:       $(BASE_PWD)"
	@echo ""
	@echo ""
## base-pwd-init
base-pwd-init:
	mkdir -p $(BASE_PWD)
	mkdir -p $(BASE_PWD_BIN)
	mkdir -p $(BASE_PWD_DATA)
	mkdir -p $(BASE_PWD_META)

## base-pwd-init-del
base-pwd-init-del:
	rm -rf $(BASE_PWD)

## base-pwd-open
base-pwd-open:
	open $(BASE_HOME)

base-pwd-list:
	ls -al $(BASE_HOME)

## base-pwd-bin-deploy
base-pwd-bin-deploy:
	# copies up to root of git repo
	cp -r $(BASE_CWD_BIN) $(BASE_PWD_BIN)



### pack ( packaging )

# This is hardcoded to use the BASE_CWD_DEP and BASE_CWD_PACK for a very good reason.
# SO that we can comppose a Proejct out of other proejcts. the .bin and .dep concept with NATS makign it work.
# DONT change this Gerard - its simple and keeps me out of writing tons of golang.
# ITS ALL ABOUT reusing open sorrce code without coding all the way up the chain.


BASE_CWD_PACK_REVERSE=$(BASE_CWD_PACK)-reverse

base-pack-print:
	@echo  ""
	@echo  "--- pack print"
	@echo "takes .dep and copies into .pack"
	@echo  ""
	@echo  "-- sources"
	@echo  "BASE_CWD_DEP:            $(BASE_CWD_DEP)"
	@echo  "BASE_CWD_BIN:            $(BASE_CWD_BIN)"
	@echo  ""
	@echo  "-- target"
	@echo  "BASE_CWD_PACK:           $(BASE_CWD_PACK)"
	@echo  ""
	@echo  "-- target-test"
	@echo  "BASE_CWD_PACK_REVERSE:   $(BASE_CWD_PACK_REVERSE)"
	
base-pack-darwin:
	# DONT DO THIS - it ust makes build hard and gives us nothign.
	# this takes the 2 darwins bins and turns it into 1 via lipo
	@echo  ""
	@echo  "--- pack amd and arm into 1: ":  $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_darwin
	@echo  ""

	@lipo -create -output $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_darwin $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_darwin_amd64 $(BASE_CWD_BIN)/$(BASE_BIN_NAME)_darwin_arm64

	# copied from: https://github.com/rwilgaard/alfred-github-search/blob/main/Makefile#L23
	# takes the many binareis and packs them. Can them put into a .app, egc
	#@lipo -create -output workflow/$(PROJECT_NAME) workflow/$(PROJECT_NAME)-amd64 workflow/$(PROJECT_NAME)-arm64
	#@rm -f workflow/$(PROJECT_NAME)-amd64 workflow/$(PROJECT_NAME)-arm64

base-pack-init: base-pack-init-del
	# So as not to destroy our dev environment, we package to the ".pack" folder
	mkdir -p $(BASE_CWD_PACK)
base-pack-init-del:
	rm -rf $(BASE_CWD_PACK)



base-pack: base-pack-init
	# This is a general Packer for ANYTHING, so that everything is forward engineering.

	# Why ?
	# To make it easy to deploy, and never miss a file, we zip everything in each folder.
	# We can then inject that zip a Docker, or NATS Stream, or S3. 
	# The point is it an immutable copy !!
	# idea from: https://github.com/mynaparrot/plugNmeet-server/blob/main/Makefile#L18

	# SOURCE and TARGET is all it needs:
	# $(BASE_CWD_PACK) is the target 
	# $(BASE_CWD_BIN) and $(BASE_CWD_DEP) are currently used.

	
	# deps
	@echo ""
	@echo "Contents of .dep"
	ls -al $(BASE_CWD_DEP)
	@echo ""

	# compress .dep folder
	$(BASE_AMP_ARC_COM_BIN_NAME) -a -f $(BASE_CWD_PACK)/dep.zip $(BASE_CWD_DEP)

	# bins
	@echo ""
	@echo "Contents of .bin"
	#ls -al $(BASE_CWD_BIN)
	@echo ""

	# compress .bin folder
	#$(BASE_AMP_ARC_COM_BIN_NAME) -a -f $(BASE_CWD_PACK)/bin.zip $(BASE_CWD_BIN)/*
	
	# Lastly the Makefile and base.mk help.mk, so that it can unpack and run itself
	cp Makefile $(BASE_CWD_PACK)
	cp $(BASE_MAKE_IMPORT)/base*.* $(BASE_CWD_PACK)
	cp $(BASE_MAKE_IMPORT)/help*.* $(BASE_CWD_PACK)

base-pack-reverse-init: base-pack-reverse-init-del
	mkdir -p $(BASE_CWD_PACK_REVERSE)
base-pack-reverse-init-del:
	# We dont do this normally
	rm -rf $(BASE_CWD_PACK_REVERSE)
base-pack-reverse: base-pack-reverse-init
	# unzip the same things.

	# The Server or Client will run this eventually.

	# We want it to put the files in the same folder that the zips files live,
	# so that no paths are needed at runtime. The Folder is our single dependency.

	# Where to put it?
	@echo ""
	@echo "- unpack "
	@echo "to $(BASE_CWD_PACK_REVERSE)"
	@echo ""

	# decompress dep.zip
	$(BASE_AMP_ARC_COM_BIN_NAME) -x -f $(BASE_CWD_PACK)/dep.zip $(BASE_CWD_PACK_REVERSE)

	# Run it ( just for now )
	#cd $(BASE_CWD_PACK) && $(MAKE) 


base-pack-gui:
	# TODO: GUI needs different shit
base-pack-service:
	# Servcies need different shit.


### inject


BASE_INJECT_TARGET=$(BASE_SRC)/.dep

## Injects whats in .dep into the src repo, so it can self run.
base-inject-print:
	@echo ""
	@echo "BASE_CWD_DEP:          $(BASE_CWD_DEP)"
	@echo "BASE_INJECT_TARGET:    $(BASE_INJECT_TARGET)"
	@echo ""
base-inject:
	# target repo...
	rm -rf $(BASE_INJECT_TARGET)
	mkdir -p $(BASE_INJECT_TARGET)

	# copy .dep into target .dep
	cp $(BASE_CWD_DEP)/*.* $(BASE_INJECT_TARGET)/

	# then delete any binaries
	rm -f $(BASE_INJECT_TARGET)/*_bin_*

### data

base-data-print:
	@echo ""
	@echo "BASE_CWD_DATA:   $(BASE_CWD_DATA)"
	@echo ""

base-data-init:
	# Only their data folder.
	mkdir -p $(BASE_CWD_DATA)
base-data-del:
	# Only their data folder.
	rm -rf $(BASE_CWD_DATA)


	
	


