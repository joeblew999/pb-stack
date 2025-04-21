# CLOUDFLARE WRANGLER

to be done.

## install

brew install cloudflare-wrangler2
brew uninstall cloudflare-wrangler2

or

npm install -g @cloudflare/wrangler

or

https://denoflare.dev/cli/
- works with deno ( which is rust )

brew install deno

/Users/apple/.deno/bin/denoflare
ℹ️  Add /Users/apple/.deno/bin to PATH
    export PATH="/Users/apple/.deno/bin:$PATH"

brew uninstall deno

deno install --unstable-worker-options --allow-read --allow-net --allow-env --allow-run --name denoflare --force https://raw.githubusercontent.com/skymethod/denoflare/v0.6.0/cli/cli.ts

 Its a binary, so easy for me to manage.

denoflare --help

deno uninstall denoflare
 deleted /Users/apple/.deno/bin/denoflare




## other

cloudflare-wrangler2 on brew
https://formulae.brew.sh/formula/cloudflare-wrangler2

The actual repo is:
https://github.com/cloudflare/workers-sdk
- no binaries !

https://github.com/cloudflare/workers-sdk/releases/tag/wrangler%403.78.7

https://registry.npmjs.org/wrangler/-/wrangler-3.78.12.tgz

https://developers.cloudflare.com/workers/wrangler/install-and-update/

ensure you have Node.js ↗ and npm ↗ installed.

## example

https://seanrmurphy.medium.com/building-an-openapi-compatible-api-for-cloudflare-workers-in-go-dff28e73dcfa
https://github.com/seanrmurphy/cf-race-api

uses swagger for API and gen db
uses golang workers !
uses pages
uses d1



## env

https://developers.cloudflare.com/workers/wrangler/system-environment-variables/


https://developers.cloudflare.com/workers/wrangler/ci-cd/

1. Wrangler requires a Cloudflare API token and account ID to authenticate with the Cloudflare API.

API token
To create an API token to authenticate Wrangler in your CI job:
Log in to the Cloudflare dashboard ↗.
Select My Profile > API Tokens.
Select Create Token > find Edit Cloudflare Workers > select Use Template.
Customize your token name.
Scope your token.
You will need to choose the account and zone resources that the generated API token will have access to. We recommend scoping these down as much as possible to limit the access of your token. For example, if you have access to three different Cloudflare accounts, you should restrict the generated API token to only the account on which you will be deploying a Worker.


2. Set up CI


The method for running Wrangler in your CI/CD environment will depend on the specific setup for your project (whether you use GitHub Actions/Jenkins/GitLab or something else entirely).
To set up your CI:
Go to your CI platform and add the following as secrets:
CLOUDFLARE_ACCOUNT_ID: Set to the Cloudflare account ID for the account on which you want to deploy your Worker.
CF_API_TOKEN: Set to the Cloudflare API token you generated.
Warning
It is important not to store the value of CF_API_TOKEN in your repository, as it gives access to deploy Workers on your account. Instead, you should utilise your CI/CD provider’s support for storing secrets.

some common ones with flare, etc

lets get this common...