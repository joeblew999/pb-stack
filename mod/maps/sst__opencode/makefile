BIN=$(PWD)/.bin/opencode

CONFIG=$(PWD)

print:


install:
	mkdir -p $(PWD)/.bin
	# https://github.com/sst/opencode/releases/tag/v0.0.52
	go install github.com/sst/opencode@v0.0.52
	mv $(GOPATH)/bin/opencode $(PWD)/.bin/opencode

run:
	$(BIN) -d -c $(PWD)

	# creates a sqlite locally.
	# ! Cant exit

run-prompt:
	$(BIN) -p "Explain the use of context in Go"

run-cli:
	$(BIN) -h
	$(BIN) --version
	$(BIN) --cwd


config:
	code $(HOME)/.opencode.json