# Project-86-Launcher

NAME=Project-86-Launcher

print:

all: src mod-tidy run mod-tidy

src: src-del
	git clone https://github.com/Project-86-Community/$(NAME) -b next
	@echo $(NAME) >> .gitignore
src-del:
	rm -rf $(NAME) 

mod-upgrade:
	cd $(NAME) && go-mod-upgrade --force
mod-tidy:
	cd $(NAME) && go-mod-upgrade --force


bin:
	cd $(NAME) && touch go.work
	cd $(NAME) && go work use .

	# slightly changed task file.
	#cp Taskfile.yml $(NAME)/Taskfile.yml

	#cd $(NAME) && task build

	cd $(NAME) && task build:windows

run:
	# http://localhost:8080/sparsesets/
	# http://localhost:8080/sparsesets/names has some data
	cd $(NAME) && task run

run-debug:
	# http://localhost:8080/sparsesets/
	# http://localhost:8080/sparsesets/names has some data
	cd $(NAME) && task run:debug




