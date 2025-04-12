# bootstrapper 
# on Desktop, gets task and golang on darwin, linux and windows.
# on CI, the actions do this.

# ISSUE: go installing task branches, DONT change the version of it.
# I use this in base.taskfile.yml, in order to upgrade task on the fly.

.DEFAULT_GOAL := print
.PHONY: default

include ./all.env
include ./mod/sops/Makefile

# assumes golang is installed. Will do for now.
BASE_OS_NAME := $(shell go env GOOS)

# https://github.com/go-task/task
TASK_BIN_NAME=task
ifeq ($(BASE_OS_NAME),windows)
	TASK_BIN_NAME=task.exe
endif
# From all.env
TASK_BIN_VERSION=$(BASE_TASK_VERSION_ENV)
TASK_BIN_WHICH=$(shell command -v $(TASK_BIN_NAME))
TASK_BIN_WHICH_VERSION=$(shell $(TASK_BIN_NAME) --version)


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
	# force it as we want the version in the all.env
	go install github.com/go-task/task/v3/cmd/task@$(TASK_BIN_VERSION)

	
