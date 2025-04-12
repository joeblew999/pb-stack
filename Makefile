# bootstrapper 
# on Desktop, gets task and golang on darwin, linux and windows.
# on CI, the actions do this.

# ISSUE: go installing task branches, DONT change the version of it.
# I use this in base.taskfile.yml, in order to upgrade task on the fly.

include ./mod/sops/Makefile

BASE_OS_NAME := $(shell go env GOOS)



# https://github.com/go-task/task
TASK_BIN_NAME=task
ifeq ($(BASE_OS_NAME),windows)
	TASK_BIN_NAME=git.exe
endif
# https://github.com/go-task/task/releases/tag/v3.42.1
TASK_BIN_VERSION=v3.42.1
# https://github.com/go-task/task/tree/recursive-config-search
TASK_BIN_VERSION=recursive-config-search
TASK_BIN_WHICH=$(shell command -v $(TASK_BIN_NAME))
TASK_BIN_WHICH_VERSION=$(shell $(TASK_BIN_NAME) --version)


# https://github.com/getsops/sops
SOPS_BIN_NAME=sops
ifeq ($(BASE_OS_NAME),windows)
	SOPS_BIN_NAME=sops.exe
endif
SOPS_BIN_VERSION=latest
# https://github.com/getsops/sops/releases/tag/v3.10.1
#SOPS_BIN_VERSION=v3.10.1
SOPS_BIN_WHICH=$(shell command -v $(SOPS_BIN_NAME))
SOPS_BIN_WHICH_VERSION=$(shell $(SOPS_BIN_NAME) --disable-version-check --version)

print:
	@echo ""
	@echo "TASK_BIN_NAME:            $(TASK_BIN_NAME)"
	@echo "TASK_BIN_VERSION:         $(TASK_BIN_VERSION)"
	@echo "TASK_BIN_WHICH:           $(TASK_BIN_WHICH)"
	@echo "TASK_BIN_WHICH_VERSION:   $(TASK_BIN_WHICH_VERSION)"
	@echo ""
	$(MAKE) sops-print
	@echo ""

task-del:
	rm -rf $(TASK_BIN_WHICH)
task:
ifeq ($(TASK_BIN_WHICH), )
	@echo ""
	@echo "$(TASK_BIN_NAME) dep check: failed, so installing ..."
	go install github.com/go-task/task/v3/cmd/task@$(TASK_BIN_VERSION)
	@echo ""
else
	@echo ""
	@echo "$(TASK_BIN_NAME) dep check: passed"
	@echo ""
endif
	
