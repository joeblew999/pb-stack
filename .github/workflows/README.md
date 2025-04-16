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


## MIRROR 

Mirrors the repo to Codeberg, as a backup, cause its easy to get deleted from Github.

Its currently turned off until i can get this working.

REFS:

- https://codeberg.org/Recommendations/Mirror_to_Codeberg

TODO:

- Start using SOPS for secrets, so its much easier, and we dont need all this shit inside github secrets.

- Might be much easier to use a golang tool to do this.

- Later get this working with SoftServer and Wishes
  - https://github.com/charmbracelet/wish
  - https://github.com/charmbracelet/soft-serve
  