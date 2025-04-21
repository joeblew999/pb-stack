### CLOUDFLARE

# https://github.com/cloudflare/cloudflare-go/

# https://developers.cloudflare.com/api

# terra / tofu
# https://github.com/cloudflare/cf-terraforming
# https://github.com/cloudflare/terraform-provider-cloudflare 


FLARE_CLI_DEP=flare-cli
FLARE_CLI_DEP_BIN_DEEP=flarectl

FLARE_CLI_DEP_REPO=cloudflare-go
FLARE_CLI_DEP_REPO_URL=https://github.com/cloudflare/cloudflare-go
FLARE_CLI_DEP_REPO_DEEP=$(FLARE_CLI_DEP_REPO)/cmd/flarectl

FLARE_CLI_DEP_MOD=github.com/cloudflare/cloudflare-go

# I WANT THIS as its github.com/cloudflare/cloudflare-go/v4
# Do NOT use this. The big version numbers do not have the cmd in it. No idea why !
# https://github.com/cloudflare/cloudflare-go/blob/v4.1.0/cmd/flarectl

# SO who is using it: https://github.com/cloudflare/cloudflare-go/network/dependents?package_id=UGFja2FnZS01NjI1ODY1MzU4 
# https://github.com/cloudflare/terraform-provider-cloudflare/blob/main/go.mod

# Uses github.com/cloudflare/cloudflare-go/, so is OLD.
# Use the little version numbers.
# https://github.com/cloudflare/cloudflare-go/blob/v0.112.0/cmd/flarectl
# https://github.com/cloudflare/cloudflare-go/blob/v0.114.0/cmd/flarectl
# https://github.com/cloudflare/cloudflare-go/blob/v0.115.0/cmd/flarectl

FLARE_CLI_DEP_MOD_DEEP=$(FLARE_CLI_DEP_MOD)/cmd/flarectl

# https://github.com/cloudflare/cloudflare-go/releases/tag/v0.112.0
# https://github.com/cloudflare/cloudflare-go/releases/tag/v0.114.0
# https://github.com/cloudflare/cloudflare-go/releases/tag/v0.115.0
FLARE_CLI_DEP_VERSION=v0.115.0

FLARE_CLI_DEP_NATIVE=$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_NATIVE)
FLARE_CLI_DEP_WHICH=$(shell $(BASE_DEP_BIN_WHICH_NAME) $(FLARE_CLI_DEP_NATIVE))

# Match Tag version
# https://github.com/cloudflare/cloudflare-go/blob/v0.115.0/cmd/flarectl/flarectl.go#L11
# we only need "main.version", because we cd into the cmd/flarectl
FLARE_CLI_GO_BUILD_LDFLAGS="-s -w -X main.version=$(FLARE_CLI_DEP_VERSION)"
FLARE_CLI_GO_BUILD_CMD=CGO_ENABLED=0 $(BASE_DEP_BIN_GO_NAME) build -ldflags $(FLARE_CLI_GO_BUILD_LDFLAGS) -v


### meta

FLARE_CLI_DEP_META_PREFIX=flare_cli
FLARE_CLI_DEP_META_NAME=$(FLARE_CLI_DEP_META_PREFIX)_meta.json
FLARE_CLI_DEP_META_WHICH=$(BASE_MAKE_IMPORT)/$(FLARE_CLI_DEP_META_NAME)

FLARE_CLI_DEP_TEMPLATE=$(PWD)/$(FLARE_CLI_DEP_META_PREFIX).mk


flare-cli-print:
	@echo ""
	@echo "- bin"
	@echo "FLARE_CLI_DEP:              $(FLARE_CLI_DEP)"
	@echo "FLARE_CLI_DEP_VERSION:      $(FLARE_CLI_DEP_VERSION)"
	@echo "FLARE_CLI_DEP_WHICH:        $(FLARE_CLI_DEP_WHICH)"
	@echo "FLARE_CLI_DEP_NATIVE:       $(FLARE_CLI_DEP_NATIVE)"
	@echo ""
	@echo "FLARE_CLI_DEP_VERSION:      $(FLARE_CLI_DEP_VERSION)"
	@echo ""
	@echo "FLARE_CLI_GO_BUILD_LDFLAGS: $(FLARE_CLI_GO_BUILD_LDFLAGS)"
	@echo "FLARE_CLI_GO_BUILD_CMD:     $(FLARE_CLI_GO_BUILD_CMD)"
	@echo ""
	@echo "- meta"
	@echo "FLARE_CLI_DEP_META_PREFIX:  $(FLARE_CLI_DEP_META_PREFIX)"
	@echo "FLARE_CLI_DEP_META_NAME:    $(FLARE_CLI_DEP_META_NAME)"
	@echo "FLARE_CLI_DEP_META_WHICH:   $(FLARE_CLI_DEP_META_WHICH)"
	@echo "FLARE_CLI_DEP_TEMPLATE:     $(FLARE_CLI_DEP_TEMPLATE)"
	@echo ""


### base

## list base files
flare-cli-base-list:
	$(BASE_DEP_BIN_TREE_NAME) -h $(BASE_MAKE_IMPORT)/$(FLARE_CLI_DEP_META_PREFIX)*
## edit base files
flare-cli-base-edit:
	$(VSCODE_BIN_NAME) $(BASE_MAKE_IMPORT)/$(FLARE_CLI_DEP_META_PREFIX)*


### dep

## edit deployed files
flare-cli-dep-edit:
	$(VSCODE_BIN_NAME) $(BASE_CWD_DEP)/$(FLARE_CLI_DEP_META_PREFIX)*

## list deployed files 
flare-cli-dep-list:
	$(BASE_DEP_BIN_TREE_NAME) $(BASE_CWD_DEP)/$(FLARE_CLI_DEP_META_PREFIX)*


## copies a config your cwd, so you can customise it for your project.
flare-cli-dep-template-config:

	@echo ""
	@echo "copying template to your cwd ..."

	# create a meta in base cwd
	#rm -rf $(FLARE_CLI_DEP_META_WHICH)
	#echo $(FLARE_CLI_DEP_VERSION) >> $(FLARE_CLI_DEP_META_WHICH)
	
	# copy the configs to base cwd
	cp $(BASE_MAKE_IMPORT)/$(FLARE_CLI_DEP_META_PREFIX)_config* $(BASE_CWD)


### dep

flare-cli-dep-template: base-dep-init
	@echo ""
	@echo "-version"
	rm -rf $(FLARE_CLI_DEP_META_WHICH)
	echo $(FLARE_CLI_DEP_VERSION) >> $(FLARE_CLI_DEP_META_WHICH)

	# templates to dep.
	cp $(BASE_MAKE_IMPORT)/$(FLARE_CLI_DEP_META_PREFIX)* $(BASE_CWD_DEP)



### dep bin

flare-cli-dep-del:
	rm -f $(FLARE_CLI_DEP_WHICH)


flare-cli-dep: 
	@echo ""
	@echo "$(FLARE_CLI_DEP) dep check ... "
	@echo ""
	@echo "FLARE_CLI_DEP_WHICH: $(FLARE_CLI_DEP_WHICH)"

ifeq ($(FLARE_CLI_DEP_WHICH), )
	@echo ""
	@echo "$(FLARE_CLI_DEP) dep check: failed"
	$(MAKE) flare-cli-dep-single
else
	@echo ""
	@echo "$(FLARE_CLI_DEP) dep check: passed"
endif

flare-cli-dep-start:
	rm -rf $(BASE_CWD_DEPTMP)
	mkdir -p $(BASE_CWD_DEPTMP)

	cd $(BASE_CWD_DEPTMP) && $(BASE_DEP_BIN_GIT_NAME) clone $(FLARE_CLI_DEP_REPO_URL) -b $(FLARE_CLI_DEP_VERSION) --single-branch
	cd $(BASE_CWD_DEPTMP) && echo $(FLARE_CLI_DEP_REPO) >> .gitignore
	cd $(BASE_CWD_DEPTMP) && touch go.work
	cd $(BASE_CWD_DEPTMP) && $(BASE_DEP_BIN_GO_NAME) work use $(FLARE_CLI_DEP_REPO)

flare-cli-dep-end:
	rm -rf $(BASE_CWD_DEPTMP)

flare-cli-dep-single: flare-cli-dep-template

	$(MAKE) flare-cli-dep-start

ifeq ($(BASE_OS_NAME),darwin)
	@echo "--- darwin ---"
	$(MAKE) flare-cli-dep-darwin
endif
ifeq ($(BASE_OS_NAME),linux)
	@echo "--- linux ---"
	$(MAKE) flare-cli-dep-linux
endif
ifeq ($(BASE_OS_NAME),windows)
	@echo "--- windows ---"
	$(MAKE) flare-cli-dep-windows
endif

	$(MAKE) flare-cli-dep-end

flare-cli-dep-all: flare-cli-dep-template

	$(MAKE) flare-cli-dep-start
	
	$(MAKE) flare-cli-dep-darwin
	$(MAKE) flare-cli-dep-linux
	$(MAKE) flare-cli-dep-windows
	
	$(MAKE) flare-cli-dep-end

flare-cli-dep-darwin:
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_CLI_DEP_REPO_DEEP) && GOOS=darwin GOARCH=amd64 $(FLARE_CLI_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_DARWIN_AMD64)
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_CLI_DEP_REPO_DEEP) && GOOS=darwin GOARCH=arm64 $(FLARE_CLI_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_DARWIN_ARM64)
flare-cli-dep-linux:
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_CLI_DEP_REPO_DEEP) && GOOS=linux GOARCH=amd64 $(FLARE_CLI_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_LINUX_AMD64)
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_CLI_DEP_REPO_DEEP) && GOOS=linux GOARCH=arm64 $(FLARE_CLI_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_LINUX_ARM64)
flare-cli-dep-windows:
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_CLI_DEP_REPO_DEEP) && GOOS=windows GOARCH=amd64 $(FLARE_CLI_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_WINDOWS_AMD64)
	cd $(BASE_CWD_DEPTMP) && cd $(FLARE_CLI_DEP_REPO_DEEP) && GOOS=windows GOARCH=arm64 $(FLARE_CLI_GO_BUILD_CMD) -o $(BASE_CWD_DEP)/$(FLARE_CLI_DEP)_$(BASE_BIN_SUFFIX_WINDOWS_ARM64)





### run

FLARE_CLI_RUN_VAR_ACCOUNT_EMAIL=xxx@gmail.com

# https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440
# Opening the Dashboard and used in CLI.
FLARE_CLI_RUN_VAR_ACCOUNT_ID=xxx


# User token or an Account owned token ?

# 1. Go see what Account Token#s exist.
# Account token's are useful for when its going to be used for a long time.
# Account Tokens Dash: 
# FLARE_CLI_RUN_DASH_ACCOUNT_TOKEN=https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/api-tokens
FLARE_CLI_RUN_DASH_ACCOUNT_TOKEN=https://dash.cloudflare.com/$(FLARE_CLI_RUN_VAR_ACCOUNT_ID)/api-tokens

# 2. Use one of them that fits your needs.
# The Actual Tokens: These SHOULD be in a Vault.
# https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/api-tokens
# A. RIGHTS: Read all resources
FLARE_CLI_RUN_VAR_ACCOUNT_TOKEN=xxx
# B. RIGHTS: Edit all DNS.
#FLARE_CLI_RUN_VAR_ACCOUNT_TOKEN=Ip3ETfHvdwQSlIlzGgwnOKDaBNeFNH0_7y1qgymF

#3. There is no 3.. I prefer Account Tokens for now.
# User Tokens are in the Profiles section.
# https://dash.cloudflare.com/profile/api-tokens
# Useful for getting API Keys.
FLARE_CLI_RUN_VAR_USER_TOKEN=?
# A. TERRAFORM: Edit lots of shit.
FLARE_CLI_RUN_VAR_USER_TOKEN=xxxx


# MAP the chosen token to the ENV variable that the CLI uses.
# It can ONLY be done via ENV variable. 
# Pick either ACCOUNT or USER token
#export CF_API_TOKEN=$(FLARE_CLI_RUN_VAR_ACCOUNT_TOKEN)
export CF_API_TOKEN=$(FLARE_CLI_RUN_VAR_USER_TOKEN)

# After creating your first API token, you can create additional API tokens via the API !!
# https://developers.cloudflare.com/fundamentals/api/how-to/create-via-api/

FLARE_CLI_RUN_VAR_ZONE=ubuntudesign.com
#FLARE_CLI_RUN_VAR_ZONE=ubuntusoftware.net


# GLOBAL OPTIONS:
#   --account-id value  Optional account ID [$CF_ACCOUNT_ID]
#   --json              show output as JSON instead of as a table (default: false)
#   --help, -h          show help
#   --version, -v       print the version

FLARE_CLI_RUN_PATH=$(PWD)
FLARE_CLI_RUN_CMD=$(FLARE_CLI_DEP_NATIVE) --account-id=$(FLARE_CLI_RUN_VAR_ACCOUNT_ID)


flare-cli-run-print:
	@echo ""
	@echo ""
	@echo "- run"
	@echo "FLARE_CLI_RUN_VAR_ACCOUNT_EMAIL:   $(FLARE_CLI_RUN_VAR_ACCOUNT_EMAIL)"
	@echo "FLARE_CLI_RUN_VAR_ACCOUNT_ID:      $(FLARE_CLI_RUN_VAR_ACCOUNT_ID)"
	@echo ""
	@echo "- run dash"
	@echo "FLARE_CLI_RUN_DASH_ACCOUNT_TOKEN:  $(FLARE_CLI_RUN_DASH_ACCOUNT_TOKEN)"
	@echo "- run token"
	@echo "FLARE_CLI_RUN_VAR_ACCOUNT_TOKEN:   $(FLARE_CLI_RUN_VAR_ACCOUNT_TOKEN)"
	
	@echo ""
	@echo "- run"
	@echo "FLARE_CLI_RUN_VAR_ZONE:            $(FLARE_CLI_RUN_VAR_ZONE)"
	
	
	@echo "CF_API_TOKEN ( EXPORTED ):         $(CF_API_TOKEN)"
	
	@echo ""
	@echo "FLARE_CLI_RUN_PATH:                $(FLARE_CLI_RUN_PATH)"
	@echo "FLARE_CLI_RUN_CMD:                 $(FLARE_CLI_RUN_CMD)"


flare-cli-run-h: flare-cli-dep
	$(FLARE_CLI_RUN_CMD) -h

flare-cli-run-version: flare-cli-dep
	$(FLARE_CLI_RUN_CMD) --version

flare-cli-run-dns-h:
	$(FLARE_CLI_RUN_CMD) dns -h
flare-cli-run-dns-list:
	#$(FLARE_CLI_RUN_CMD) dns list -h
	# --id value       record id
	# --zone value     zone name
	$(FLARE_CLI_RUN_CMD) dns list --zone=$(FLARE_CLI_RUN_VAR_ZONE)
flare-cli-run-dns-create:
    # I use Create or Update so its idempotent
    # https://www.thegeekdiary.com/flarectl-official-cli-for-cloudflare/
	$(FLARE_CLI_RUN_CMD) dns create-or-update -h
    # --zone value      zone name
    # --name value      record name
    # --content value   record content
    #$(FLARE_CLI_RUN_CMD) dns create-or-update --zone="$(FLARE_CLI_RUN_VAR_ZONE)" --name="app" --type="CNAME" --content="myapp.herokuapp.com" --proxy
flare-cli-run-dns-delete:
	#  Delete a DNS record
	$(FLARE_CLI_RUN_CMD) dns delete -h
    # --zone value  zone name
    # --id value    record id
	$(FLARE_CLI_RUN_CMD) dns delete --zone=$(FLARE_CLI_RUN_VAR_ZONE) --id="?"

### firewall


flare-cli-run-firewall-rules-list-h:
	$(FLARE_CLI_RUN_CMD) firewall rules list -h
    #   --zone value        zone name
    #   --account value     account name
    #   --value value       rule value
    #   --scope-type value  rule scope
    #   --mode value        rule mode
    #   --notes value       rule notes
flare-cli-run-firewall-rules-list:
	# works
	$(FLARE_CLI_RUN_CMD) firewall rules list --zone=$(FLARE_CLI_RUN_VAR_ZONE)
	# works also with extra values
	$(FLARE_CLI_RUN_CMD) firewall rules list --zone=$(FLARE_CLI_RUN_VAR_ZONE)
flare-cli-run-firewall-rules-create-h:
	# Block a specific IP
	# https://www.thegeekdiary.com/flarectl-official-cli-for-cloudflare/
	$(FLARE_CLI_RUN_CMD) firewall rules create -h
    #   --zone value     zone name
    #   --account value  account name
    #   --value value    rule value
    #   --mode value     rule mode
    #   --notes value    rule notes
flare-cli-run-firewall-rules-create:
	# Block a specific IP
	# https://www.thegeekdiary.com/flarectl-official-cli-for-cloudflare/
	$(FLARE_CLI_RUN_CMD) firewall rules create --zone=$(FLARE_CLI_RUN_VAR_ZONE) --account="gedw99@gmail.com" --value="8.8.8.8" --mode="block" --notes="Block bad actor"
	
### ip

flare-cli-run-ip-h:
	# Print Cloudflare IP ranges
	$(FLARE_CLI_RUN_CMD) ips -h
    #   --ip-type value  type of IPs ( ipv4 | ipv6 | all ) (default: "all")
    #   --ip-only        show only addresses (default: false)
flare-cli-run-ip:
	# works
	$(FLARE_CLI_RUN_CMD) ips --ip-type all

### ocrc

flare-cli-run-ocrc-h:
	# WORKS- NO AUTH needed
	# Print Origin CA Root Certificate (in PEM format)
	$(FLARE_CLI_RUN_CMD) ocrc -h
    # --algorithm value  certificate algorithm ( ecc | rsa )
flare-cli-run-ocrc:
	$(FLARE_CLI_RUN_CMD) ocrc --algorithm=ecc

### pagerules ( Page Rules endpoint does not support account owned tokens. (1011) )

flare-cli-run-pagerules-h:
	$(FLARE_CLI_RUN_CMD) pagerules -h
flare-cli-run-pagerules-list:
	$(FLARE_CLI_RUN_CMD) pagerules list --zone=$(FLARE_CLI_RUN_VAR_ZONE)
	
### railgun ( DEPRECATED )

flare-cli-run-railgun-h:
	# THEY ARE DEPRECATED THIS. CF TUNNEL replaced it
	# https://blog.cloudflare.com/deprecating-railgun/
	# NOTE: Page Rules endpoint does not support account owned tokens. (1011)
	$(FLARE_CLI_RUN_CMD) railgun -h
flare-cli-run-railgun-list:
	$(FLARE_CLI_RUN_CMD) railgun --zone=$(FLARE_CLI_RUN_VAR_ZONE)
	

flare-cli-run-user-h:
	$(FLARE_CLI_RUN_CMD) user -h
    # info, i    User details
    # update, u  Update user details

flare-cli-run-user-info:
	# Works - not yet
	# Print  User details
	$(FLARE_CLI_RUN_CMD) user info
flare-cli-run-user-update-h:
	# Update user details
	$(FLARE_CLI_RUN_CMD) user update -h

### agents ( to block user agents )

flare-cli-run-user-agents-list-h:
	$(FLARE_CLI_RUN_CMD) user-agents list -h
    #   --zone value  zone name
    #   --page value  result page to return (default: 0)
flare-cli-run-user-agents-list:
	# works
	$(FLARE_CLI_RUN_CMD) user-agents list --zone=$(FLARE_CLI_RUN_VAR_ZONE)
flare-cli-run-user-agents-create-h:
	$(FLARE_CLI_RUN_CMD) user-agents create -h
    #   --zone value         zone name
    #   --mode value         the blocking mode: block, challenge, js_challenge, whitelist
    #   --value value        the exact User-Agent to block
    #   --paused             whether the rule should be paused (default: false) (default: false)
    #   --description value  a description for the rule
flare-cli-run-user-agents-create:
	# needs WRITE token
	$(FLARE_CLI_RUN_CMD) user-agents create --zone=$(FLARE_CLI_RUN_VAR_ZONE) --mode=block --value=?? --description="your blocked dude"
flare-cli-run-user-agents-update-h:
	$(FLARE_CLI_RUN_CMD) user-agents update -h
    #   --zone value         zone name
    #   --id value           User-Agent blocking rule ID
    #   --mode value         the blocking mode: block, challenge, js_challenge, whitelist
    #   --value value        the exact User-Agent to block
    #   --paused             whether the rule should be paused (default: false) (default: false)
    #   --description value  a description for the rule
flare-cli-run-user-agents-update:
	$(FLARE_CLI_RUN_CMD) user-agents update --zone=$(FLARE_CLI_RUN_VAR_ZONE) --id=1 
flare-cli-run-user-agents-delete-h:
	$(FLARE_CLI_RUN_CMD) user-agents delete -h
    #   --zone value  zone name
    #   --id value    User-Agent blocking rule ID
flare-cli-run-user-agents-delete:
	$(FLARE_CLI_RUN_CMD)user-agents delete --zone=$(FLARE_CLI_RUN_VAR_ZONE) --id=1

### zone

flare-cli-run-zone-h:
	$(FLARE_CLI_RUN_CMD) zone -h
#   list, l       List all zones on an account
#   create, c     Create a new zone
#   delete        Delete a zone
#   check         Initiate a zone activation check
#   info, i       Information on one zone
#   lockdown, lo  Lockdown a zone based on config
#   plan, p       Plan information for one zone
#   settings, s   Settings for one zone
#   purge         (Selectively) Purge the cache for a zone
#   dns, d        DNS records for a zone
#   railgun, r    Railguns for a zone
#   certs, ct     Custom SSL certificates for a zone
#   keyless, k    Keyless SSL for a zone
#   export, x     Export DNS records for a zone
#   help, h       Shows a list of commands or help for one command
flare-cli-run-zone-list:
	$(FLARE_CLI_RUN_CMD) zone list

flare-cli-run-zone-create-h:
	# Create many new Cloudflare zones automatically with names from domains.txt...
	# Information on one zone
	$(FLARE_CLI_RUN_CMD) zone create -h
    #   --zone value        zone name
    #   --jumpstart         automatically fetch DNS records (default: false)
    #   --account-id value  account ID
flare-cli-run-zone-create:
	$(FLARE_CLI_RUN_CMD) zone create --zone=$(FLARE_CLI_RUN_VAR_ZONE)

flare-cli-run-zone-info-h:
	# Create many new Cloudflare zones automatically with names from domains.txt...
	# Information on one zone
	$(FLARE_CLI_RUN_CMD) zone info -h
    #    --zone value  zone name
flare-cli-run-zone-info:
	$(FLARE_CLI_RUN_CMD) zone info --zone=$(FLARE_CLI_RUN_VAR_ZONE)

flare-cli-run-zone-plan-h:
	$(FLARE_CLI_RUN_CMD) zone plan -h
    #    --zone value  zone name
flare-cli-run-zone-plan:
	# They have not finished it yet.
	$(FLARE_CLI_RUN_CMD) zone plan --zone=$(FLARE_CLI_RUN_VAR_ZONE)
    #    --zone value  zone name

flare-cli-run-zone-delete-h:
	$(FLARE_CLI_RUN_CMD) zone plan
    #    --zone value  zone name
flare-cli-run-zone-delete:
	# VERY MUCH CHECK FIRST !!!
	#$(FLARE_CLI_RUN_CMD) zone delete --zone=$(FLARE_CLI_RUN_VAR_ZONE)

flare-cli-run-zone-check:
	$(FLARE_CLI_RUN_CMD) zone check --zone=$(FLARE_CLI_RUN_VAR_ZONE)
flare-cli-run-zone-settings:
	# Does not do anything yet. 
	$(FLARE_CLI_RUN_CMD) zone settings


### terraforming

flare-cli-run-terra-dep:
	# https://github.com/cloudflare/cf-terraforming/tree/master?tab=readme-ov-file#example-usage
	# https://github.com/cloudflare/cf-terraforming/releases/tag/v0.23.2
	go install github.com/cloudflare/cf-terraforming/cmd/cf-terraforming@v0.23.2

flare-cli-run-terra-h: flare-cli-run-terra-dep
	cf-terraforming -h
flare-cli-run-terra-version: flare-cli-run-terra-dep
	cf-terraforming version
flare-cli-run-terra-generate: flare-cli-run-terra-dep
	# Fetch resources from the Cloudflare API and generate the respective Terraform stanzas

	$(MAKE) tofu-run-init

	cf-terraforming --provider-registry-hostname registry.opentofu.org \
		--terraform-binary-path $(TOFU_DEP_WHICH) \
		--modern-import-block \
		--email $(FLARE_CLI_RUN_VAR_ACCOUNT_EMAIL) \
		--token $(CF_API_TOKEN) \
		generate \
		--zone $(FLARE_CLI_RUN_VAR_ZONE) \
		--resource-type "cloudflare_zone"

flare-cli-run-terra-import: flare-cli-run-terra-dep
	# Output `terraform import` compatible commands in order to import resources into state

	$(MAKE) tofu-run-init
	

	cf-terraforming import -h

	#cf-terraforming import --provider-registry-hostname registry.opentofu.org \
		--terraform-binary-path $(TOFU_DEP_WHICH) \
		--modern-import-block \
		--email $(FLARE_CLI_RUN_VAR_ACCOUNT_EMAIL) \
		--token $(CF_API_TOKEN) \
		--zone $(FLARE_CLI_RUN_VAR_ZONE) \
		--resource-type cloudflare_zone









	