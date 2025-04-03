

# os
MAKE_OS_NAME := $(shell go env GOOS)
MAKE_OS_ARCH := $(shell go env GOARCH)

# COMSPEC
MAKE_SHELL=${SHELL}
ifeq ($(MAKE_OS_NAME),windows)
	MAKE_SHELL=$(shell COMSPEC)
endif

# task
MAKE_TASK_VERSION=$(shell task --version)
MAKE_TASK_EXP=$(shell task --experiments)

print:
	@echo ""
	@echo "MAKE_OS_NAME:         $(MAKE_OS_NAME)"
	@echo ""
	@echo "MAKE_SHELL:           $(MAKE_SHELL)"
	@echo "MAKE_SHELL_VERSION:   $(shell zsh --version)"
	
	@echo ""
	@echo "MAKE_TASK_VERSION:    $(MAKE_TASK_VERSION)"
	@echo "MAKE_TASK_EXP:        $(MAKE_TASK_EXP)"

git-print:
	

dep-print:
	@echo ""
	@echo ls -al /etc/bash_completion.d/
	@echo ls -al /usr/local/share/zsh/site-functions
	@echo ls -al ~/.config/fish/completions/

dep-edit:
	code /Users/apple/.zshrc
	
dep:
	# assumes go installed

	# https://taskfile.dev/installation/	
	env GOBIN=$(PWD)/.bin go install github.com/go-task/task/v3/cmd/task@latest


	# NOT recommended. Install we added a line to the .zshrc
	# https://taskfile.dev/installation/
	task --completion bash  > $(PWD)/.bin/task.bash
	task --completion zsh  > $(PWD)/.bin/task.zsh
	task --completion fish  > $(PWD)/.bin/task.fish

	# copy to right location:
	#task --completion bash > /etc/bash_completion.d/task

	# zsh perm install
	#mkdir -p /usr/local/share/zsh/site-functions
	#sudo touch /usr/local/share/zsh/site-functions/_task
	#sudo chmod 777 /usr/local/share/zsh/site-functions/_task
	#task --completion zsh  > /usr/local/share/zsh/site-functions/_task

	

	#task --completion fish > ~/.config/fish/completions/task.fish
dep-del:
	sudo rm -f /usr/local/share/zsh/site-functions/_task

list:
	task --list-all