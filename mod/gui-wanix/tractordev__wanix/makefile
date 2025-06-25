

git:
	git clone https://github.com/tractordev/wanix/
	echo wanix >> .gitignore

bin:
	cd wanix && $(MAKE) all
###

wanix-help:
	# curl from https://github.com/tractordev/wanix/releases/tag/v0.3-preview
	./wanix-bin-bin

wanix-serve:
	# Runs the Seerver and the wanix browser gui.
	# http://localhost:7654
	./wanix-bin serve

wanix-mount:
	./wanix-bin mount
	# 2025/05/28 22:52:44 main.go:35: no FUSE mount utility found
wanix-console:
	./wanix-bin console
	# 2025/05/28 22:53:19 non-root open: wanix.bundle.js (=> fskit.UnionFS wanix.bundle.js)
wanix-export:
	./wanix-bin export
	
	