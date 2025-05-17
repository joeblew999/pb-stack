print:

all: src mod-tidy run mod-tidy

src: src-del
	git clone https://github.com/delaneyj/geck
	@echo geck >> .gitignore
src-del:
	rm -rf geck 

mod-upgrade:
	cd geck && go-mod-upgrade --force
mod-tidy:
	cd geck && go-mod-upgrade --force

bin:
	cd geck && touch go.work
	cd geck && go work use .

	# slightly changed task file.
	cp Taskfile.yml geck/Taskfile.yml

	cd geck && task tools
	cd geck && task pb
	cd geck && task qtc
	cd geck && task templ

run:
	# http://localhost:8080/sparsesets/
	# http://localhost:8080/sparsesets/names has some data
	cd geck && task example




