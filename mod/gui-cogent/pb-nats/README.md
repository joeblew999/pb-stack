# pb-nats

pb-nats is a golang project that runs Pocketbase ( https://github.com/pocketbase/pocketbase) and nats.go ( https://github.com/nats-io/nats.go )  so that changes to Pocketbase are streamed to NATS Jetstream, ready to be used by Benthos ( https://github.com/redpanda-data/benthos ).

The backend and frontend are decoupled via NATS Jetstream.

## Taskfile

A Taskfile kicks off the 2 servers with the frontend and backend command.

## Backend folder

Pocketbase runs here.

The main.go runs all the golang code.

## Frontend folder

The DataStar Web severs runs here.

We use a Datastar ( https://github.com/starfederation/datastar ) based Web GUI for managing Pocketbase.

golang modules used are:

"github.com/go-chi/chi/v5"

"github.com/starfederation/datastar"

The main.go runs all the golang code.

The index.html has the datastar javascript and tailwind css.

## Environment variables

Allow these to be overridden:

- NATS_SERVER_URL is the NATS server URL. Default: nats://localhost:4222

- POCKETBASE_SERVER_URL is the Pocketbase server URL. Default: http://localhost:8090

- POCKETBASE_ADMIN_USER is the Pocketbase admin user. Default: admin

- POCKETBASE_DATA_PATH is the pocketbase data directory. Default: ./data
