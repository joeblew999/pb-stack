# pb-stack

https://github.com/joeblew999/pb-stack

Gerard Webb

Describes a Stack based on Pocketbase that where:

1. Server Golang Code is generated as much as possible, so that developers are not hand writing code causing bugs and security leaks, allowing rapid extension of the system. A Fix only requires a change to the generator allowing rapid remediation.

2. Authentication and Authorisation is 100% controlled by the database.

3. GUI is 100% controlled using the HTMX principles.  Web and Native ( Desktop and  app ) is based off the same code, using WebViews and DeepLinks to align the Web with Native.

4. Both Cloud and On Premise so that Organisations can control their own data. 

## Stack

The following is what i currently use.

I am currently porting the Makefiles to Taskfiles.

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

Provides Synchronisation of the Pocketbase DB and Files in a master / master approach. 

Marmot runs as a side car ( using PC ). 

A Global NATS Jetstream cluster is the central rendezvous point.

Features:

- Scale out - The Load balancer automatically forwards any request to the nearest Data Center.

- Network failure tolerant - The NATS Cluster will ensure any PB will catchup. 

### DataStar ( DS )

https://github.com/starfederation/datastar

Datastar brings the functionality provided by libraries like Alpine.js (frontend reactivity) and htmx (backend reactivity) together, into one cohesive solution. It’s a lightweight, extensible framework that allows you to:

Manage state and build reactivity into your frontend using HTML attributes.

Modify the DOM and state by sending events from your backend.

With Datastar, you can build any UI that a full-stack framework like React, Vue.js or Svelte can, but with a much simpler, hypermedia-driven approach.


It reached API stability on 1st April, 2024: https://github.com/starfederation/datastar/releases/tag/v1.0.0-beta.11

Tests: https://data-star.dev/tests

Bundler: https://data-star.dev/bundler

Examples: https://data-star.dev/examples


NATS Jetstream ( and the NATS CLI ) can make calls into it, allowing many PB systems to be easily composed together globally.

## Code generation

Task files to drive the code generators as linear steps.

Development:

Its important to understand that this will be an iterative process where as we go through this we will go back and refactor these generators.

We can refactor without fear off causing refactor pain for our Developer users, because Code generation allows us to update our developers code, saving painful refactoring cycles. 

### STEP 1

Generate the golang code off the PB system.

https://github.com/Snonky/pocketbase-gogen

### STEP 2 - Open API M2M

Generate the Open API, so that the system can be used for M2M use cases.

TODO

- Use https://github.com/ogen-go/example/tree/main as a reference. It has many examples which gives me confidence that this is the right system to use.

- Write a generator that creates the Open API file off the PB Schema.

We do not know yet how far we can take this ... Lets find out.

### STEP 3 - Open API DATASTAR

Generate an OpenAPI that incorporates SSE, so that is can be used with DS.

TODO:

1. Use https://github.com/ogen-go/ogen/issues/1375#issuecomment-2766653824 as a reference as it shows the that Open API can works with SSE. 


### STEP 4 - DATASTAR WEB GUI

- Generate the DataStar, so that developers can easily write DS Web GUI.

Once we get into it we will see obvious things to code generator like:

1. Each PB Table needs a real time editor: 





















