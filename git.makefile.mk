
### signing of src code	

BASE_SRC_SIGNING_CONFIG=$(HOME)/.gitconfig
BASE_SRC_SIGNING_CONFIG_LOCAL=$(BASE_SRC)/.git/config

BASE_SRC_SIGNING_USER_NAME=gedw99
BASE_SRC_SIGNING_USER_EMAIL=gedw99@gmail.com

# Thre is ssh and smime options. Pick one

# ssh based
# https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key#telling-git-about-your-ssh-key

BASE_SRC_SIGNING_KEY_PRIV=$(HOME)/.ssh/gedw99_github.com
BASE_SRC_SIGNING_KEY=$(HOME)/.ssh/gedw99_github.com.pub
BASE_SRC_SIGNING_PROGRAM=ssh
BASE_SRC_SIGNING_FORMAT=ssh

# PKI based using smime
# https://docs.github.com/en/authentication/managing-commit-signature-verification/telling-git-about-your-signing-key#telling-git-about-your-x509-key-1
#BASE_SRC_SIGNING_PATH=""
#BASE_SRC_SIGNING_KEY=e9abcabcef55053ead1dda7c297c6efd6c146a86
#BASE_SRC_SIGNING_PROGRAM=smimesign
#BASE_SRC_SIGNING_FORMAT=x509

BASE_SRC_SIGNING_BIN_NAME=smimesign
ifeq ($(BASE_OS_NAME),windows)
	BASE_SRC_SIGNING_BIN_NAME=smimesign.exe
endif
BASE_SRC_SIGNING_BIN_WHICH=$(shell command -v $(BASE_SRC_SIGNING_BIN_NAME))


base-src-sign-print:
	
	@echo "" 
	@echo "- bin" 
	@echo "BASE_SRC_SIGNING_BIN_NAME:      $(BASE_SRC_SIGNING_BIN_NAME)"
	@echo "BASE_SRC_SIGNING_BIN_WHICH:     $(BASE_SRC_SIGNING_BIN_WHICH)"
	@echo "" 
	@echo "" 
	@echo "- config" 
	@echo "BASE_SRC_SIGNING_CONFIG:        $(BASE_SRC_SIGNING_CONFIG)"
	@echo "BASE_SRC_SIGNING_CONFIG_LOCAL:  $(BASE_SRC_SIGNING_CONFIG_LOCAL)"
	@echo "" 
	@echo "" 
	@echo "- var" 
	@echo "" 
	@echo "BASE_SRC_SIGNING_USER_NAME:     $(BASE_SRC_SIGNING_USER_NAME)"
	@echo "BASE_SRC_SIGNING_USER_EMAIL:    $(BASE_SRC_SIGNING_USER_EMAIL)"
	@echo "" 

	@echo "BASE_SRC_SIGNING_KEY_CONFIG:    $(HOME)/.ssh/config"
	@echo "BASE_SRC_SIGNING_KEY_PRIV:      $(BASE_SRC_SIGNING_KEY_PRIV)"
	@echo "BASE_SRC_SIGNING_KEY:           $(BASE_SRC_SIGNING_KEY)"
	@echo "BASE_SRC_SIGNING_PROGRAM:       $(BASE_SRC_SIGNING_PROGRAM)"
	@echo "BASE_SRC_SIGNING_FORMAT:        $(BASE_SRC_SIGNING_FORMAT)"
	@echo ""

base-src-sign-print-scope:
	@echo ""
	@echo "- git config scope"
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --list --show-scope --show-origin
	@echo ""

base-src-sign-dep:
	@echo ""
	@echo " $(BASE_SRC_SIGNING_BIN_NAME) dep check ... "


ifeq ($(BASE_SRC_SIGNING_BIN_WHICH), )
	@echo ""
	@echo "$(BASE_SRC_SIGNING_BIN_NAME) dep check: failed"
	$(MAKE) base-src-sign-dep-single
else
	@echo ""
	@echo "$(BASE_SRC_SIGNING_BIN_NAME) dep check: passed"
endif

base-src-sign-dep-single:
	# https://github.com/github/smimesign
	#go install github.com/github/smimesign@latest
	brew install smimesign

base-src-sign-run-h: base-src-sign-dep
	$(BASE_SRC_SIGNING_BIN_NAME) -h
	$(BASE_SRC_SIGNING_BIN_NAME) -v
base-src-sign-run-list: base-src-sign-dep
	$(BASE_SRC_SIGNING_BIN_NAME) --list-keys
base-src-sign-run-sign: base-src-sign-dep
	#$(BASE_SRC_SIGNING_BIN_NAME) --sign -h
	$(BASE_SRC_SIGNING_BIN_NAME) --sign USER-ID=$(BASE_SRC_SIGNING_USER_NAME)
	

base-src-sign-set:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.name $(BASE_SRC_SIGNING_USER_NAME)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.email  $(BASE_SRC_SIGNING_USER_EMAIL)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.signingkey $(BASE_SRC_SIGNING_KEY)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.program $(BASE_SRC_SIGNING_PROGRAM)
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.format $(BASE_SRC_SIGNING_FORMAT)
base-src-sign-get:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --get user.name
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --get user.email
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.signingkey
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.program
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.format
base-src-sign-del:
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.name ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.email ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local user.signingkey ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.program ""
	cd $(BASE_SRC) && $(BASE_DEP_BIN_GIT_NAME) config --local gpg.format ""

