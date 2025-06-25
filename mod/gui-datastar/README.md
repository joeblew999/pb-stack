# datastar


A HTMX web system that uses SSE for real time Web GUI.

There is no main.go for datastar. there could be If the runtime has typescript.

```sh
go get -tool github.com/a-h/templ/cmd/templ@latest
```

It can run an embedded NATS, so that you can send and review with nats KV, for example.

This makes it self autonomous, in that the NATS Leaf Server can store anything we want and run offline and then catchup when the network resolves.

code:

https://github.com/starfederation/datastar
https://github.com/starfederation/datastar/releases/tag/v1.0.0-beta.11

docs:

There is no main cli, just an SDK at: 

https://github.com/starfederation/datastar/tree/develop/sdk/go

STATUS: broken for now.

## todo

datastar.taskfile.yml needs refactoring to work off base paths.

## usage

I need to make a base package..


## examples

https://github.com/starfederation/datastar/tree/develop/examples





