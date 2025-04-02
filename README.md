# pb-stack

Describes a Stack based on Pocketbase.

A Forward engineering stack where:

1. Server Golang Code is generated as much as possible, so that developers are not hand writing code causing bugs and security leaks, and allowing rapid extension of the system. This makes security audits much easier because its the code generator that is determining the security of the system. A Fix only requires a change to the generator allowing rapid remediation.

2. Authentication and Authorisation is 100% controlled by the database, so that its highly secure.

3. GUI is 100% controlled using the HTMX principles, with Web, Desktop and Mobile is based off the same code, so development is rapid, security is highest possible. 


## Stack

This describes all parts to build a system that allows both Cloud and On Premise so that Organisations can control their own data. 

### Process Compose ( PC ) 

https://github.com/F1bonacc1/process-compose

Provides a configuration based way to run many processes on Desktops and Servers such that:

1. Provides the ability to easily compose the system together with the parts running where needed.

2. Removes the need for K8, Docker, Docker compose, so that its simple and easy to run anywhere.

On the Cloud Servers:

- PC runs on the Host ( VPS Server ), managing the docker runtime

- PC runs inside each Docker starting up the binaries.

On the Edge / On Premise Servers:

- A Single instance of PC runs and is identical to the PC used on the Cloud Servers.

TODO: 

- Add an OS level Service runner for Desktops and Servers, so that it can be booted up by the OS easily and in a consistent way. 

### TaskFile ( TF )

https://github.com/go-task/task

Provide an alternative to Makefiles that runs on all Desktops and Servers.

This augments what PC does, allowing commands to be run. 

It has a Web GUI and Terminal GUI. 

The Web GUI reflects the commands into a Web GUI, allowing a Web based provisioning system to aid with debugging. The Web GUI is not needed once TOFU takes over and drives the provisions and configuration.

TODO:

- Build up Remote task files to control the Binaries of the Stack, so that Developers and Users have a consistent set of commands that conform to the conventions of the Stack.

### Tofu ( TOFU )

https://github.com/opentofu/opentofu

Provides a way to deploy the Stack in a consistent way to any cloud.

TOFU calls into the PC and TF sub systems to automate the provisioning and configuration of the systems.

It is idempotent, so that changes can be made and reconciled.

### Cloudflare tunnel ( CF Tunnel )

https://github.com/cloudflare/cloudflared

Provides a network tunnel to Cloud VPS Servers and Edge / On Premise systems, so that:

1. Systems are better protected from DDOS attacks.

2. Edge Servers and On Premise Servers are easy to expose and reduces security risks.


### Pocketbase ( PB )

https://github.com/pocketbase/pocketbase

Provides a SQLite DB and Web GUI Editor to allow developers to very quickly develop a database structure.

TODO: 

- Extend so that the Admin settings are part of the API, so that TOFU can manage it.


### Marmot ( MA )

https://github.com/maxpert/marmot 

Provides Synchronisation of the Pocketbase DB and Files in a master / master approach using CRDT. 

This uses NATS Jetstream.

### DataStar ( DS )

https://github.com/starfederation/datastar

It reached API stability on 1st April, 2024: https://github.com/starfederation/datastar/releases/tag/v1.0.0-beta.11

It has proper tests, which HTMX does not.: https://data-star.dev/tests/aliased

It is 100% based on SSE, which is the way forward. 

Standard HTMX requires alpine.js in order to have client side features, but DS is a single JS for everything.

Its works with NATS Jetstream, allowing many PB systems to be composed together globally.

- NATS CLI can be used to interact with it.

## Code generation

Task files to drive the code generators as linear steps.

Development:

Its important to understand that this will be an iterative process where as we go through this we will go back and refactor these generators.

We can refactor without fear off causing refactor pain for our Developer users, because Code generation allows us to update our developers code, saving painful refactoring cycles. 

### STEP 1

Generate the golang code off the PB system.

https://github.com/Snonky/pocketbase-gogen

### STEP 2

Generate the Open API, so that the system can be used for M2M use cases.

TODO

- Use https://github.com/ogen-go/example/tree/main as a reference. It has many examples which gives me confidence that this is the right system to use.

- Write a generator that creates the Open API file off the PB Schema.

We do not know yet how far we can take this ... Lets find out.

### STEP 3

Generate an OpenAPI that incorporates SSE, so that is can be used with DS.

TODO:

1. Use https://github.com/ogen-go/ogen/issues/1375#issuecomment-2766653824 as a reference as it shows the that Open API can works with SSE. 

- Generate the DataStar, so that developers can easily write DS Web GUI.

Once we get into it we will see obvious pattern like:

1. Each PB Tables needs a real time editor: 

















