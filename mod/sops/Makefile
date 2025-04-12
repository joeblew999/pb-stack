
SOPS_BASE_OS_NAME := $(shell go env GOOS)


# https://github.com/getsops/sops
SOPS_BIN_NAME=sops
ifeq ($(SOPS_BASE_OS_NAME),windows)
	SOPS_BIN_NAME=sops.exe
endif
SOPS_BIN_VERSION=latest
# https://github.com/getsops/sops/releases/tag/v3.10.1
#SOPS_BIN_VERSION=v3.10.1
SOPS_BIN_WHICH=$(shell command -v $(SOPS_BIN_NAME))
SOPS_BIN_WHICH_VERSION=$(shell $(SOPS_BIN_NAME) --disable-version-check --version)

sops-print:
	@echo ""
	@echo ""
	@echo "SOPS_BIN_NAME:            $(SOPS_BIN_NAME)"
	@echo "SOPS_BIN_VERSION:         $(SOPS_BIN_VERSION)"
	@echo "SOPS_BIN_WHICH:           $(SOPS_BIN_WHICH)"
	@echo "SOPS_BIN_WHICH_VERSION:   $(SOPS_BIN_WHICH_VERSION)"
	@echo ""
sops-del:
	rm -rf $(SOPS_BIN_WHICH)
sops:
ifeq ($(SOPS_BIN_WHICH), )
	@echo ""
	@echo "$(SOPS_BIN_NAME) dep check: failed, so installing ..."
	go install github.com/getsops/sops/v3/cmd/sops@$(SOPS_BIN_VERSION)
else
	@echo ""
	@echo "$(SOPS_BIN_NAME) dep check: passed"
endif
sops-run-encrypt-h: sops
	$(SOPS_BIN_NAME) encrypt -h
sops-run-encrypt: sops
	@echo ""
	@echo "encrypt ..."
	@echo ""
	rm -f test.sops.enc.*

	# .env
	#$(SOPS_BIN_NAME) encrypt --input-type dotenv --age age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw test.sops.env > test.sops.enc.env
	# yml
	$(SOPS_BIN_NAME) encrypt --input-type yaml --age age1yt3tfqlfrwdwx0z0ynwplcr6qxcxfaqycuprpmy89nr83ltx74tqdpszlw test.sops.yml > test.sops.enc.yml
	
	$(SOPS_BIN_NAME) updatekeys test.sops.enc.yml
sops-run-decrypt: sops
	@echo ""
	@echo "decrypt ..."
	@echo ""
	rm -f test.sops.dec.*

	# .env
	#cat test.sops.enc.env | $(SOPS_BIN_NAME) decrypt --input-type dotenv --output-type dotenv /dev/stdin > test.sops.dec.env
	# yml
	cat test.sops.enc.yml | $(SOPS_BIN_NAME) decrypt --input-type yaml --output-type yaml /dev/stdin > test.sops.dec.yml

	
	
	

