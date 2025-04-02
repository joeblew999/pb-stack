# https://github.com/titpetric/task-ui
# https://github.com/yonasBSD/task-ui seems to be ahead by 4 months..

BASE_SRC_NAME=task-ui
BASE_SRC_URL=https://github.com/titpetric/task-ui
# no tags
BASE_SRC_VERSION=main

BASE_BIN_NAME=task_web
BASE_BIN_MOD=.
BASE_BIN_ENTRY=.

#EX=$(PWD)/../ex

task-web-print: task-print
	@echo ""
	@echo "BASE_SRC_NAME:     $(BASE_SRC_NAME)"
	@echo "BASE_BIN_NAME:     $(BASE_BIN_NAME)"
	@echo "BASE_RUN:          $(BASE_RUN)"
	@echo ""

task-web-all: task-web-dep task-web-src task-web-bin

task-web-dep: base-dep

task-web-src: base-src 

task-web-bin: base-bin-init base-bin-init-golang base-bin
task-web-bin-all: task-web-bin base-bin-all
task-web-bin-obf: task-web-bin base-bin-obf

### run

TASK_UI_RUN_PATH=$(BASE_SRC)/examples
TASK_UI_DOCKER_PATH=$(BASE_SRC)/docker

TASK_UI_RUN_VAR_NAME=docker-compose
TASK_UI_RUN_VAR_NAME=doctl
TASK_UI_RUN_VAR_NAME=git-pull
TASK_UI_RUN_VAR_NAME=go
#TASK_UI_RUN_VAR_NAME=rclone-dropbox

TASK_UI_RUN_PATH_WHICH=$(TASK_UI_RUN_PATH)/$(TASK_UI_RUN_VAR_NAME)

# --history-enable
# requires TTY, so skip it.
#TASK_UI_RUN_CMD=cd $(TASK_UI_RUN_PATH_WHICH) && $(BASE_BIN_TARGET) --history-enable
TASK_UI_RUN_CMD=cd $(TASK_UI_RUN_PATH_WHICH) && $(BASE_BIN_TARGET)

task-web-run-print:
	@echo ""
	@echo "BASE_BIN_TARGET:           $(BASE_BIN_TARGET)"
	@echo ""
	@echo "TASK_UI_RUN_PATH:          $(TASK_UI_RUN_PATH)"
	@echo "TASK_UI_DOCKER_PATH:       $(TASK_UI_DOCKER_PATH)"
	@echo ""
	@echo "TASK_UI_RUN_VAR_NAME:      $(TASK_UI_RUN_VAR_NAME)"
	@echo "TASK_UI_RUN_PATH_WHICH:    $(TASK_UI_RUN_PATH_WHICH)"
	@echo ""
	@echo "TASK_UI_RUN_CMD:           $(TASK_UI_RUN_CMD)"
	@echo ""


task-web-run-h:
	$(BASE_RUN) -h
task-web-run-version:
	$(BASE_BIN_TARGET) --version

task-web-run-docker:
	cd $(TASK_UI_DOCKER_PATH) && task run

task-web-run-server:
	# http://localhost:3000
	$(TASK_UI_RUN_CMD) 


