
# add sqlite to path

SQLITE_BIN_PATH=/opt/homebrew/opt/sqlite/bin
SQLITE_LIB_PATH=/opt/homebrew/opt/sqlite/lib

export PATH:="/opt/homebrew/opt/sqlite/bin:$(PATH)"


#For compilers to find sqlite you may need to set:
export LDFLAGS="-L/opt/homebrew/opt/sqlite/lib"
export CPPFLAGS="-I/opt/homebrew/opt/sqlite/include"

BIN_ROOT=$(PWD)/.bin

print:

	echo 
	echo $(BIN_ROOT)
	echo
	echo $(shell which sqlite3)
	echo

	type -a sqlite3 
	echo
	sqlite3 --version
	echo

	cd $(SQLITE_BIN_PATH) && ./sqlite3 --version
	echo

	ls -al $(SQLITE_BIN_PATH)
	echo

dep:
	# we need non apple sqlite
	brew install sqlite3
dep-del:
	brew uninstall sqlite3

git:
	git clone https://git.sr.ht/~mil/mobsql
	echo mobsql >> .gitignore

test:
	cd mobsql && go test -tags=sqlite_math_functions -v ./...

build:
	mkdir -p $(BIN_ROOT)
	cd mobsql && go build -o $(BIN_ROOT)/mobsql -v -tags=sqlite_math_functions cli/*.go

install:

	go install -tags=sqlite_math_functions git.sr.ht/~mil/mobsql/cli@latest

	#go install git.sr.ht/~mil/transito@latest
	# has the db inside.
	# also using https://github.com/gioui-plugins/gio-plugins