# kamal

TODO

- TASK file 
- TOFO file

## kamal-proxy

Kamal allows doing Blue / Green Upgrades in place with zero downtime.

You just need a Docker and nothing else with your Compose configuration and binaries inside. The base TASK Pack Command does this for us.

https://github.com/basecamp/kamal-proxy
https://github.com/basecamp/kamal-proxy/releases/tag/v0.8.7

Example: https://github.com/basecamp/kamal-proxy/tree/main/example

The Logic is
- Green is deployed as a Docker but not exposed
- then Blue is drained.
- then Green is migrated.
- then Green is exposed.

When used with Tofo and Marmot, we do the same logic but allow for how many Servers there are globally.

We can add Caddy modules on top also and that will reverse Proxy into the Docker being run by Kamel.







