# ds task

Using Datastar to manage repositories.

## Running

Taskfile in the root can run everything.

## Stack

https://github.com/starfederation/datastar

Examples: 

- https://github.com/starfederation/datastar/tree/develop/examples/go/hello-world

- https://github.com/starfederation/datastar/tree/develop/site

## Pattern

Datastar can update the Web GUI, in response to changes on the backend. 

We express the Operations in the Web GUI with Datastar using basic CQRS patterns.

Datastar works wel with NATS Jetstream also.

The Go backend expose specific HTTP API endpoints that the Datastar-powered frontend will interact with for these operations.

A user initiates an operation by filling a form, and then the operation matching the form occurs on the server, and then the Web GUI updates to reflect this.

## Operations

Standard git operations, assuming git is installed. clone, commit, push, status, pull

