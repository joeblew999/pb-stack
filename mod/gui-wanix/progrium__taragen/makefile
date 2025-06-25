#REPO=taragen/web

REPO=wanix.sh

print:

dep:
	go install github.com/progrium/taragen/cmd/taragen@latest

git:
	#git clone https://github.com/progrium/taragen
	#echo taragen >> .gitignore

	git clone https://github.com/wanixdev/wanix.sh
	echo wanix.sh >> .gitignore


taragen-serve:
	# http://localhost:8088
	cd $(REPO) && taragen serve
taragen-build:
	cd $(REPO) && taragen build

wanix:
	# curl from https://github.com/tractordev/wanix/releases/tag/v0.3-preview
	./wanix
	
	