# github workflows

REF: https://github.com/openmcp-project/mcp-operator uses TASK like we do.

https://github.com/mvdan/github-actions-golang has good best practices and workarounds too.


Its a mono repo, driven by TASK files, so its easy to deploy as many projects as you want from this repo with a single setup. 

Tests for all platforms are also easy.

## CI

Builds whatever is configured by the Root TASK file.

The binaries are stored in github but not signed, so the Pack manifest will call the OS to ask for them to be run the first time. 

TODO:

- Sign use Apple Creds from .env file.

## PUBLISH 

Deploys whatever is configured by the Root TASK file. 

TODO:

- Roundtrip to your desktop. So your local binaries will self update from committing into Github.

- Deploy to Hetzner VPS. 