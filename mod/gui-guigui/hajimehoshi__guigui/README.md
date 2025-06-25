# README

Various projects using guigui ( https://github.com/hajimehoshi/guigui ) and datastar ( https://github.com/starfederation/datastar ) to built real time streaming GUI for web and desktop.

We use golang's go.work to separate out the modules.

We use taskfile ( https://github.com/taskfile/taskfile ) and process-compose ( https://github.com/F1bonacc1/process-compose  ) to build and run the projects.

The taskfile in this folder uses a higher level Taskfile via the taskfile include directive. It provides many common tasks that a developer needs, like git, go builds and running.  Its a bit of a hack, but it works. Maybe Augment can help us here ?

## All Golang based

No python, no OS dependent tools are used.

Because we have so much reuse of golang tools, we can have a very consistent and reproducible dev environment. But we need a shared area for these binaries and hence share task files.

One way is to build to the bin folder, and then copy them up to a higher folder level for reuse. 

This pattern can work for everything, including tools that are golang based. 

We can then keep the binaries out of the GOBIN path. 

There maybe a way for Task and Process compose to see all those binaries and use them.

Personally i prefer to use the .bin and .dep folder concept to describe what is a binaries and what are dependencies. 

The idea is that NATS Jetstream is all that is needed, because it can store the binaries and dependencies in a global place.

Devs and Users are then just consumers of those binaries.
Devs can create new binaries, and Users can consume them.

NATS and DataStar work together because changes in NATS can trigger a rebuild of the UI. 


Maybe NATS and narun (https://github.com/akhenakh/narun ) and 
https://github.com/akhenakh/nats-s3  can help with this ?  It can be a global place to store and co-ordinate everything in DEV and PROD.

Maybe the Augment AI can then use these constructs to make it easier for it to control things. NATS could control Process compose even ?







### Roadmap

do not do this yet until the taskfiles and shared bin binaries is stables.

https://github.com/nlepage/go-wasm-http-server and module github.com/nlepage/go-wasm-http-server/v2


docs. Make a docs folder thats designed or end users to use this system.  Easiest might be hugo ? 

Use NATS to help with controlling various things.

For NATS controlling everything.

https://github.com/akhenakh/narun

https://github.com/akhenakh/nats-s3

For MCP and A2A later:

https://github.com/sinadarbouy/mcp-nats

https://github.com/hwclass/a2a-bus








































































































