# github workflows

Its a mono repo, driven by TASK files, so its easy to deploy as many projects as you want from this repo with a single setup.

## CI

Builds whatever is configured by the Root TASK file.

The binaries are stored in <github> but not signed for Apple or Microsoft, and so the Pack manifest will call the OS to ask for them to be run.

TODO

- Signing use Apple Creds in .env

## CD 

Deploys whatever is configured by the Root TASK file. 

TODO:

- Roundtrip to your desktop. So your local binaries will self update from committing into Github.

- Deploy to Hetzner VPS. 