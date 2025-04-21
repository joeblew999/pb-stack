# cloudflare

- TASK file 
- TOFU file for automating this.

## Setup of tokens

Your .env needs to have:

```sh

# https://dash.cloudflare.com
CLOUDFLARE_ACCOUNT_ID=7384af54e33b8a54ff240371ea368440

# https://dash.cloudflare.com/$CLOUDFLARE_ACCOUNT_ID
# example: https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/api-tokens
CF_API_TOKEN=xxx
CF_API_TOKEN=2kVTxWum_I7Ts2wf8IzlIh09ZD3so6o-J-1h8MSM

```

### 0. Get your Account ID

CLOUDFLARE_ACCOUNT_ID is found on your Cloudflare home page at:

https://dash.cloudflare.com/$CLOUDFLARE_ACCOUNT_ID

Example:

https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440


### 1. Create an Account Token

CF_API_TOKEN is created at:

https://dash.cloudflare.com/$CLOUDFLARE_ACCOUNT_ID/api-tokens

Example: 

https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/api-tokens


### 2. Create API Tokens

User API Tokens

Manage access and permissions for your accounts, sites, and products

Create your API token: https://dash.cloudflare.com/profile/api-tokens

TOFU can automate creation of these for you.

## cloudflare cli

Does not exist anymore. They deleted the cmd. bugger !!

https://github.com/cloudflare/cloudflare-go/forks?include=active&page=1&period=2y&sort_by=last_updated

## cloudflare tunnel

https://github.com/cloudflare/cloudflared

For protecting and exposing Servers in VPS and On Premise.



## wrangler workers

https://github.com/syumai/workers

For writing things in golang.

https://developers.cloudflare.com/workers/wrangler/system-environment-variables/

Env variables.

## terraform

https://github.com/cloudflare/terraform-provider-cloudflare

For provisioning and configuring.

---

https://github.com/cloudflare/cf-terraforming

Retrieves your configurations from the Cloudflare API and converting them to Terraform configurations

