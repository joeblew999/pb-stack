# pb-stack

https://github.com/joeblew999/pb-stack

This project came about due to the need for Developers to be Users and Users to be Developers. 

Hopefully the way things work will make sense to you based on this need.


A Stack based on Pocketbase ( PB )  where:

1. Server Golang Code is generated as much as possible, so that developers are not hand writing code causing bugs and security leaks, allowing rapid extension of the system. A Fix only requires a change to the generator allowing rapid remediation.

2. Authentication and Authorisation is 100% controlled by the PB database. All in one place.

3. GUI is 100% controlled using the HTMX principles, using DataStar.  Web and Native ( Desktop and  app ) is based off the same code, using WebViews and DeepLinks to align the Web with Native. DataStar and Webviews does this.

4. Both Cloud and On Premise so that Organisations can control their own data. Cloud Flare tunnel does this.

5. Each Developer and / or User can pick a Data Center within their region for GDPR, and then choose other regions for replication. Marmot does this.

You do not have to use all the stack. Its a composition where you can pick what you want.

## Repo Structure

I use mono repos, cause life is short. Long story.

Users and Developers have the exact same structure, but Users will not have the Modules, unless they are developing their own Modules of course.

So as you might have guessed this is Git OPS based, where each user has Git on their Desktop. 

- Root has the root TASK file that dictates everything else.

- MOD ( Modules ) has folders for each Stack part and a TASK file for running it.  You can make PlayGround here too to experiment and refactor.

- PROJ ( Projects ) has folders for each Project / PlayGround.  These use the Modules via the TASK files.

## Project Structure.

A Playground is a Project or Module experiment.

Each PlayGround has a TASK file, with the common TASK file manipulating these git folder:

- .bin for produced binaries.
- .dep for consumed binaries.
- .pack for packing of everything for deployment.
- .src for source code.

Task files do 2 things:

- Build the binary of what they represent.
- Run the binary of what they represent. 

We can embed Task in the main binary also

https://github.com/hypershift-community/hyper-console

We can also embed NATS.go, so that its reachable over NATS Servers, and thus eliminates SSH or HTTP for management, which is a security nightmare.


## Stack of Modules


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

Provides an alternative to Makefiles that runs on all Desktops and Servers.

Augments what PC does, allowing commands to be run. 

Has a Web GUI and Terminal GUI. 

The Web GUI reflects the commands into a Web GUI, allowing a Web based provisioning system to aid with debugging. The Web GUI is not needed once TOFU takes over and drives the provisions and configuration.

Used locally and in github actions. Single source of truth.

Can be shared across repositories using Remote Task files.


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

- Add Generator: https://github.com/Snonky/pocketbase-gogen

- Add OpenAPI: https://github.com/ogen-go/ogen/issues/1375#issuecomment-2772711703

- Structure to be like https://github.com/go-goyave/goyave 

- Extend so that the Admin settings are part of the API, so that TOFU can manage it.


### Marmot ( MA )

https://github.com/maxpert/marmot 

Example: https://github.com/maxpert/marmot-pocketbase-flyio

Provides Synchronisation of the Pocketbase DB and Files in a master / master approach. 

Marmot runs as a side car ( using PC ). 

A Global NATS Jetstream cluster is the central rendezvous point.

Features:

- Scale out - The Load balancer automatically forwards any request to the nearest Data Center.

- Network failure tolerant - The NATS Cluster will ensure any PB will catchup. 

TODO

- A PB Schema change requires a global "stop the world and sync " to be co-ordinated. Write a basic CLI to do this, which will be put into PC and later used by TOFU to make it automatic.

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

We might then want to add code generation of the standard things like:

- View a Table:  https://data-star.dev/examples/infinite_scroll

- Edit a Row in a Table: https://data-star.dev/examples/click_to_edit

- Delete a row in a table: https://data-star.dev/examples/bulk_update






















