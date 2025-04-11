# bootstrapper 
# on Desktop, gets task and golang on darwin, linux and windows.
# on CI, the actions do this.

# ISSUE: go installing task branches, DONT change the version of it.
# I use this in base.taskfile.yml, in order to upgrade task on the fly.


TASK_BIN_NAME=task
ifeq ($(BASE_OS_NAME),windows)
	TASK_BIN_NAME=git.exe
endif
# https://github.com/go-task/task/releases/tag/v3.42.0
# https://github.com/go-task/task/releases/tag/v3.42.1
TASK_BIN_VERSION=v3.42.1
# https://github.com/go-task/task/tree/recursive-config-search
TASK_BIN_VERSION=recursive-config-search
TASK_BIN_WHICH=$(shell command -v $(TASK_BIN_NAME))
TASK_BIN_WHICH_VERSION=$(shell $(TASK_BIN_NAME) --version)
TASK_BIN_WHICH_VERSION_SELF=$(shell $(TASK_BIN_NAME) --version)

print:
	@echo ""
	@echo "TASK_BIN_NAME:            $(TASK_BIN_NAME)"
	@echo "TASK_BIN_VERSION:         $(TASK_BIN_VERSION)"
	@echo "TASK_BIN_WHICH:           $(TASK_BIN_WHICH)"
	@echo "TASK_BIN_WHICH_VERSION:   $(TASK_BIN_WHICH_VERSION)"
	
	@echo ""

task-del:
	rm -rf $(TASK_BIN_WHICH)
task:
	go install github.com/go-task/task/v3/cmd/task@$(TASK_BIN_VERSION)
	

